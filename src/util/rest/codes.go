package rest

const (
	Success             = "200"
	BadRequest          = "400"
	ServerInternalError = "500"
	NotFound            = "404"
	AlreadyExists       = "409"
	Unauthorized        = "401"
)

var ErrorCode = map[string]string{
	Success:             "Success",
	BadRequest:          "Bad Request",
	NotFound:            "Not Found",
	AlreadyExists:       "Already Exists",
	ServerInternalError: "Server Internal Error",
	Unauthorized:        "Unauthorized",
}

func Error(error error) (string, string) {
	err := ErrorCode[error.Error()]
	if err == "" {
		return ServerInternalError, ErrorCode[ServerInternalError]
	}
	return error.Error(), err
}
