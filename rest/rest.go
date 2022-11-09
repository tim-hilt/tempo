package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/paths"
	"github.com/tim-hilt/tempo/util"
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

type SearchWorklogsResult struct {
	TempoWorklogId  int    `json:"tempoWorklogId"`
	DurationSeconds int    `json:"timeSpentSeconds"`
	Issue           issue  `json:"issue"`
	DateTime        string `json:"dateCreated"`
}

func (a *Api) FindWorklogsInRange(from string, to string) (*[]SearchWorklogsResult, error) {
	log.Info().Str("from", from).Str("to", to).Msg("started searching for worklogs")
	resp, err := a.client.R().
		SetBody(searchWorklogBody{From: from, To: to, Users: []string{a.UserId}}).
		SetResult([]SearchWorklogsResult{}).
		Post(paths.FindWorklogsPath())
	status := resp.StatusCode()

	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, errors.New("error when searching for worklogs in range " + from +
			" to " + to + ": Response was HTTP-status " + fmt.Sprint(status))
	}

	log.Info().Str("from", from).Str("to", to).Msg("finished searching for worklogs")

	worklogs := resp.Result().(*[]SearchWorklogsResult)
	return worklogs, nil
}

func (a *Api) findWorklogIdsOn(day string) (*[]SearchWorklogsResult, error) {
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
	errs, _ := errgroup.WithContext(context.Background())

	for _, worklog := range *worklogs {
		worklog := worklog // Necessary as of https://go.dev/doc/faq#closures_and_goroutines
		errs.Go(func() error {
			now := time.Now()
			worklogCreated, err := time.Parse(util.DATE_TIME_FORMAT, worklog.DateTime)
			if !worklogCreated.Before(now.Add(time.Duration(-1) * time.Minute)) {
				// If worklog isn't older than one minute, then don't delete it
				// Prevents race conditions with submitting
				return nil
			}

			if err != nil {
				return err
			}

			worklogId := fmt.Sprint(worklog.TempoWorklogId)
			log.Info().Str("ticket", worklog.Issue.Ticket).Str("description", worklog.Issue.Description).Msg("started deleting worklog")

			// TODO: Only delete if worklog was created > 5 minutes in the past
			resp, err := a.client.R().Delete(paths.DeleteWorklogPath(worklogId))
			status := resp.StatusCode()

			if err != nil {
				return err
			} else if status != http.StatusNoContent {
				return errors.New("HTTP-response was " + fmt.Sprint(status))
			}

			log.Info().Str("ticket", worklog.Issue.Ticket).Msg("finished deleting worklog")
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
	log.Info().Str("ticket", ticket).Msg("start creating worklog")
	day = strings.TrimSuffix(day, ".md")
	worklog := worklog{Ticket: ticket, Comment: comment, Seconds: seconds, Day: day, UserId: a.UserId}

	resp, err := a.client.R().
		SetBody(worklog).
		Post(paths.CreateWorklogPath())
	status := resp.StatusCode()

	if err != nil {
		return err
	} else if status != http.StatusOK {
		return errors.New("error when creating worklog for ticket " + ticket + ": HTTP-status was " + fmt.Sprint(status) + " with body " + fmt.Sprint(worklog))
	}

	log.Info().Str("ticket", ticket).Msg("finished creating worklog")

	return nil
}
