package middleware

import (
	"fmt"
	"net/http"

	"github.com/CristianVega28/goserver/utils"
)

func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs := utils.Logger{}
		log := logs.Create()
		// output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		// log.Logger = zerolog.New(output).With().Timestamp().Logger()
		fmt.Println("login")
		log.Msg(fmt.Sprintf("Method: %s, Path: %s", r.Method, r.URL.Path))
		f(w, r)
	}
}
