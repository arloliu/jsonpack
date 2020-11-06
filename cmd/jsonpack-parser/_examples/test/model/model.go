package model // import "example.com/example/test/model"

type DeviceSt struct {
	Id       string `json:"id,omitempty"`
	Location string `json:"location,omitempty"`
}

type ModelSt struct {
	Model string `json:"model,omitempty"`
}

type ModelEmbed struct {
	ModelName string `json:"model_name,omitempty"`
	ModelNum  int32  `json:"model_num,omitempty"`
}
