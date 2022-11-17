package mock_gomock

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInf := NewMockInf(ctrl)
	mockInf.EXPECT().Bar().Return(1).AnyTimes()
	mockInf.EXPECT().IsGood(gomock.Eq("zhangsan")).Return(true).AnyTimes()
	mockInf.EXPECT().IsOldPerson(newGtAgeMatcher(10)).Return(true).AnyTimes()
	a := &A{inf: mockInf}

	assert.Equal(t, true, a.IsOne())
}

type gtAgeMatcher struct {
	age int
}

func newGtAgeMatcher(age int) *gtAgeMatcher {
	return &gtAgeMatcher{age: age}
}

func (self *gtAgeMatcher) Matches(x interface{}) bool {
	p, ok := x.(*People)
	if !ok {
		return false
	}
	if self.age < p.age {
		return true
	}

	return false
}

func (self *gtAgeMatcher) String() string {
	return fmt.Sprintf("age > %d", self.age)
}
