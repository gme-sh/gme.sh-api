package shortreq

// OK
var (
	ResponseOkCreate = &Response{
		InternalCode: +2001,
		StatusCode:   201,
		Message:      "created",
	}
)

// ERR
var (
	ResponseErrInvalidURL = &Response{
		InternalCode: -2001,
		StatusCode:   400,
		Message:      "invalid url",
	}
	ResponseErrAliasOccupied = &Response{
		InternalCode: -2002,
		StatusCode:   409,
		Message:      "alias is already occupied",
	}
	ResponseErrDomainBlocked = &Response{
		InternalCode: -2003,
		StatusCode:   403,
		Message:      "domain is blocked",
	}
	ResponseErrGeneratedAliasNotAvailable = &Response{
		InternalCode: -2004,
		StatusCode:   500,
		Message:      "no generated alias available",
	}
	ResponseErrDatabaseSave = &Response{
		InternalCode: -2005,
		StatusCode:   503,
		Message:      "error saving",
	}
	ResponseErrInvalidID = &Response{
		InternalCode: -2006,
		StatusCode:   400,
		Message:      "invalid id (alias)",
	}
)
