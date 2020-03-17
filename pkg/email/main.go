package email

type Sender interface {
	Send(Mail) error
}

type Mail struct {
	To      string
	From    string
	Subject string
	Body    string
}
