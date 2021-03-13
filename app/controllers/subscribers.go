package controllers

import (
	"encoding/json"
	"github.com/goware/emailx"
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
		_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
			Data: pkg.Data{
				Id:        uuid.NewV4().String(),
				UiMessage: "something went wrong. contact admin!",
			},
			Meta: pkg.Meta{
				Timestamp:     time.Now(),
				TransactionId: uuid.NewV4().String(),
				Status:        "SUCCESS",
			},
		})
		return
	}

	logs.Logger.Info(string(body))

	err = json.Unmarshal(body, &subscriber)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
			Data: pkg.Data{
				Id:        uuid.NewV4().String(),
				UiMessage: "something went wrong. contact admin!",
			},
			Meta: pkg.Meta{
				Timestamp:     time.Now(),
				TransactionId: uuid.NewV4().String(),
				Status:        "SUCCESS",
			},
		})
		return
	}

	logs.Logger.Info(subscriber)

	err = emailx.Validate(subscriber.Email)
	if err != nil {
		logs.Logger.Error("Email is not valid.")

		if err == emailx.ErrInvalidFormat {
			logs.Logger.Error("Wrong format.")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
				Data: pkg.Data{
					Id:        uuid.NewV4().String(),
					UiMessage: "Invalid email format",
				},
				Meta: pkg.Meta{
					Timestamp:     time.Now(),
					TransactionId: uuid.NewV4().String(),
					Status:        "FAIL",
				},
			})
			return
		}

		if err == emailx.ErrUnresolvableHost {
			logs.Logger.Error("Unresolvable host.")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
				Data: pkg.Data{
					Id:        uuid.NewV4().String(),
					UiMessage: "invalid email host",
				},
				Meta: pkg.Meta{
					Timestamp:     time.Now(),
					TransactionId: uuid.NewV4().String(),
					Status:        "FAIL",
				},
			})
			return
		}
	}

	query := `INSERT INTO shiftr.postit_subscribers(subscriber_id, subscriber_email, subscriber_phone_number, email_sent_status) VALUES ($1, $2, $3, $4);`

	subscriber.Id = uuid.NewV4()
	logs.Logger.Info("Subscriber ID: ", subscriber.Id)

	_, err = db.Connection.Exec(query, subscriber.Id, subscriber.Email, subscriber.PhoneNumber, false)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
			Data: pkg.Data{
				Id:        uuid.NewV4().String(),
				UiMessage: "email or phone number already taken.",
			},
			Meta: pkg.Meta{
				Timestamp:     time.Now(),
				TransactionId: uuid.NewV4().String(),
				Status:        "SUCCESS",
			},
		})
		return
	}

	logs.Logger.Info("Email is valid")
	retry, err := pkg.SendEmail([]string{subscriber.Email})
	if err != nil {
		logs.Logger.Error(err)
		if retry {
			_, err = pkg.SendEmail([]string{subscriber.Email})
			logs.Logger.Error(err)
		}
		w.WriteHeader(http.StatusCreated)
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
		return
	}

	logs.Logger.Info("Email Sent")

	query = `UPDATE shiftr.postit_subscribers SET email_sent_status = $1 WHERE subscriber_id = $2`
	_, err = db.Connection.Exec(query, true, subscriber.Id)
	if err != nil {
		logs.Logger.Info(err)
	}

	w.WriteHeader(http.StatusCreated)
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