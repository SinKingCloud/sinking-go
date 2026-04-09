package node

type SelectNode struct {
	Group           string
	Name            string
	Status          string
	OnlineStatus    string
	Address         string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

type UpdateNode struct {
	Group        interface{}
	Name         interface{}
	Address      interface{}
	OnlineStatus interface{}
	Status       interface{}
	LastHeart    interface{}
}
