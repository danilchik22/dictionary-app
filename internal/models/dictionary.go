package models

type Dictionary struct {
	ID         int    `json:"id"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type QuerySearch struct {
	Query string `json:"query"`
	IsOne bool   `json:"is_one"`
	Limit int    `json:"limit" env-default:"10"`
}
