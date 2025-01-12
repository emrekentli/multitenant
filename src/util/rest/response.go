package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type Meta struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Page[T any] struct {
	Size    int   `json:"size"`
	Total   int   `json:"total"`
	Content *[]*T `json:"content"`
}

func jsonResponse(c fiber.Ctx, data any, code string, message string) error {
	return c.JSON(&fiber.Map{
		"data": data,
		"meta": Meta{
			Code:    code,
			Message: message,
		},
	})
}

// wrap data with meta success
func Res(c fiber.Ctx, err error, value any) error {
	if err != nil {
		log.Error(err)
		code, message := Error(err)
		return ErrorRes(c, code, message)
	}

	return jsonResponse(c, value, Success, ErrorCode[Success])
}

// wrap data with meta error
func ErrorRes(c fiber.Ctx, code string, message string) error {
	log.Error(message)
	log.Error(code)
	log.Error(c.Path())
	log.Error(c.Method())
	return jsonResponse(c, nil, code, message)
}

func PageToResponseList[T any, R any](page *Page[T], convertFunc func(*T) *R) *Page[R] {
	if page == nil {
		return nil
	}

	var resContent []*R
	for _, item := range *page.Content {
		resContent = append(resContent, convertFunc(item))
	}

	return &Page[R]{
		Size:    page.Size,
		Total:   page.Total,
		Content: &resContent,
	}
}
func ListToResponseList[T any, R any](list []T, convertFunc func(T) R) []R {
	if list == nil {
		return []R{}
	}

	var resContent []R
	for _, item := range list {
		resContent = append(resContent, convertFunc(item))
	}

	return resContent
}
