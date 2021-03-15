package shortreq

// OK
var (
	ResponseOkDeleted = &Response{
		InternalCode: +3001,
		StatusCode:   200,
		Message:      "deleted",
	}
)

// ERR
var (
	ResponseErrEmptyID = &Response{
		InternalCode: -3006,
		StatusCode:   400,
		Message:      "empty short-id",
	}
	ResponseErrLocked = &Response{
		InternalCode: -3007,
		StatusCode:   403,
		Message:      "url is locked",
	}
	ResponseErrSecretMismatch = &Response{
		InternalCode: -3008,
		StatusCode:   401,
		Message:      "secret mismatch",
	}
)
