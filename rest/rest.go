package rest

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/endpoints"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/logging"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logging.SetLoglevel()
}

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

func (b *Api) initUser() {
	log.Info().Msg("Started getting userId")

	resp, err := b.client.R().
		SetResult(userIdResponse{}).
		Get(endpoints.UserIdPath())

	util.HandleErr(err, "error when getting myself")
	util.HandleErronousHttpStatus(resp.StatusCode())

	userId := resp.Result().(*userIdResponse).UserId
	log.Info().Msg("Finished getting userId: " + userId)
	b.UserId = userId
}

type worklog struct {
	Ticket  string `json:"originTaskId"`
	Comment string `json:"comment"`
	Seconds int    `json:"timeSpentSeconds"`
	Day     string `json:"started"`
	UserId  string `json:"worker"`
}

type searchWorklogBody struct {
	From  string   `json:"from"`
	To    string   `json:"to"`
	Users []string `json:"worker"`
}

type searchWorklogsResult struct {
	TempoWorklogId  int    `json:"tempoWorklogId"`
	Ticket          string `json:"issue.key"`
	Description     string `json:"issue.summary"`
	DurationSeconds int    `json:"timeSpentSeconds"`
}

func (a *Api) FindWorklogsInRange(from string, to string) *[]searchWorklogsResult {
	log.Info().Msg("Started searching for worklogs in range " + from + " - " + to)
	resp, err := a.client.R().
		SetBody(searchWorklogBody{From: from, To: to, Users: []string{a.UserId}}).
		SetResult([]searchWorklogsResult{}).
		Post(endpoints.FindWorklogsPath())

	util.HandleErr(err, "error while searching for worklogs in range "+from+" - "+to)
	util.HandleErronousHttpStatus(resp.StatusCode())
	log.Info().Msg("Finished searching for worklogs in range " + from + " - " + to)
	worklogs := resp.Result().(*[]searchWorklogsResult)
	return worklogs
}

func (a *Api) findWorklogIdsOn(day string) *[]searchWorklogsResult {
	worklogs := a.FindWorklogsInRange(day, day)
	return worklogs
}

func (a *Api) DeleteWorklogs(day string) {
	worklogs := a.findWorklogIdsOn(day)

	for _, worklog := range *worklogs {
		worklogId := fmt.Sprint(worklog.TempoWorklogId)
		log.Info().Msg("Started deleting worklog for ticket " + worklog.Ticket + " with description: " + worklog.Description)

		resp, err := a.client.R().Delete(endpoints.DeleteWorklogPath(worklogId))

		util.HandleErr(err, "error while deleting worklog with id "+worklogId)
		util.HandleErronousHttpStatus(resp.StatusCode())
		log.Info().Msg("Finished deleting worklog for ticket " + worklog.Ticket)
	}
}

func (a *Api) CreateWorklog(ticket string, comment string, seconds int, day string) {
	log.Info().Msg("Start creating worklog for " + ticket)

	resp, err := a.client.R().
		SetBody(worklog{Ticket: ticket, Comment: comment, Seconds: seconds, Day: day, UserId: a.UserId}).
		Post(endpoints.CreateWorklogPath())

	util.HandleErr(err, "error when creating worklog")
	util.HandleErronousHttpStatus(resp.StatusCode())

	log.Info().Msg("Finished creating worklog for " + ticket)
}
