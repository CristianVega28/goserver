package middleware

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/CristianVega28/goserver/core/models"
)

type (
	SecurityMiddleware struct{}
)

func (security *SecurityMiddleware) Csrf() MiddlewareFunction {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
}

func (security *SecurityMiddleware) Cors() MiddlewareFunction {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
}

func (security *SecurityMiddleware) RateLimit() MiddlewareFunction {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			rateLimit := models.RateLimit{}
			rateLimit.SetTableName("rate_limits")
			rateLimit.SetPrimaryKey("ip")

			columns := rateLimit.ParserColumn(rateLimit.GetMigration())
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			model, _ := rateLimit.Select(host, columns)

			if len(model) == 0 {
				rateLimit.CurrentCount = 1
				rateLimit.LastCount = 0
				rateLimit.TimestampStart = time.Now().UnixMilli()
				rateLimit.Ip = host

				rateLimit.InsertData()

			} else {
				currentCount := model[0]["current_count"].(int64)

				rateLimit.CurrentCount = int(currentCount) + 1
				rateLimit.LastCount = int(model[0]["last_count"].(int64))

				timestamp, _ := strconv.ParseInt(model[0]["timestamp_start"].(string), 10, 64)
				rateLimit.TimestampStart = timestamp
				rateLimit.Ip = host
				rateLimit.UpdateData(host)

			}

			// overlap :=
			f(w, r)
		}

	}
}
