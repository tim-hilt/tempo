package paths

import "github.com/tim-hilt/tempo/util"

func UserIdPath() string {
	return util.GetConfigParams().JiraHost + "/rest/api/2/myself"
}

func CreateWorklogPath() string {
	return util.GetConfigParams().JiraHost + "/rest/tempo-timesheets/4/worklogs"
}

func FindWorklogsPath() string {
	return util.GetConfigParams().JiraHost + "/rest/tempo-timesheets/4/worklogs/search"
}

func DeleteWorklogPath(worklogId string) string {
	return util.GetConfigParams().JiraHost + "/rest/tempo-timesheets/4/worklogs/" + worklogId
}
