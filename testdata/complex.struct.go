package testdata

type ComplextStatus struct {
	Group string  `json:"group"`
	Msg   *string `json:"msg"`
}
type ComplexUser struct {
	Name          *string        `json:"name"`
	Email         string         `json:"email"`
	CurrentStatus ComplextStatus `json:"currentStatus"`
}

type Complex struct {
	Category  uint32        `json:"category"`
	Ips       []string      `json:"ips"`
	Positions []uint8       `json:"positions"`
	User      *ComplexUser  `json:"user"`
	Accounts  []ComplexUser `json:"accounts"`
}
