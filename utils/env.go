package utils

import "os"

type (
	Env struct{}
)

var lgbk LoggerI = &Logger{}
var log Logger = lgbk.Create()

func (e *Env) GetEnv(key string) (bool, string) {
	value, exists := os.LookupEnv(key)
	return exists, value
}

func (e *Env) SetEnv(key string, value string) error {
	return os.Setenv(key, value)
}

func (e *Env) Log() {
	log.Msg("Environment Variables:")

	v, _ := e.GetEnv("rate_limit_requests")
	log.Slice("rate_limit_request", v)

	v, _ = e.GetEnv("rate_limit_time")
	log.Slice("rate_limit_time", v)

	v, _ = e.GetEnv("rate_limit_scope")
	log.Slice("rate_limit_scope", v)

}
