package constants

const (
	StatusSuccess      = "00"
	StatusBadRequest   = "400"
	StatusUnauthorized = "401"
	StatusConflict     = "409"
	StatusFailed       = "05"
)

const (
	MessageSuccess      = "Success"
	MessageBadRequest   = "Invalid Request Format"
	MessageUnauthorized = "Unauthorized"
	MessageConflict     = "Conflict"
	MessageFailed       = "Something went wrong"
)

type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}
