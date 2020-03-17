package backends

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"photo/pkg/email"
)

type SendGrid struct {
	apiKey string
	url    string
}

// NewSendGrid instantiate new Sender with SendGrid backend
//
// Example:
// sender := backends.NewSendGrid("YOUR_API_KEY")
// mail := email.Mail{
//   To:      "rokkerruslan@protonmail.com",
//   From:    "noreply@photogram.live",
//   Subject: "Code",
//   Body:    "123",
// }
// sender.Send(mail)
//
func NewSendGrid(apiKey string) *SendGrid {
	return &SendGrid{
		apiKey: apiKey,
		url:    "https://api.sendgrid.com/v3/mail/send",
	}
}

func (s *SendGrid) Send(mail email.Mail) error {
	baseErr := "SendGrid.Send fails: %w"

	data := request{
		Personalizations: []per{{
			Subject: mail.Subject,
			To: []to{{
				Email: mail.To,
			}},
		}},
		From: from{
			Email: mail.From,
		},
		Content: []content{{
			Type:  "text/plain",
			Value: mail.Body,
		}},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}

	req, err := http.NewRequest("POST", s.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}

	fmt.Println(res)
	fmt.Println(string(body))

	return nil
}

// Internal structures

type to struct {
	Email string `json:"email"`
}

type from struct {
	Email string `json:"email"`
}

type per struct {
	Subject string `json:"subject"`
	To      []to   `json:"to"`
}

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type request struct {
	Personalizations []per `json:"personalizations"`
	From             from  `json:"from"`
	Content          []content
}
