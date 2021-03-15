package shortreq

// OK
var (
	ResponseOkStats = &Response{
		InternalCode: +4001,
		StatusCode:   200,
		Message:      "ok",
	}
)

// ERR
var (
	ResponseErrRequestedFile = &Response{
		InternalCode: -5001,
		StatusCode:   404,
		Message:      "requested file",
	}
	ResponseErrExpired = &Response{
		InternalCode: -5002,
		StatusCode:   410,
		Message:      "expired",
	}
)
