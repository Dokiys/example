package localcache

import (
	"testing"
)

func TestName(t *testing.T) {
	type S [2]any
	s := S{"1", 2}
	s2 := S{"1", 2}
	if s == s2 {
		t.Log("==")
	}
}

var _ IModel = (*MockIModel)(nil)

type MockIModel struct{}

func (m MockIModel) TableName() string {
	// TODO implement me
	panic("implement me")
}
func (m MockIModel) PK() string {
	// TODO implement me
	panic("implement me")
}

func TestNew(t *testing.T) {
	_ = NewLocalCache[[2]any, MockIModel](nil, nil, PkDbTags("", ""))
}
