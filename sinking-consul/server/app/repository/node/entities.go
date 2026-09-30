package node

type SelectNode struct {
	Group           *string
	Name            *string
	Status          *int
	OnlineStatus    *int
	Address         *string
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}

type UpdateNode struct {
	Group        *string
	Name         *string
	Address      *string
	OnlineStatus *int
	Status       *int
	LastHeart    *int64
}
