package rest

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/rest/paths"
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
		Get(paths.UserIdPath())

	util.HandleResponse(resp.StatusCode(), err, "error when getting myself")

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

func (a *Api) FindWorklogsInRange(from string, to string) (worklogs *[]searchWorklogsResult) {
	log.Info().Msg("Started searching for worklogs in range " + from + " - " + to)
	resp, err := a.client.R().
		SetBody(searchWorklogBody{From: from, To: to, Users: []string{a.UserId}}).
		SetResult([]searchWorklogsResult{}).
		Post(paths.FindWorklogsPath())

	util.HandleResponse(resp.StatusCode(), err, "error while searching for worklogs in range "+from+" - "+to)

	log.Info().Msg("Finished searching for worklogs in range " + from + " - " + to)

	worklogs = resp.Result().(*[]searchWorklogsResult)
	return
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

		resp, err := a.client.R().Delete(paths.DeleteWorklogPath(worklogId))
		util.HandleResponse(resp.StatusCode(), err, "error while deleting worklog with id "+worklogId)

		log.Info().Msg("Finished deleting worklog for ticket " + worklog.Ticket)
	}
}

func (a *Api) CreateWorklog(ticket string, comment string, seconds int, day string) {
	log.Info().Msg("Start creating worklog for " + ticket)

	resp, err := a.client.R().
		SetBody(worklog{Ticket: ticket, Comment: comment, Seconds: seconds, Day: day, UserId: a.UserId}).
		Post(paths.CreateWorklogPath())

	util.HandleResponse(resp.StatusCode(), err, "error when creating worklog")

	log.Info().Msg("Finished creating worklog for " + ticket)
}
