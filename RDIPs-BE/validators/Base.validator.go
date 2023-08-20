package validators

type BaseValidator struct {
}
type Validator interface {
	Post([]byte) error
	Put([]byte) error
}

func (*BaseValidator) Post() error {
	return nil
}

func (*BaseValidator) Put() error {
	return nil
}
