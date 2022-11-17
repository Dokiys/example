package mock_gomock

//go:generate mockgen -source=./mock.go -destination=./mock_inf.go -package=tdd
type Inf interface {
	Bar() int
	IsGood(name string) bool
	IsOldPerson(p *People) bool
}

type People struct {
	age int
}

type A struct {
	inf Inf
}

func (self *A) IsOne() bool {
	if r := self.inf.Bar(); r == 1 {
		return true
	} else {
		return false
	}
}
