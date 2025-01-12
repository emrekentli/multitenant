package user

import (
	"github.com/emrekentli/multitenant-boilerplate/src/util/rest"
	"github.com/gofiber/fiber/v3"
)

func getAll(c fiber.Ctx) error {
	limit, offset := rest.GetPageParams(c)
	res, err := GetAll(limit, offset)
	return rest.Res(c, err, rest.PageToResponseList(res, modalToResponse))
}

func login(c fiber.Ctx) error {
	var modalLoginRequest ModalRequest
	err := rest.SetBodyAndValidate(c, &modalLoginRequest)
	if err != nil {
		return rest.Res(c, err, nil)
	}

	res, err := Login(&modalLoginRequest)
	return rest.Res(c, err, res)
}

func create(c fiber.Ctx) error {
	var modal ModalRequest
	err := rest.SetBodyAndValidate(c, &modal)
	if err != nil {
		return rest.Res(c, err, nil)
	}

	res, err := Create(requestToModal(&modal))
	return rest.Res(c, err, modalToResponse(res))
}

func update(c fiber.Ctx) error {
	var modal Modal
	err := rest.SetBodyAndValidate(c, &modal)
	if err != nil {
		return rest.Res(c, err, nil)
	}

	id := c.Params("id")
	err = Update(id, &modal)
	return rest.Res(c, err, &modal)
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
	group := router.Group("/user")
	group.Get("/", getAll)
	group.Post("/login", login)

	group.Post("/", create)
	group.Put("/:id", update)
	group.Delete("/", deleteByIds)
}
