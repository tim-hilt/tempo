package rest

import (
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/tim-hilt/tempo/rest/paths"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/logging"
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

func (b *Api) initUser() {
	logging.Logger.Info().Msg("Started getting userId")

	resp, err := b.client.R().
		SetResult(userIdResponse{}).
		Get(paths.UserIdPath())

	util.HandleResponse(resp.StatusCode(), err, "error when getting myself")

	userId := resp.Result().(*userIdResponse).UserId
	logging.Logger.Info().Msg("Finished getting userId: " + userId)
	b.UserId = userId
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

func (a *Api) FindWorklogsInRange(from string, to string) (worklogs *[]searchWorklogsResult) {
	logging.Logger.Info().Msg("Started searching for worklogs in range " + from + " - " + to)
	resp, err := a.client.R().
		SetBody(searchWorklogBody{From: from, To: to, Users: []string{a.UserId}}).
		SetResult([]searchWorklogsResult{}).
		Post(paths.FindWorklogsPath())

	util.HandleResponse(resp.StatusCode(), err, "error while searching for worklogs in range "+from+" - "+to)

	logging.Logger.Info().Msg("Finished searching for worklogs in range " + from + " - " + to)

	worklogs = resp.Result().(*[]searchWorklogsResult)
	return
}

func (a *Api) findWorklogIdsOn(day string) *[]searchWorklogsResult {
	worklogs := a.FindWorklogsInRange(day, day)
	return worklogs
}

func (a *Api) DeleteWorklogs(day string) {
	worklogs := a.findWorklogIdsOn(day)
	var wg sync.WaitGroup

	for _, worklog := range *worklogs {
		wg.Add(1)
		go func(worklog searchWorklogsResult) {
			defer wg.Done()
			worklogId := fmt.Sprint(worklog.TempoWorklogId)
			logging.Logger.Info().Msg("Started deleting worklog for ticket " + worklog.Issue.Ticket + " with description: " + worklog.Issue.Description)

			resp, err := a.client.R().Delete(paths.DeleteWorklogPath(worklogId))
			util.HandleResponse(resp.StatusCode(), err, "error while deleting worklog with id "+worklogId)

			logging.Logger.Info().Msg("Finished deleting worklog for ticket " + worklog.Issue.Ticket)
		}(worklog)
	}
	wg.Wait()
}

type worklog struct {
	Ticket  string `json:"originTaskId"`
	Comment string `json:"comment"`
	Seconds int    `json:"timeSpentSeconds"`
	Day     string `json:"started"`
	UserId  string `json:"worker"`
}

func (a *Api) CreateWorklog(ticket string, comment string, seconds int, day string) {
	logging.Logger.Info().Msg("Start creating worklog for " + ticket)

	resp, err := a.client.R().
		SetBody(worklog{Ticket: ticket, Comment: comment, Seconds: seconds, Day: day, UserId: a.UserId}).
		Post(paths.CreateWorklogPath())

	util.HandleResponse(resp.StatusCode(), err, "error when creating worklog")

	logging.Logger.Info().Msg("Finished creating worklog for " + ticket)
}
