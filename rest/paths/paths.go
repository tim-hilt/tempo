package paths

const (
	prefix         = "https://jira.in-technology.de:443/rest/" // TODO: This should be configurable
	MyselfUrl      = prefix + "api/2/myself"
	WorklogsPrefix = prefix + "tempo-timesheets/4/worklogs"
)

func UserIdPath() string {
	return MyselfUrl
}

func CreateWorklogPath() string {
	return WorklogsPrefix
}

func FindWorklogsPath() string {
	return WorklogsPrefix + "/search"
}

func DeleteWorklogPath(worklogId string) string {
	return WorklogsPrefix + "/" + worklogId
}
