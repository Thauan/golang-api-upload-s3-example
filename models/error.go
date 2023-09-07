package models

type Error struct {
	Success  bool   `json:"success"`
	Status   int    `json:"status"`
	Messages string `json:"messages"`
}

func (*Error) Build(success bool, status int, messages string) *Error {

	error := &Error{
		Success:  success,
		Status:   status,
		Messages: messages,
	}

	return error
}
