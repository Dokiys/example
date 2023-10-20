package basic

import "io"

type Generator interface {
	Gen(wr io.Writer) error
}
