package tag

func requestToModal(request ModalRequest) *Modal {
	return &Modal{
		Name: request.Name,
	}
}

func ModalToResponse(modal Modal) *ModalResponse {
	return &ModalResponse{
		Id:   modal.Id,
		Name: modal.Name,
	}
}
