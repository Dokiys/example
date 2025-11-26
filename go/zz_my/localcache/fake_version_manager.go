package localcache

import (
	"context"
	"fmt"
	"time"
)

var _ VersionManager = (*FakeVersionManager)(nil)

type FakeVersionManager struct{}

func NewFakeVersionManager() *FakeVersionManager {
	return &FakeVersionManager{}
}

func (r *FakeVersionManager) Version(_ context.Context, _ string) (string, error) {
	return fmt.Sprint(time.Now().UnixMicro()), nil
}
func (r *FakeVersionManager) IncrVersion(_ context.Context, _ string) (string, error) {
	return fmt.Sprint(time.Now().UnixMicro()), nil
}
