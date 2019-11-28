// Package client provides a customersvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/examples/customersvc"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

// New returns a service that's load-balanced over instances of customersvc found
// in the provided Consul server. The mechanism of looking up customersvc
// instances in Consul is hard-coded into the client.
func New(consulAddr string, logger log.Logger) (customersvc.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of customersvc, we declare and enforce these
	// parameters for all of the customersvc consumers.
	var (
		consulService = "customersvc"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints customersvc.Endpoints
	)
	{
		factory := factoryFor(customersvc.MakePostCustomerEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostCustomerEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakeGetCustomerEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetCustomerEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakePutCustomerEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PutCustomerEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakePatchCustomerEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PatchCustomerEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakeDeleteCustomerEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.DeleteCustomerEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakeGetAddressesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetAddressesEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakeGetAddressEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetAddressEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakePostAddressEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostAddressEndpoint = retry
	}
	{
		factory := factoryFor(customersvc.MakeDeleteAddressEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.DeleteAddressEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(customersvc.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := customersvc.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
