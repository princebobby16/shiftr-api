package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/pbobby001/shiftr/app/controllers"
	"gitlab.com/pbobby001/shiftr/app/controllers/websockets"
	"net/http"
)

//Route Create a single route object
type Route struct {
	Name    string
	Path    string
	Method  string
	Handler http.HandlerFunc
}

//Routes Create an object of different routes
type Routes []Route

// InitRoutes Set up routes
func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	routes := Routes{
		// health check
		Route{
			Name:    "Health Check",
			Path:    "/",
			Method:  http.MethodGet,
			Handler: controllers.HealthCheckHandler,
		},

		Route{
			Name:    "Email Notification Service",
			Path:    "/send-email",
			Method:  http.MethodPost,
			Handler: controllers.EmailNotificationService,
		},

		// websockets
		Route{
			Path:    "/shiftr/postit/subscribers",
			Method:  http.MethodGet,
			Handler: websockets.ScheduleStatus,
		},
	}

	for _, route := range routes {
		router.Name(route.Name).
			Methods(route.Method).
			Path(route.Path).
			Handler(route.Handler)
		//router.Handle("/metrics", promhttp.Handler())
	}

	return router
}
