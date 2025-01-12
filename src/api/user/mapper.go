package user

func requestToModal(request *ModalRequest) *Modal {
	return &Modal{
		Email:    request.Email,
		Password: request.Password,
	}
}

func modalToResponse(modal *Modal) *ModalResponse {
	return &ModalResponse{
		Id:    modal.Id,
		Email: modal.Email,
	}
}
