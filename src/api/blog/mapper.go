package blog

import (
	"github.com/emrekentli/multitenant-boilerplate/src/api/tag"
	"github.com/emrekentli/multitenant-boilerplate/src/util/rest"
)

func requestToModal(request *ModalRequest) *Modal {
	return &Modal{
		Body:  request.Body,
		Image: request.Image,
		Tags:  request.Tags,
	}
}

func ModalToResponse(modal *Modal) *ModalResponse {
	return &ModalResponse{
		Id:    modal.Id,
		Body:  modal.Body,
		Image: modal.Image,
		Tags:  rest.ListToResponseList(modal.Tags, tag.ModalToResponse),
	}
}
