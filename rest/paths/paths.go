/**
 * The paths were defined according to https://www.tempo.io/server-api-documentation/timesheets
 */
package paths

import "github.com/tim-hilt/tempo/util/config"

func UserIdPath() string {
	return config.GetHost() + "/rest/api/2/myself"
}

func CreateWorklogPath() string {
	return config.GetHost() + "/rest/tempo-timesheets/4/worklogs"
}

func FindWorklogsPath() string {
	return config.GetHost() + "/rest/tempo-timesheets/4/worklogs/search"
}

func DeleteWorklogPath(worklogId string) string {
	return config.GetHost() + "/rest/tempo-timesheets/4/worklogs/" + worklogId
}
