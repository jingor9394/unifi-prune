package utils

type InputParams struct {
	Model string
	IP    string
	Port  string
	User  string
	Pass  string
}

func NewInputParams() (*InputParams, error) {
	return &InputParams{}, nil
}
