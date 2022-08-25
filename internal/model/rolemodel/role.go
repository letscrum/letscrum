package rolemodel

import "github.com/letscrum/letscrum/internal/model"

type Role struct {
	model.Model

	Name string `json:"name"`
}
