package testdata

type TestArrayStructObj struct {
	Name string `json:"name"`
}

type TestArrayStruct struct {
	Num1   int32              `json:"num1"`
	Num2   float64            `json:"num2"`
	Object TestArrayStructObj `json:"obj"`
}
