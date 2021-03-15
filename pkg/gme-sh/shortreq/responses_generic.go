package shortreq

// OK
var (
	ResponseOkRedirectDry = &Response{
		InternalCode: +5001,
		StatusCode:   200,
		Message:      "found",
	}
)

// ERR
var (
	ResponseErrURLNotFound = &Response{
		InternalCode: -1001,
		StatusCode:   0,
		Message:      "url not found",
	}
)
