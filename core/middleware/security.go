package middleware

import (
	"net"
	"net/http"
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

				currentCount := model[0]["current_count"].(int)
				timestart_start := model[0]["timestamp_start"].(int64)
				previousCount := model[0]["current_count"].(int)
				rateLimit.CurrentCount = currentCount
				rateLimit.LastCount = previousCount
				rateLimit.TimestampStart = timestart_start
				rateLimit.Ip = host

				nowMs := time.Now().UnixMilli()
				currentSecond := nowMs / 1000

				if currentSecond != timestart_start/1000 {
					elapsedMs := nowMs % 1000
					weight := (1000 - elapsedMs) / 1000

					total := int64(currentCount) + int64(previousCount)*weight

					if total > int64(rateLimit.GetEnvLimit()) {
						w.WriteHeader(http.StatusTooManyRequests)
						w.Write([]byte("Rate limit exceeded"))
						return
					}
					rateLimit.CurrentCount = currentCount + 1
					rateLimit.LastCount = previousCount
					rateLimit.TimestampStart = timestart_start
					rateLimit.Ip = host

					f(w, r)
				}

			}

		}

	}
}
