package service1

import "fmt"

type ServiceA struct{}

type AIn struct {
	Time int64
	V    int
}
type AOut struct {
	V    string
	Time int64
}

const ServiceA_A = "ServiceA.A"
const ServiceA_B = "ServiceA.B"

func (a *ServiceA) A(args *AIn, reply *AOut) error {
	args.V++
	reply.Time = args.Time
	reply.V = fmt.Sprintf("ServiceA.A %d", args.V)
	return nil
}

type BIn struct {
	B int
}
type BOut struct {
	B string
}

func (a *ServiceA) B(args *BIn, reply *BOut) error {
	reply.B = fmt.Sprintf("ServiceA.B %d", args.B)
	return nil
}

const ServiceA_C = "ServiceA.C"

type CIn struct {
	C int
}
type COut struct {
	C string
}

func (a *ServiceA) C(args *CIn, reply *COut) error {
	reply.C = fmt.Sprintf("ServiceA.C %d", args.C*2)
	return nil
}
