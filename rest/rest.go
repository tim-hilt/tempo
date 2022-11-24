package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/paths"
	"github.com/tim-hilt/tempo/util/config"
)

type Api struct {
	client *resty.Client
	UserId string
}

type issue struct {
	Ticket string `json:"key"`
}

type SearchWorklogsResult struct {
	Description     string `json:"comment"`
	TempoWorklogId  int    `json:"tempoWorklogId"`
	DurationSeconds int    `json:"timeSpentSeconds"`
	Issue           issue  `json:"issue"`
	DateTime        string `json:"dateCreated"`
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

func (a *Api) initUser() error {
	log.Info().Msg("Started getting userId")

	resp, err := a.client.R().
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
			Str("error", fmt.Sprintf("%+v", errResponse)).
			Msg("unexpected http-status when searching for userId")
		return errors.New("error when getting userId")
	}

	userId := resp.Result().(*userIdResponse).UserId
	log.Info().Str("userId", userId).Msg("finished getting userId")
	a.UserId = userId

	return nil
}
