package test

import "example.com/example/test/model"

type NumAlias int64

type MainEmbedType int
type MainEmbed struct {
	EmbedName string   `json:"embed_name,omitempty"`
	EmbedNum  NumAlias `json:"embed_num,omitempty"`
}

type SubSt struct {
	SubName string
}

type MainSt struct {
	Name          string           `json:"name"`
	Address       []string         `json:"address,omitempty"`
	Pointer       *string          `json:"pointer,omitempty"`
	TestSubSt     SubSt            `json:"test_sub_st,omitempty"`
	Sibling       []SiblingSt      `json:"sibling,omitempty"`
	Devices       []model.DeviceSt `json:"devices,omitempty"`
	Model         model.ModelSt    `json:"model,omitempty"`
	MainEmbed     `json:"main_embed,omitempty"`
	MainEmbedType `json:"main_embed_type,omitempty"`
	model.ModelEmbed
}
