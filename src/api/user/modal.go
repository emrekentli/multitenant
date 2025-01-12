package user

import "time"

type Modal struct {
	Id       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type ModalResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

type ModalRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type ModalDeleteRequest struct {
	IdList []int64 `json:"idList" validate:"required"`
}

type JwtResponse struct {
	Token string `json:"token"`
}
