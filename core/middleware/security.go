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
				rateLimit.TimestampStart = time.Now().Unix()
				rateLimit.Ip = host

				rateLimit.InsertData()

			} else {
				currentCount := model[0]["current_count"].(string)
				timestart_start := model[0]["timestamp_start"].(string)
				previousCount := model[0]["current_count"].(string)

				intCurrentCount, _ := strconv.ParseInt(currentCount, 10, 64)
				intPreviousCount, _ := strconv.ParseInt(previousCount, 10, 64)
				intTimestart_start, _ := strconv.ParseInt(timestart_start, 10, 64)

				rateLimit.CurrentCount = intCurrentCount
				rateLimit.LastCount = intPreviousCount
				rateLimit.TimestampStart = intTimestart_start
				rateLimit.Ip = host

				nowMs := time.Now().UnixMilli()
				currentSecond := nowMs

				if currentSecond != intTimestart_start {
					elapsedMs := nowMs
					weight := (1000 - elapsedMs)

					total := intCurrentCount + intPreviousCount*weight

					if total > int64(rateLimit.GetEnvLimit()) {
						w.WriteHeader(http.StatusTooManyRequests)
						w.Write([]byte("Rate limit exceeded"))
						return
					}
					rateLimit.CurrentCount = intCurrentCount + 1
					rateLimit.LastCount = intPreviousCount
					rateLimit.TimestampStart = intTimestart_start
					rateLimit.Ip = host

					rateLimit.UpdateData(host)
				}

			}

			f(w, r)

		}

	}
}
