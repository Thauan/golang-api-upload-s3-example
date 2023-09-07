package models

type Response struct {
	Success  bool   `json:"success"`
	Status   int    `json:"status"`
	Data     File   `json:"data, omitempty"`
	Messages string `json:"messages"`
}

func (*Response) Build(success bool, status int, data File, messages string) *Response {

	error := &Response{
		Success:  success,
		Data:     data,
		Status:   status,
		Messages: messages,
	}

	return error
}
