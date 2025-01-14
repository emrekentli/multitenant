package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

const MaxLimit = 100

func GetPageParams(c fiber.Ctx) (int, int) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		offset = 0
	}
	return limit, offset
}

func GetPageParamsWithSort(c fiber.Ctx) (int, int, string) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		offset = 0
	}

	order := c.Query("order", "desc")
	return limit, offset, order
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func validateRequest(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func SetBodyAndValidate(c fiber.Ctx, body interface{}) error {
	err := c.Bind().Body(body)
	if err != nil {
		return err
	}
	err = validateRequest(body)
	if err != nil {
		return err
	}
	return nil
}
