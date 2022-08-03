package data

import "time"

//go:generate pinfo
type B struct {
	CreateAt time.Time
}

type B2 struct {
	CreateAt time.Time
}

type A struct {
	Id   int
	Name string
	B    B
}

type A2 struct {
	Id   int
	Name string
	B    B2
}
