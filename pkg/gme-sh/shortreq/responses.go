package shortreq

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	InternalCode int
	StatusCode   int
	Message      string
}

type Successable struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) successable(data interface{}) *Successable {
	return &Successable{
		Success: r.StatusCode >= 200 && r.StatusCode < 400,
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
