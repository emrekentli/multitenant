package role

import (
	"app/app/middlewares/tenantcontext"
	"app/src/general/util/rest"
	"github.com/gofiber/fiber/v3"
)

func getAll(c fiber.Ctx) error {
	limit, offset := rest.GetPageParams(c)
	schemaName := tenantcontext.GetTenantSchemaName(c)
	res, err := GetAll(schemaName, limit, offset)
	return rest.Res(c, err, rest.ListToResponseList(res, modalToResponse))
}

func create(c fiber.Ctx) error {
	var modal ModalRequest
	err := rest.SetBodyAndValidate(c, &modal)
	if err != nil {
		return rest.Res(c, err, nil)
	}
	schemaName := tenantcontext.GetTenantSchemaName(c)
	res, err := Create(schemaName, requestToModal(&modal))
	return rest.Res(c, err, modalToResponse(res))
}

func update(c fiber.Ctx) error {
	id := c.Params("id")
	var modal ModalRequest
	err := rest.SetBodyAndValidate(c, &modal)
	if err != nil {
		return rest.Res(c, err, nil)
	}
	schemaName := tenantcontext.GetTenantSchemaName(c)
	err = Update(schemaName, id, requestToModal(&modal))
	return rest.Res(c, err, nil)
}

func deleteByIds(c fiber.Ctx) error {
	var modalDeleteRequest ModalDeleteRequest
	err := rest.SetBodyAndValidate(c, &modalDeleteRequest)
	if err != nil {
		return rest.Res(c, err, nil)
	}
	schemaName := tenantcontext.GetTenantSchemaName(c)
	err = Delete(schemaName, &modalDeleteRequest)
	return rest.Res(c, err, nil)
}

func Register(router fiber.Router) {
	group := router.Group("/roles")
	group.Get("/", getAll)
	group.Post("/", create)
	group.Put("/:id", update)
	group.Delete("/", deleteByIds)
}
