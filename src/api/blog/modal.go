package blog

import (
	"github.com/emrekentli/multitenant-boilerplate/src/api/tag"
	"time"
)

type Modal struct {
	Id       int64       `json:"id"`
	Created  time.Time   `json:"created"`
	Modified time.Time   `json:"modified"`
	Body     string      `json:"body"`
	Image    string      `json:"image"`
	Tags     []tag.Modal `json:"tags"`
	Slug     string      `json:"slug"`
}

type ModalResponse struct {
	Id    int64                `json:"id"`
	Body  string               `json:"body"`
	Image string               `json:"image"`
	Tags  []*tag.ModalResponse `json:"tags"`
	Slug  string               `json:"slug"`
}

type ModalRequest struct {
	Body  string      `json:"body"`
	Image string      `json:"image"`
	Tags  []tag.Modal `json:"tags"`
	Slug  string      `json:"slug"`
}
type ModalDeleteRequest struct {
	IdList []int64 `json:"idList"`
}
