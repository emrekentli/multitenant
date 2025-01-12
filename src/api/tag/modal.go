package tag

import (
	"time"
)

type Modal struct {
	Id       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Name     string    `json:"name"`
}

type ModalResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ModalRequest struct {
	Name string `json:"name"`
}
type ModalDeleteRequest struct {
	IdList []int64 `json:"idList"`
}
