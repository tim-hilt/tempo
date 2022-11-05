package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/paths"
	"golang.org/x/sync/errgroup"
)

type Api struct {
	client *resty.Client
	UserId string
}

func New(user string, password string) *Api {
	apiClient := resty.New()
	apiClient.SetBasicAuth(user, password)

	tempo := &Api{client: apiClient}
	tempo.initUser()

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
	log.Info().Msg("Finished getting userId: " + userId)
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
	log.Info().Msg("Started searching for worklogs in range " + from + " - " + to)
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

	log.Info().Msg("Finished searching for worklogs in range " + from + " - " + to)

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
	worklogs, err := a.findWorklogIdsOn(day)
	if err != nil {
		return err
	}
	errs, _ := errgroup.WithContext(context.Background())

	for _, worklog := range *worklogs {
		errs.Go(func() error {
			worklogId := fmt.Sprint(worklog.TempoWorklogId)
			log.Info().Msg("Started deleting worklog for ticket " + worklog.Issue.Ticket +
				" with description: " + worklog.Issue.Description)

			resp, err := a.client.R().Delete(paths.DeleteWorklogPath(worklogId))
			status := resp.StatusCode()

			if err != nil {
				return err
			} else if status != http.StatusOK {
				return errors.New("error when deleting worklog for ticket " + worklog.Issue.Ticket +
					": HTTP-response was " + fmt.Sprint(status))
			}

			log.Info().Msg("Finished deleting worklog for ticket " + worklog.Issue.Ticket)
			return nil
		})
	}
	return errs.Wait()
}

type worklog struct {
	Ticket  string `json:"originTaskId"`
	Comment string `json:"comment"`
	Seconds int    `json:"timeSpentSeconds"`
	Day     string `json:"started"`
	UserId  string `json:"worker"`
}

func (a *Api) CreateWorklog(ticket string, comment string, seconds int, day string) error {
	log.Info().Msg("Start creating worklog for " + ticket)

	resp, err := a.client.R().
		SetBody(worklog{Ticket: ticket, Comment: comment, Seconds: seconds, Day: day, UserId: a.UserId}).
		Post(paths.CreateWorklogPath())
	status := resp.StatusCode()

	if err != nil {
		return err
	} else if status != http.StatusOK {
		return errors.New("error when creating worklog for ticket " + ticket + ": HTTP-status was " + fmt.Sprint(status))
	}

	log.Info().Msg("Finished creating worklog for " + ticket)

	return nil
}
