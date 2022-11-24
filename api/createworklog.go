package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/api/paths"
)

type worklog struct {
	Ticket      string `json:"originTaskId"`
	Description string `json:"comment"`
	Seconds     int    `json:"timeSpentSeconds"`
	Day         string `json:"started"`
	UserId      string `json:"worker"`
}

func (a Api) CreateWorklog(ticket string, description string, seconds int, day string) error {
	log.Info().
		Str("ticket", ticket).
		Msg("start creating worklog")
	worklog := worklog{
		Ticket:      ticket,
		Description: description,
		Seconds:     seconds,
		Day:         day,
		UserId:      a.UserId,
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
			Str("description", description).
			Str("day", day).
			Int("httpStatus", status).
			Str("error", fmt.Sprintf("%+v", errResponse)).
			Msg("unexpected http-status when creating worklog")
		return errors.New("error when creating worklog")
	}

	log.Info().
		Str("ticket", ticket).
		Msg("finished creating worklog")

	return nil
}
