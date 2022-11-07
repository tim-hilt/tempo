package paths

import "github.com/tim-hilt/tempo/util/config"

func UserIdPath() string {
	return config.GetConfigParams().JiraHost + "/rest/api/2/myself"
}

func CreateWorklogPath() string {
	return config.GetConfigParams().JiraHost + "/rest/tempo-timesheets/4/worklogs"
}

func FindWorklogsPath() string {
	return config.GetConfigParams().JiraHost + "/rest/tempo-timesheets/4/worklogs/search"
}

func DeleteWorklogPath(worklogId string) string {
	return config.GetConfigParams().JiraHost + "/rest/tempo-timesheets/4/worklogs/" + worklogId
}
