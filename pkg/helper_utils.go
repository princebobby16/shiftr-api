package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/shiftr/pkg/logs"
	"html/template"
	"net/http"
	"net/smtp"
	"path/filepath"
	"time"
)

/* Validate header is a function used to make sure that the required  headers are sent to the API
It takes the http request and extracts the headers from it and returns a map of the needed headers
and an error. Other headers are essentially ignored.*/
func ValidateHeaders(r *http.Request) (map[string]string, error) {
	//Group the headers
	receivedHeaders := make(map[string]string)
	requiredHeaders := []string{"trace-id", "tenant-namespace"}

	for _, header := range requiredHeaders {
		value := r.Header.Get(header)
		if value != "" {
			receivedHeaders[header] = value
		} else if value == "" {
			return nil, errors.New("Required header: " + header + " not found")
		} else {
			return nil, errors.New("No headers received be sure to send some headers")
		}
	}

	return receivedHeaders, nil
}

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

/* Helper function to send the email to shiftrgh@gmail.com */
func SendEmail() (bool, error) {
	from := "princebobby506@gmail.com"
	password := "yoforreal.com"

	to := []string{
		"shiftrgh@gmail.com",
	}
	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	emailBody, err := parseTemplate("index.html")
	if err != nil {
		return false, err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "THANK YOU FOR SUBSCRIBING TO POSTIT" + "!\n"
	headers := []byte(subject + mime + "\n" + string(emailBody))
	var body bytes.Buffer

	body.Write(headers)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	retry := false
	// Sending email.
	err = smtp.SendMail(smtpServer.Address(), auth, from, to, body.Bytes())
	if err != nil {
		retry = true
		return retry, err
	}

	return retry, err
}

func parseTemplate(s string) ([]byte, error) {

	path, err := filepath.Abs(fmt.Sprintf("pkg/12/%s", s))
	if err != nil {
		return nil, err
	}
	logs.Logger.Info(path)

	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, nil
	}
	logs.Logger.Info(t.Name())

	buff := new(bytes.Buffer)
	err = t.Execute(buff, nil)
	if err != nil {
		return nil, nil
	}

	return buff.Bytes(), nil
}
/* Helper func to handle error */
func SendErrorResponse(w http.ResponseWriter, tId uuid.UUID, traceId string, err error, httpStatus int) {
	w.WriteHeader(httpStatus)
	logs.Logger.Error(err)
	_ = json.NewEncoder(w).Encode(StandardResponse {
		Data: Data{
			Id:        "",
			UiMessage: "Something went wrong! Contact Admin!",
		},
		Meta: Meta{
			Timestamp:     time.Now(),
			TransactionId: tId.String(),
			TraceId:       traceId,
			Status:        "FAILED",
		},
	})
	return
}