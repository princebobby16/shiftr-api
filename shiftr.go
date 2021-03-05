package main

import (
	"flag"
	"github.com/gorilla/handlers"
	_ "github.com/joho/godotenv/autoload"
	"gitlab.com/pbobby001/shiftr/app/middlewares"
	"gitlab.com/pbobby001/shiftr/app/router"
	"gitlab.com/pbobby001/shiftr/db"
	"gitlab.com/pbobby001/shiftr/pkg/logs"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Is this better?
	db.Connect()

	r := router.InitRoutes()

	origins := handlers.AllowedOrigins([]string{ /*"https://postit-ui.herokuapp.com"*/ "*"})
	headers := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Content-Length",
		"Content-Event-Type",
		"X-Requested-With",
		"Accept-Encoding",
		"Accept",
		"Authorization",
		"User-Agent",
		"tenant-namespace",
		"trace-id",
		"Access-Control-Allow-Origin",
	})
	methods := handlers.AllowedMethods([]string{
		http.MethodPost,
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPut,
	})

	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "5379"
	}

	server := &http.Server{
		Addr: ":" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CORS(origins, headers, methods)(r), // Pass our instance of gorilla/mux in.
	}

	r.Use(middlewares.JSONMiddleware)

	defer db.Disconnect()
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		// TODO: Fetch port from store
		logs.Logger.Info("Server running on port", port)
		if err := server.ListenAndServe(); err != nil {
			logs.Logger.Warn(err)
		}
	}()

	channel := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.

	signal.Notify(channel, os.Interrupt)
	// Block until we receive our signal.
	<-channel

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = server.Shutdown(ctx)

	// Optionally, you could run server.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logs.Logger.Warn("shutting down")
	os.Exit(0)
}
