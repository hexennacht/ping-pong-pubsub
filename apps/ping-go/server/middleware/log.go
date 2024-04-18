package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func NewZapLoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	if logger == nil {
		logger = zap.Must(zap.NewProduction())
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			now := time.Now()

			var buffer bytes.Buffer

			wrw.Tee(&buffer)

			defer func() {
				reqLogger := logger.With(
					zap.String("proto", r.Proto),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("reqId", middleware.GetReqID(r.Context())),
					zap.Duration("lat", time.Since(now)),
					zap.Int("status", wrw.Status()),
					zap.Int("size", wrw.BytesWritten()),
					zap.String("body", buffer.String()),
				)
				ref := wrw.Header().Get("Referer")
				if ref == "" {
					ref = r.Header.Get("Referer")
				}

				if ref != "" {
					reqLogger = reqLogger.With(zap.String("ref", ref))
				}

				ua := wrw.Header().Get("User-Agent")
				if ua == "" {
					ua = r.Header.Get("User-Agent")
				}

				if ua != "" {
					reqLogger = reqLogger.With(zap.String("ua", ua))
				}

				reqLogger.Info("Served")
			}()

			next.ServeHTTP(wrw, r)
		})
	}
}
