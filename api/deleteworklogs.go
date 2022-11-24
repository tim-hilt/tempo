package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/api/paths"
	"golang.org/x/sync/errgroup"
)

func (a Api) DeleteWorklogs(worklogs *[]SearchWorklogsResult) error {

	errs, _ := errgroup.WithContext(context.Background())

	for _, worklog := range *worklogs {
		worklog := worklog // Necessary as of https://go.dev/doc/faq#closures_and_goroutines
		errs.Go(func() error {
			worklogId := fmt.Sprint(worklog.TempoWorklogId)
			log.Info().
				Str("ticket", worklog.Issue.Ticket).
				Str("description", worklog.Description).
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
					Str("error", fmt.Sprintf("%+v", errResponse)).
					Str("ticket", worklog.Issue.Ticket).
					Str("description", worklog.Description).
					Msg("unexpected http-status when deleting worklog")
				return errors.New("error when deleting worklog")
			}

			log.Info().
				Str("ticket", worklog.Issue.Ticket).
				Str("description", worklog.Description).
				Msg("finished deleting worklog")
			return nil
		})
	}
	return errs.Wait()
}
