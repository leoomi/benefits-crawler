package models

type Benefits struct {
	CPF      string   `json:"cpf"`
	Benefits []string `json:"benefits"`
}
