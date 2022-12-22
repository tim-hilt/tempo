package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type columns struct {
	Tickets   string `validate:"required"`
	Comments  string `validate:"required"`
	Durations string `validate:"required"`
}

type configParams struct {
	User           string `validate:"required_without=JiraUserId TempoApiToken"`
	Password       string `validate:"required_without=JiraUserId TempoApiToken"`
	JiraUserId     string `validate:"required_with=TempoApiToken"`
	TempoApiToken  string `validate:"required_with=JiraUserId"`
	Notesdir       string `validate:"required,dir"`
	JiraHost       string `validate:"required,url"`
	Loglevel       int    `validate:"gte=-1,lte=5"`
	DailyWorkhours int    `validate:"gte=0,lte=10"`
	Columns        columns
}

func GetCredentials() (string, string) {
	return viper.GetString(USER_CONFIG_VAL), viper.GetString(PASSWORD_CONFIG_VAL)
}

func GetJiraUserId() string {
	return viper.GetString(JIRA_USER_CONFIG_VAL)
}

func GetTempoApiToken() string {
	return viper.GetString(TEMPO_TOKEN_CONFIG_VAL)
}

func HasTempoApiToken() bool {
	return viper.IsSet(JIRA_USER_CONFIG_VAL) && viper.IsSet(TEMPO_TOKEN_CONFIG_VAL)
}

func GetHost() string {
	return viper.GetString(HOST_CONFIG_VAL)
}

func GetNotesdir() string {
	notesDir := viper.GetString(NOTESDIR_CONFIG_VAL)
	if strings.HasPrefix(notesDir, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("can't get users homedir")
		}
		notesDir = homedir + notesDir[1:]
	}
	return notesDir
}

func GetWorkhours() int {
	return viper.GetInt(DAILYWORKHOURS_CONFIG_VAL)
}

func GetLoglevel() int {
	return viper.GetInt(LOGLEVEL_CONFIG_VAL)
}

func GetColumns() columns {
	return columns{
		Tickets:   viper.GetString(TICKETS_COLUMN_CONFIG_VAL),
		Comments:  viper.GetString(COMMENTS_COLUMN_CONFIG_VAL),
		Durations: viper.GetString(DURATIONS_COLUMN_CONFIG_VAL),
	}
}

func DebugEnabled() bool {
	return viper.GetBool(DEBUG_ENABLED_CONFIG_VAL)
}

func Validate() {
	user, password := GetCredentials()
	config := configParams{
		User:           user,
		Password:       password,
		JiraUserId:     GetJiraUserId(),
		TempoApiToken:  GetTempoApiToken(),
		JiraHost:       GetHost(),
		Notesdir:       GetNotesdir(),
		Loglevel:       GetLoglevel(),
		DailyWorkhours: GetWorkhours(),
		Columns:        GetColumns(),
	}
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		log.Fatal().Err(err).Msg("error when validating config-params")
	}
}
