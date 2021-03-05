package pkg

import (
	"github.com/twinj/uuid"
	"time"
)

type (

	Subscriber struct {
		Id uuid.UUID 		`json:"id"`
		Email string 		`json:"email"`
		PhoneNumber string 	`json:"phone_number"`
		CreateAt time.Time 	`json:"create_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PostItSubscribers struct {
		SubId string 				`json:"sub_id"`
		SubCount int 				`json:"sub_count"`
		Subscribers []Subscriber	`json:"subscribers"`
	}

	EmailRequest struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Message     string `json:"message"`
	}

	StandardResponse struct {
		Data Data `json:"data"`
		Meta Meta `json:"meta"`
	}

	Data struct {
		Id        string `json:"id"`
		UiMessage string `json:"ui_message"`
	}

	Meta struct {
		Timestamp     time.Time `json:"timestamp"`
		TransactionId string    `json:"transaction_id"`
		TraceId       string    `json:"trace_id"`
		Status        string    `json:"status"`
	}
)
