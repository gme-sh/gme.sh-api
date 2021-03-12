package shortreq

import (
	"github.com/gofiber/fiber/v2"
	"math"
)

var (
	// GENERIC
	// err
	ResponseErrURLNotFound = &Response{
		InternalCode: -1001,
		StatusCode:   0,
		Message:      "url not found",
	}

	// CREATE
	// ok
	ResponseOkCreate = &Response{
		InternalCode: +2001,
		StatusCode:   201,
		Message:      "created",
	}
	// err
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

	// DELETE
	// ok
	ResponseOkDeleted = &Response{
		InternalCode: +3001,
		StatusCode:   200,
		Message:      "deleted",
	}
	// err
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

	// STATS
	// ok
	ResponseOkStats = &Response{
		InternalCode: +4001,
		StatusCode:   200,
		Message:      "ok",
	}

	// REDIRECT
	// ok
	ResponseOkRedirectDry = &Response{
		InternalCode: +5001,
		StatusCode:   200,
		Message:      "found",
	}
	// err
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

////////////////////////////////////////////////////////////////////////////////////////////////

type Response struct {
	InternalCode int
	StatusCode   int
	Message      string
}

func (r *Response) successable(data interface{}) *Successable {
	return &Successable{
		Success: math.Floor(float64(r.StatusCode)/100) == 2,
		Message: r.Message,
		Code:    r.InternalCode,
		Data:    data,
	}
}

func (r *Response) Send(ctx *fiber.Ctx) error {
	return r.SendWithData(ctx, nil)
}
func (r *Response) SendWithData(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(r.StatusCode).JSON(r.successable(data))
}

func (r *Response) SendWithMessage(ctx *fiber.Ctx, message string) error {
	return r.SendWithMessageData(ctx, message, nil)
}
func (r *Response) SendWithMessageData(ctx *fiber.Ctx, message string, data interface{}) error {
	s := r.successable(data)
	s.Message = message
	return ctx.Status(r.StatusCode).JSON(s)
}
