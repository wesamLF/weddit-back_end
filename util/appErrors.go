package util

type AppError struct {
	Message    string `json:"message"`     // what client sees
	StatusCode int    `json:"status_code"` // HTTP code
	Err        error  `json:"-"`           // internal (not sent to client)
}
