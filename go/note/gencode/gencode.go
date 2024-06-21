package gencode

import (
	"fmt"
	"time"
)

type WorkIdPicker interface {
	PickWorkId() uint32
}

type Generator struct {
	picker   WorkIdPicker
	ring     *Ring
	hashSalt string

	prefix string
}

func NewGenerator(hashSalt string, prefix string, picker WorkIdPicker) *Generator {
	return &Generator{picker: picker, ring: newRing(1000, 9999), hashSalt: hashSalt, prefix: prefix}
}

func (g *Generator) Code() string {
	current, seq := g.ring.Code()
	// sn(16) = 机器码(2) + 时间戳(10) + 序列号(4)
	sn := g.workerId() + g.timestamp(current) + fmt.Sprintf("%d", seq)
	return g.prefix + hash16([]byte(sn+g.hashSalt))
}

// sn(10) = 两位年(2) + 当前时间在今年的秒数(8)
func (g *Generator) timestamp(now time.Time) string {
	// 两位年
	year := now.Year() % 100
	// 当前时间在今年的秒数
	tmStart := time.Date(now.Year(), 0, 0, 0, 0, 0, 0, time.Local)
	seconds := now.Sub(tmStart).Nanoseconds() / 1e9
	return fmt.Sprintf("%d", int64(year*1e8)+seconds)
}

// sn(2) = 两位机器码
func (g *Generator) workerId() string {
	return fmt.Sprintf("%02d", g.picker.PickWorkId()%100)
}

func djbHash64(str []byte) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(str); i++ {
		hash = (hash<<5 + hash) + uint64(str[i])
	}
	return hash
}

func hash16(str []byte) string {
	return fmt.Sprintf("%016d", djbHash64(str))
}
