package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/paths"
	"github.com/tim-hilt/tempo/util/config"
	"golang.org/x/sync/errgroup"
)

type Api struct {
	client *resty.Client
	UserId string
}

func New(user string, password string) *Api {
	apiClient := resty.New()
	apiClient.SetBasicAuth(user, password)
	apiClient.SetDebug(config.DebugEnabled())

	tempo := &Api{client: apiClient}
	err := tempo.initUser()

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("error when getting user-id")
	}

	return tempo
}

type userIdResponse struct {
	UserId string `json:"key"`
}

type errorResponse struct {
	Errors struct {
		TimeSpentSeconds string `json:"timeSpentSeconds"`
	} `json:"errors"`
	ErrorMessages []string `json:"errorMessages"`
	Reasons       []string `json:"reasons"`
}

func (b *Api) initUser() error {
	log.Info().Msg("Started getting userId")

	resp, err := b.client.R().
		SetResult(userIdResponse{}).
		SetError(errorResponse{}).
		Get(paths.UserIdPath())
	status := resp.StatusCode()

	if err != nil {
		return err
	} else if status >= http.StatusBadRequest {
		errResponse := resp.Error().(*errorResponse)
		log.Trace().
			Int("httpStatus", status).
			Str("body", fmt.Sprintf("%+v", errResponse)).
			Msg("unexpected http-status when searching for userId")
		return errors.New("error when getting userId")
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
	log.Info().
		Str("from", from).
		Str("to", to).
		Msg("started searching for worklogs")
	resp, err := a.client.R().
		SetBody(searchWorklogBody{
			From:  from,
			To:    to,
			Users: []string{a.UserId},
		}).
		SetResult([]SearchWorklogsResult{}).
		SetError(errorResponse{}).
		Post(paths.FindWorklogsPath())
	status := resp.StatusCode()

	if err != nil {
		return nil, err
	} else if status >= http.StatusBadRequest {
		errResponse := resp.Error().(*errorResponse)
		log.Trace().
			Int("status", status).
			Str("body", fmt.Sprintf("%+v", errResponse)).
			Str("from", from).
			Str("to", to).
			Msg("unexpected http-status when searching for worklogs")
		return nil, errors.New("error when searching for worklogs")
	}

	log.Info().
		Str("from", from).
		Str("to", to).
		Msg("finished searching for worklogs")

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
	worklogs, err := a.findWorklogIdsOn(day)
	if err != nil {
		return err
	}
	errs, _ := errgroup.WithContext(context.Background())

	for _, worklog := range *worklogs {
		worklog := worklog // Necessary as of https://go.dev/doc/faq#closures_and_goroutines
		errs.Go(func() error {
			worklogId := fmt.Sprint(worklog.TempoWorklogId)
			log.Info().
				Str("ticket", worklog.Issue.Ticket).
				Str("description", worklog.Issue.Description).
				Msg("started deleting worklog")

			resp, err := a.client.R().
				SetError(errorResponse{}).
				Delete(paths.DeleteWorklogPath(worklogId))
			status := resp.StatusCode()

			if err != nil {
				return err
			} else if status >= http.StatusBadRequest {
				errResponse := resp.Error().(*errorResponse)
				log.Trace().
					Int("status", status).
					Str("body", fmt.Sprintf("%+v", errResponse)).
					Str("day", day).
					Str("ticket", worklog.Issue.Ticket).
					Str("comment", worklog.Issue.Description).
					Msg("unexpected http-status when deleting worklog")
				return errors.New("error when deleting worklog")
			}

			log.Info().
				Str("ticket", worklog.Issue.Ticket).
				Msg("finished deleting worklog")
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
	log.Info().
		Str("ticket", ticket).
		Msg("start creating worklog")
	worklog := worklog{
		Ticket:  ticket,
		Comment: comment,
		Seconds: seconds,
		Day:     day,
		UserId:  a.UserId,
	}

	resp, err := a.client.R().
		SetBody(worklog).
		SetError(errorResponse{}).
		Post(paths.CreateWorklogPath())
	status := resp.StatusCode()

	if err != nil {
		return err
	} else if status >= http.StatusBadRequest {
		errResponse := resp.Error().(*errorResponse)
		log.Trace().
			Str("ticket", ticket).
			Str("comment", comment).
			Str("day", day).
			Int("httpStatus", status).
			Str("body", fmt.Sprintf("%+v", errResponse)).
			Msg("unexpected http-status when creating worklog")
		return errors.New("error when creating worklog")
	}

	log.Info().
		Str("ticket", ticket).
		Msg("finished creating worklog")

	return nil
}
