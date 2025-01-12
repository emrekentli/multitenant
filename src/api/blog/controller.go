package blog

import (
	"github.com/emrekentli/multitenant-boilerplate/src/util/rest"
	"github.com/gofiber/fiber/v3"
)

func getAll(c fiber.Ctx) error {
	limit, offset, order := rest.GetPageParamsWithSort(c)
	res, err := GetAll(limit, offset, order)
	return rest.Res(c, err, res)
}

func create(c fiber.Ctx) error {
	var modal ModalRequest
	err := rest.SetBodyAndValidate(c, &modal)
	if err != nil {
		return rest.Res(c, err, nil)
	}
	res, err := Create(requestToModal(&modal))
	return rest.Res(c, err, ModalToResponse(res))
}

func update(c fiber.Ctx) error {
	id := c.Params("id")
	var modal ModalRequest
	err := rest.SetBodyAndValidate(c, &modal)
	if err != nil {
		return rest.Res(c, err, nil)
	}
	err = Update(id, requestToModal(&modal))
	return rest.Res(c, err, nil)
}

func deleteByIds(c fiber.Ctx) error {
	var modalDeleteRequest ModalDeleteRequest
	err := rest.SetBodyAndValidate(c, &modalDeleteRequest)
	if err != nil {
		return rest.Res(c, err, nil)
	}

	err = Delete(&modalDeleteRequest)
	return rest.Res(c, err, nil)
}

func Register(router fiber.Router) {
	group := router.Group("/blog")
	group.Get("/", getAll)
	group.Post("/", create)
	group.Put("/:id", update)
	group.Delete("/", deleteByIds)
}
