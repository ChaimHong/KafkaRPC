package service1

import "fmt"

type ServiceA struct{}

type AIn struct {
	V int
}
type AOut struct {
	V string
}

const ServiceA_A = "ServiceA.A"
const ServiceA_B = "ServiceA.B"

func (a *ServiceA) A(args *AIn, reply *AOut) error {
	args.V++
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
