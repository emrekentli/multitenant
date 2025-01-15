package role

func requestToModal(request *ModalRequest) *Modal {
	return &Modal{
		Name:        request.Name,
		Description: request.Description,
	}
}

func modalToResponse(modal *Modal) *ModalResponse {
	return &ModalResponse{
		Id:          modal.Id,
		Name:        modal.Name,
		Description: modal.Description,
	}
}
