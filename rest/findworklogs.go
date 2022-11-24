package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/paths"
)

type searchWorklogBody struct {
	From    string   `json:"from,omitempty"`
	To      string   `json:"to,omitempty"`
	Tickets []string `json:"taskKey,omitempty"`
	Users   []string `json:"worker"`
}

func (a Api) FindWorklogs(searchBody searchWorklogBody) (*[]SearchWorklogsResult, error) {
	searchBody.Users = []string{a.UserId}
	log.Info().
		Str("query", fmt.Sprintf("%+v", searchBody)).
		Msg("started searching for worklogs")
	resp, err := a.client.R().
		SetBody(searchBody).
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
			Str("error", fmt.Sprintf("%+v", errResponse)).
			Str("query", fmt.Sprintf("%+v", searchBody)).
			Msg("unexpected http-status when searching for worklogs")
		return nil, errors.New("error when searching for worklogs")
	}

	log.Info().
		Str("query", fmt.Sprintf("%+v", searchBody)).
		Msg("finished searching for worklogs")

	worklogs := resp.Result().(*[]SearchWorklogsResult)
	return worklogs, nil
}

func (a Api) FindWorklogsInRange(from string, to string) (*[]SearchWorklogsResult, error) {
	worklogs, err := a.FindWorklogs(searchWorklogBody{
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, err
	}
	return worklogs, nil
}

func (a Api) FindWorklogsForTicket(ticket string) (*[]SearchWorklogsResult, error) {
	worklogs, err := a.FindWorklogs(searchWorklogBody{Tickets: []string{ticket}})
	if err != nil {
		return nil, err
	}
	return worklogs, nil
}
