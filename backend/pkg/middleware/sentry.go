package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

// SentryRecovery catches panics and reports them to Sentry if initialized
func SentryRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				err, ok := rcv.(error)
				if !ok {
					err = fmt.Errorf("%v", rcv)
				}

				hub := sentry.CurrentHub().Clone()
				hub.Scope().SetRequest(r)
				hub.CaptureException(err)

				log.Printf("[PANIC RECOVERED] %v", rcv)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
