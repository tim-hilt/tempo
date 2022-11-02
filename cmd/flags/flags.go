package flags

import (
	"github.com/spf13/viper"
)

var (
	Loglevel       int
	DailyWorkhours int
	JiraUser       string
	Password       string
	NotesDir       string
)

func SetFlagvars() {
	Loglevel = viper.GetInt("loglevel")
	DailyWorkhours = viper.GetInt("dailyWorkhours")
	JiraUser = viper.GetString("jiraUser")
	Password = viper.GetString("password")
	NotesDir = viper.GetString("notesDir")
}
