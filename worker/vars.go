package worker

var (
	config                = Configure()
	https          string = "https"
	firDatabaseURL string = config.DBUrl
	apiKey         string = config.ApiKey
)
