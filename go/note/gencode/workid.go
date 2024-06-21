package gencode

import "sync/atomic"

var atomInt32 atomic.Int32

var _ WorkIdPicker = (*localWorkIdPicker)(nil)

type localWorkIdPicker struct{}

func (l *localWorkIdPicker) PickWorkId() uint32 {
	return uint32(atomInt32.Add(1))
}
