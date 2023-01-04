package installer

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kraisitdev/assessment/app/rest/handler"
)

func SetupLogging() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Info().Msg("Config Logging Success.")
}

func SetupMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("method", v.Method).
				Str("latency_human", v.Latency.String()).
				Msg("")

			return nil
		},
	}))

	e.Use(middleware.Recover())

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "November",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "10, 2009", nil
		},
	}))
}

func SetupEndPoint(e *echo.Echo) {
	h := handler.NewApp(true)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/expenses", h.InsertExpense)

	e.GET("/expenses/:id", h.GetExpenseById)
}

func SetupServer(e *echo.Echo) {
	// Graceful shutdown
	go func() {
		// Start server
		log.Info().Msg("Starting Server...")
		err := e.Start(":" + os.Getenv("PORT"))
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Start Server Error: %s", err.Error())
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Shutdown Server Error: %s", err.Error())
	}
	log.Info().Msg("Shutdown Server Success.")
}
