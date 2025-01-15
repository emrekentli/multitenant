package role

import "time"

type Modal struct {
	Id          int64     `json:"id"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type ModalResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModalRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ModalDeleteRequest struct {
	IdList []int64 `json:"idList" validate:"required"`
}
