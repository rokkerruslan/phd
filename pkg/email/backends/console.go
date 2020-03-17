package backends

type ConsoleSender struct {

}

func NewConsoleSender() *ConsoleSender {
	s := ConsoleSender{}
	return &s
}

func (s *ConsoleSender) Send() error {
	panic("implement me")
}
