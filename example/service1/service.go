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

func (a *ServiceA) B(args *int, reply *string) error {
	*reply = fmt.Sprintf("ServiceA.B %d", *args)
	return nil
}
