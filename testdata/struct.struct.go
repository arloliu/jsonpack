package testdata

type PhoneNumber struct {
	Area   uint32 `json:"area"`
	Number string `json:"number"`
}

type Address struct {
	Area    uint32 `json:"area"`
	Address string `json:"address"`
}

type Skill struct {
	Languages []string `json:"languages"`
}

type MapInfo struct {
	MapName  string `json:"mapName"`
	Location string `json:"location"`
}

type TestStruct struct {
	Name      string        `json:"name"`
	Sex       uint32        `json:"sex"`
	Nicknames []string      `json:"nicknames"`
	Maps      *MapInfo      `json:"maps"`
	Phones    []PhoneNumber `json:"phones"`
	Address   Address       `json:"address"`
	Skill
}
