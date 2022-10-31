package tempo

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/rest/endpoints"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/logging"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logging.SetLoglevel()
}

type Tempo struct {
	// TODO: Should this functionality be factored away? Separation of concerns! A Book instance doesn't need to know about REST
	tempoApiClient *resty.Client
	UserId         string
}

func New(user string, password string) *Tempo {
	apiClient := resty.New()
	apiClient.SetBasicAuth(user, password)

	tempo := &Tempo{tempoApiClient: apiClient}
	tempo.initUser()

	return tempo
}

type userIdResponse struct {
	UserId string `json:"key"`
}

func (b *Tempo) initUser() {
	log.Info().Msg("Started getting userId")

	resp, err := b.tempoApiClient.R().
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

func (t *Tempo) findWorklogsInRange(from string, to string) *[]searchWorklogsResult {
	log.Info().Msg("Started searching for worklogs in range " + from + " - " + to)
	resp, err := t.tempoApiClient.R().
		SetBody(searchWorklogBody{From: from, To: to, Users: []string{t.UserId}}).
		SetResult([]searchWorklogsResult{}).
		Post(endpoints.FindWorklogsPath())

	util.HandleErr(err, "error while searching for worklogs in range "+from+" - "+to)
	util.HandleErronousHttpStatus(resp.StatusCode())
	log.Info().Msg("Finished searching for worklogs in range " + from + " - " + to)
	worklogs := resp.Result().(*[]searchWorklogsResult)
	return worklogs
}

func (t *Tempo) FindWorklogIdsOn(day string) *[]searchWorklogsResult {
	worklogs := t.findWorklogsInRange(day, day)
	return worklogs
}

func (t *Tempo) DeleteWorklogs(day string) {
	worklogs := t.FindWorklogIdsOn(day)

	for _, worklog := range *worklogs {
		worklogId := fmt.Sprint(worklog.TempoWorklogId)
		log.Info().Msg("Started deleting worklog for ticket " + worklog.Ticket + " with description: " + worklog.Description)

		resp, err := t.tempoApiClient.R().Delete(endpoints.DeleteWorklogPath(worklogId))

		util.HandleErr(err, "error while deleting worklog with id "+worklogId)
		util.HandleErronousHttpStatus(resp.StatusCode())
		log.Info().Msg("Finished deleting worklog for ticket " + worklog.Ticket)
	}
}

func (t *Tempo) CreateWorklog(ticket string, comment string, seconds int, day string) {
	log.Info().Msg("Start creating worklog for " + ticket)

	resp, err := t.tempoApiClient.R().
		SetBody(worklog{Ticket: ticket, Comment: comment, Seconds: seconds, Day: day, UserId: t.UserId}).
		Post(endpoints.CreateWorklogPath())

	util.HandleErr(err, "error when creating worklog")
	util.HandleErronousHttpStatus(resp.StatusCode())

	log.Info().Msg("Finished creating worklog for " + ticket)
}

func (t *Tempo) SubmitDay(day string) {
	ticketEntries := noteparser.ParseDailyNote(day)

	// Clean up first
	t.DeleteWorklogs(day)

	// ...then book on clean state
	workedMinutes := 0
	for _, ticket := range ticketEntries {
		t.CreateWorklog(ticket.Ticket, ticket.Comment, ticket.DurationMinutes*60, day)
		workedMinutes += ticket.DurationMinutes
	}

	hours, minutes := util.Divmod(workedMinutes, util.MINUTES_IN_HOUR)
	fmt.Println("successfully logged " + fmt.Sprint(hours) + " hours and " + fmt.Sprint(minutes) + " minutes on " + day)
}

func (t *Tempo) GetMonthlyHours() {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, -1)
	worklogs := t.findWorklogsInRange(start.Format(util.DATE_FORMAT), end.Format(util.DATE_FORMAT))

	bookedTimeSeconds := 0
	for _, worklog := range *worklogs {
		bookedTimeSeconds += worklog.DurationSeconds
	}

	hours, minutes := util.Divmod(bookedTimeSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
	fmt.Println("worked " + fmt.Sprint(hours) + " hours and " + fmt.Sprint(minutes) + " minutes in current month")
}

// TODO: Fill with life
func WatchNotes() {
	// 1. Watch note-dir
	// 2. On Create/Change...
	// 2.1. Wait 5 minutes
	// 2.2. If 5 minutes passed and no new event happened, Delete all worklogs for that day and CreateWorklog for all the other entries
}
