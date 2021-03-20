package shortreq

// OK
var (
	ResponseOkPoolGet = &Response{
		InternalCode: +6001,
		StatusCode:   200,
		Message:      "",
	}
	ResponseOkPoolUpdating = &Response{
		InternalCode: +6002,
		StatusCode:   200,
		Message:      "updated",
	}
)

// ERR
var (
	ResponseErrPoolNotFound = &Response{
		InternalCode: -6001,
		StatusCode:   404,
		Message:      "pool not found",
	}
	ResponseErrPoolSecretMismatch = &Response{
		InternalCode: -6002,
		StatusCode:   400,
		Message:      "secret mismatch",
	}
	ResponseErrPoolUpdating = &Response{
		InternalCode: -6003,
		StatusCode:   400,
		Message:      "error updating",
	}
	ResponseErrPoolInvalidURL = &Response{
		InternalCode: -6004,
		StatusCode:   400,
		Message:      "invalid url",
	}
)
