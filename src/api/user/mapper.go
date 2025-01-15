package user

import "app/src/general/util/hash"

func requestToModal(request *ModalRequest) *Modal {
	modal := &Modal{
		Email:    request.Email,
		Password: hash.Hash(request.Password),
	}
	return modal
}

func modalToResponse(modal *Modal) *ModalResponse {
	return &ModalResponse{
		Id:    modal.Id,
		Email: modal.Email,
	}
}
