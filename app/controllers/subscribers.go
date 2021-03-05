package controllers

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/shiftr/db"
	"gitlab.com/pbobby001/shiftr/pkg"
	"gitlab.com/pbobby001/shiftr/pkg/logs"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetSubscriber(w http.ResponseWriter, r *http.Request) {
	var subscriber pkg.Subscriber

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logs.Logger.Info(string(body))

	err = json.Unmarshal(body, &subscriber)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logs.Logger.Info(subscriber)

	query := `INSERT INTO shiftr.postit_subscribers(subscriber_id, subscriber_email, subscriber_phone_number) VALUES ($1, $2, $3);`

	_, err = db.Connection.Exec(query, uuid.NewV4().String(), subscriber.Email, subscriber.PhoneNumber)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
		Data: pkg.Data{
			Id:        uuid.NewV4().String(),
			UiMessage: "subscription received!, We will contact you when POSTIT is ready!",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: uuid.NewV4().String(),
			Status:        "SUCCESS",
		},
	})
}