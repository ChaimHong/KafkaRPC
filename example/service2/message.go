package service2

type Raw []byte
type Cst int8

const (
	C_A Cst = iota
)

type A struct {
	F1 *Raw
}
