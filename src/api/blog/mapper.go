package blog

import (
	"app/src/api/tag"
	"app/src/general/util/rest"
)

func requestToModal(request *ModalRequest) *Modal {
	return &Modal{
		Body:  request.Body,
		Image: request.Image,
		Slug:  request.Slug,
		Tags:  request.Tags,
	}
}

func ModalToResponse(modal *Modal) *ModalResponse {
	return &ModalResponse{
		Id:    modal.Id,
		Body:  modal.Body,
		Image: modal.Image,
		Slug:  modal.Slug,
		Tags:  rest.ListToResponseList(modal.Tags, tag.ModalToResponse),
	}
}
