// +heroku install .
module gitlab.com/pbobby001/shiftr

// +heroku goVersion go1.15
go 1.15

require (
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/goware/emailx v0.2.0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.9.0
	github.com/myesui/uuid v1.0.0 // indirect
	github.com/twinj/uuid v1.0.0
	go.uber.org/fx v1.13.1
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	gopkg.in/stretchr/testify.v1 v1.2.2 // indirect
)
