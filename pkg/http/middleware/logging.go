package xmiddleware

import (
	"time"

	"github.com/labstack/echo/v4"

	xlogger "thomas.vn/hr_recruitment/pkg/logger"
)

func RequestLogging(logger *xlogger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			err := next(c)
			stop := time.Now()

			fields := []xlogger.Field{
				xlogger.String("host", req.Host),
				xlogger.String("method", req.Method),
				xlogger.String("uri", req.RequestURI),
				xlogger.String("user_agent", req.UserAgent()),
				xlogger.Int("status", res.Status),
				xlogger.String("id", req.Header.Get(echo.HeaderXRequestID)),
				xlogger.Int64("latency", int64(stop.Sub(start))),
				xlogger.String("latency_human", stop.Sub(start).String()),
				xlogger.Int64("bytes_in", req.ContentLength),
				xlogger.Int64("bytes_out", res.Size),
			}

			if err != nil {
				fields = append(fields, xlogger.Error(err))
				logger.Error("request failed", fields...)
			} else {
				logger.Info("request completed", fields...)
			}

			return err
		}
	}
}
