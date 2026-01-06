package utils

import (
	"os"
)

type (
	Env struct{}
)

func (e *Env) GetEnv(key string) (bool, string) {
	value, exists := os.LookupEnv(key)
	return exists, value
}

func (e *Env) SetEnv(key string, value string) error {
	return os.Setenv(key, value)
}

func (e *Env) Log() {
	Log.Msg("Environment Variables:")

	v, _ := e.GetEnv("rate_limit_requests")
	Log.Slice("rate_limit_request", v)

	v, _ = e.GetEnv("rate_limit_time")
	Log.Slice("rate_limit_time", v)

	v, _ = e.GetEnv("rate_limit_scope")
	Log.Slice("rate_limit_scope", v)

}
