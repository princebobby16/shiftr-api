package websockets

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/shiftr/db"
	"gitlab.com/pbobby001/shiftr/pkg"
	"gitlab.com/pbobby001/shiftr/pkg/logs"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func Writer(conn *websocket.Conn, connection *sql.DB) {
	for {
		ticker := time.NewTicker(2 * time.Second)

		for t := range ticker.C {
			logs.Logger.Info("Updating status: ", t)
			subscribers, err := FetchStatuses(connection)
			if err != nil {
				_ = logs.Logger.Error(err)
				return
			}

			jsonBytes, err := json.Marshal(subscribers)
			if err != nil {
				_ = logs.Logger.Error(err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, jsonBytes)
			if err != nil {
				_ = logs.Logger.Error(err)
				return
			}

		}
	}
}

func FetchStatuses(connection *sql.DB) (*pkg.PostItSubscribers, error) {

	//  Prepare the query
	query := fmt.Sprintf("SELECT * FROM shiftr.postit_subscribers")

	// run the query
	rows, err := connection.Query(query)
	if err != nil {
		return &pkg.PostItSubscribers{}, err
	}

	var subscribers []pkg.Subscriber
	for rows.Next() {
		var subscriber pkg.Subscriber
		err = rows.Scan(
			&subscriber.Id,
			&subscriber.Email,
			&subscriber.PhoneNumber,
			&subscriber.EmailSentStatus,
			&subscriber.CreateAt,
			&subscriber.UpdatedAt,
		)
		if err != nil {
			return &pkg.PostItSubscribers{}, err
		}
		subscribers = append(subscribers, subscriber)
	}

	postitSubscribers := &pkg.PostItSubscribers{
		SubId:       uuid.NewV4().String(),
		SubCount:    len(subscribers),
		Subscribers: subscribers,
	}

	return postitSubscribers, nil
}

func ScheduleStatus(w http.ResponseWriter, r *http.Request) {

	logs.Logger.Info("connecting to websocket")

	ws, err := upgrade(w, r)
	if err != nil {
		_ = logs.Logger.Error(err)
		return
	}

	logs.Logger.Info("connection upgraded")
	go Writer(ws, db.Connection)
}
