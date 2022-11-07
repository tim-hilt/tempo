package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/paths"
)

type Api struct {
	client *resty.Client
	UserId string
}

func New(user string, password string) *Api {
	apiClient := resty.New()
	apiClient.SetBasicAuth(user, password)

	tempo := &Api{client: apiClient}
	err := tempo.initUser()

	if err != nil {
		log.Fatal().Err(err).Msg("error when getting user-id")
	}

	return tempo
}

type userIdResponse struct {
	UserId string `json:"key"`
}

func (b *Api) initUser() error {
	log.Info().Msg("Started getting userId")

	resp, err := b.client.R().
		SetResult(userIdResponse{}).
		Get(paths.UserIdPath())
	status := resp.StatusCode()

	if err != nil {
		return err
	} else if status != 200 {
		return errors.New("error when getting userId: Response was HTTP-Status" + fmt.Sprint(status))
	}

	userId := resp.Result().(*userIdResponse).UserId
	log.Info().Str("userId", userId).Msg("finished getting userId")
	b.UserId = userId

	return nil
}

type searchWorklogBody struct {
	From  string   `json:"from"`
	To    string   `json:"to"`
	Users []string `json:"worker"`
}

type issue struct {
	Ticket      string `json:"key"`
	Description string `json:"summary"`
}

type searchWorklogsResult struct {
	TempoWorklogId  int    `json:"tempoWorklogId"`
	DurationSeconds int    `json:"timeSpentSeconds"`
	Issue           issue  `json:"issue"`
	Date            string `json:"started"`
}

func (a *Api) FindWorklogsInRange(from string, to string) (*[]searchWorklogsResult, error) {
	log.Info().Str("from", from).Str("to", to).Msg("started searching for worklogs")
	resp, err := a.client.R().
		SetBody(searchWorklogBody{From: from, To: to, Users: []string{a.UserId}}).
		SetResult([]searchWorklogsResult{}).
		Post(paths.FindWorklogsPath())
	status := resp.StatusCode()

	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, errors.New("error when searching for worklogs in range " + from +
			" to " + to + ": Response was HTTP-status " + fmt.Sprint(status))
	}

	log.Info().Str("from", from).Str("to", to).Msg("finished searching for worklogs")

	worklogs := resp.Result().(*[]searchWorklogsResult)
	return worklogs, nil
}

func (a *Api) findWorklogIdsOn(day string) (*[]searchWorklogsResult, error) {
	worklogs, err := a.FindWorklogsInRange(day, day)
	if err != nil {
		return nil, err
	}
	return worklogs, nil
}

func (a *Api) DeleteWorklogs(day string) error {
	day = strings.TrimSuffix(day, ".md")
	worklogs, err := a.findWorklogIdsOn(day)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for _, worklog := range *worklogs {
		// TODO: errgroups didn't work, because they don't allow to be parameterized.
		//       error-handling is not quite possible with go-routines. I still need
		//       a better way to solve this and handle errors on the top-level.
		wg.Add(1)
		go func(worklog searchWorklogsResult) {
			defer wg.Done()
			worklogId := fmt.Sprint(worklog.TempoWorklogId)
			log.Info().Str("ticket", worklog.Issue.Ticket).Str("description", worklog.Issue.Description).Msg("started deleting worklog")

			resp, err := a.client.R().Delete(paths.DeleteWorklogPath(worklogId))
			status := resp.StatusCode()

			if err != nil {
				log.Error().Err(err).Str("ticket", worklog.Issue.Ticket).Msg("error when deleting worklog")
			} else if status != http.StatusNoContent {
				log.Error().Err(errors.New("HTTP-response was "+fmt.Sprint(status))).Str("ticket", worklog.Issue.Ticket).Msg("error when deleting worklog")
			}

			log.Info().Str("ticket", worklog.Issue.Ticket).Msg("finished deleting worklog")
		}(worklog)
	}
	wg.Wait()
	return nil
}

type worklog struct {
	Ticket  string `json:"originTaskId"`
	Comment string `json:"comment"`
	Seconds int    `json:"timeSpentSeconds"`
	Day     string `json:"started"`
	UserId  string `json:"worker"`
}

func (a *Api) CreateWorklog(ticket string, comment string, seconds int, day string) error {
	log.Info().Str("ticket", ticket).Msg("start creating worklog")
	day = strings.TrimSuffix(day, ".md")

	resp, err := a.client.R().
		SetBody(worklog{Ticket: ticket, Comment: comment, Seconds: seconds, Day: day, UserId: a.UserId}).
		Post(paths.CreateWorklogPath())
	status := resp.StatusCode()

	if err != nil {
		return err
	} else if status != http.StatusOK {
		return errors.New("error when creating worklog for ticket " + ticket + ": HTTP-status was " + fmt.Sprint(status))
	}

	log.Info().Str("ticket", ticket).Msg("finished creating worklog")

	return nil
}
