package customersvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints collects all of the endpoints that compose a customer service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	PostCustomerEndpoint   endpoint.Endpoint
	GetCustomerEndpoint    endpoint.Endpoint
	PutCustomerEndpoint    endpoint.Endpoint
	PatchCustomerEndpoint  endpoint.Endpoint
	DeleteCustomerEndpoint endpoint.Endpoint
	GetAddressesEndpoint   endpoint.Endpoint
	GetAddressEndpoint     endpoint.Endpoint
	PostAddressEndpoint    endpoint.Endpoint
	DeleteAddressEndpoint  endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a customersvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostCustomerEndpoint:   MakePostCustomerEndpoint(s),
		GetCustomerEndpoint:    MakeGetCustomerEndpoint(s),
		PutCustomerEndpoint:    MakePutCustomerEndpoint(s),
		PatchCustomerEndpoint:  MakePatchCustomerEndpoint(s),
		DeleteCustomerEndpoint: MakeDeleteCustomerEndpoint(s),
		GetAddressesEndpoint:   MakeGetAddressesEndpoint(s),
		GetAddressEndpoint:     MakeGetAddressEndpoint(s),
		PostAddressEndpoint:    MakePostAddressEndpoint(s),
		DeleteAddressEndpoint:  MakeDeleteAddressEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a customersvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{
		PostCustomerEndpoint:   httptransport.NewClient("POST", tgt, encodePostCustomerRequest, decodePostCustomerResponse, options...).Endpoint(),
		GetCustomerEndpoint:    httptransport.NewClient("GET", tgt, encodeGetCustomerRequest, decodeGetCustomerResponse, options...).Endpoint(),
		PutCustomerEndpoint:    httptransport.NewClient("PUT", tgt, encodePutCustomerRequest, decodePutCustomerResponse, options...).Endpoint(),
		PatchCustomerEndpoint:  httptransport.NewClient("PATCH", tgt, encodePatchCustomerRequest, decodePatchCustomerResponse, options...).Endpoint(),
		DeleteCustomerEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteCustomerRequest, decodeDeleteCustomerResponse, options...).Endpoint(),
		GetAddressesEndpoint:   httptransport.NewClient("GET", tgt, encodeGetAddressesRequest, decodeGetAddressesResponse, options...).Endpoint(),
		GetAddressEndpoint:     httptransport.NewClient("GET", tgt, encodeGetAddressRequest, decodeGetAddressResponse, options...).Endpoint(),
		PostAddressEndpoint:    httptransport.NewClient("POST", tgt, encodePostAddressRequest, decodePostAddressResponse, options...).Endpoint(),
		DeleteAddressEndpoint:  httptransport.NewClient("DELETE", tgt, encodeDeleteAddressRequest, decodeDeleteAddressResponse, options...).Endpoint(),
	}, nil
}

// PostCustomer implements Service. Primarily useful in a client.
func (e Endpoints) PostCustomer(ctx context.Context, p Customer) error {
	request := postCustomerRequest{Customer: p}
	response, err := e.PostCustomerEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postCustomerResponse)
	return resp.Err
}

// GetCustomer implements Service. Primarily useful in a client.
func (e Endpoints) GetCustomer(ctx context.Context, id string) (Customer, error) {
	request := getCustomerRequest{ID: id}
	response, err := e.GetCustomerEndpoint(ctx, request)
	if err != nil {
		return Customer{}, err
	}
	resp := response.(getCustomerResponse)
	return resp.Customer, resp.Err
}

// PutCustomer implements Service. Primarily useful in a client.
func (e Endpoints) PutCustomer(ctx context.Context, id string, p Customer) error {
	request := putCustomerRequest{ID: id, Customer: p}
	response, err := e.PutCustomerEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putCustomerResponse)
	return resp.Err
}

// PatchCustomer implements Service. Primarily useful in a client.
func (e Endpoints) PatchCustomer(ctx context.Context, id string, p Customer) error {
	request := patchCustomerRequest{ID: id, Customer: p}
	response, err := e.PatchCustomerEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(patchCustomerResponse)
	return resp.Err
}

// DeleteCustomer implements Service. Primarily useful in a client.
func (e Endpoints) DeleteCustomer(ctx context.Context, id string) error {
	request := deleteCustomerRequest{ID: id}
	response, err := e.DeleteCustomerEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteCustomerResponse)
	return resp.Err
}

// GetAddresses implements Service. Primarily useful in a client.
func (e Endpoints) GetAddresses(ctx context.Context, customerID string) ([]Address, error) {
	request := getAddressesRequest{CustomerID: customerID}
	response, err := e.GetAddressesEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(getAddressesResponse)
	return resp.Addresses, resp.Err
}

// GetAddress implements Service. Primarily useful in a client.
func (e Endpoints) GetAddress(ctx context.Context, customerID string, addressID string) (Address, error) {
	request := getAddressRequest{CustomerID: customerID, AddressID: addressID}
	response, err := e.GetAddressEndpoint(ctx, request)
	if err != nil {
		return Address{}, err
	}
	resp := response.(getAddressResponse)
	return resp.Address, resp.Err
}

// PostAddress implements Service. Primarily useful in a client.
func (e Endpoints) PostAddress(ctx context.Context, customerID string, a Address) error {
	request := postAddressRequest{CustomerID: customerID, Address: a}
	response, err := e.PostAddressEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postAddressResponse)
	return resp.Err
}

// DeleteAddress implements Service. Primarily useful in a client.
func (e Endpoints) DeleteAddress(ctx context.Context, customerID string, addressID string) error {
	request := deleteAddressRequest{CustomerID: customerID, AddressID: addressID}
	response, err := e.DeleteAddressEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteAddressResponse)
	return resp.Err
}

// MakePostCustomerEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postCustomerRequest)
		e := s.PostCustomer(ctx, req.Customer)
		return postCustomerResponse{Err: e}, nil
	}
}

// MakeGetCustomerEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCustomerRequest)
		p, e := s.GetCustomer(ctx, req.ID)
		return getCustomerResponse{Customer: p, Err: e}, nil
	}
}

// MakePutCustomerEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putCustomerRequest)
		e := s.PutCustomer(ctx, req.ID, req.Customer)
		return putCustomerResponse{Err: e}, nil
	}
}

// MakePatchCustomerEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePatchCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchCustomerRequest)
		e := s.PatchCustomer(ctx, req.ID, req.Customer)
		return patchCustomerResponse{Err: e}, nil
	}
}

// MakeDeleteCustomerEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteCustomerRequest)
		e := s.DeleteCustomer(ctx, req.ID)
		return deleteCustomerResponse{Err: e}, nil
	}
}

// MakeGetAddressesEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetAddressesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAddressesRequest)
		a, e := s.GetAddresses(ctx, req.CustomerID)
		return getAddressesResponse{Addresses: a, Err: e}, nil
	}
}

// MakeGetAddressEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetAddressEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAddressRequest)
		a, e := s.GetAddress(ctx, req.CustomerID, req.AddressID)
		return getAddressResponse{Address: a, Err: e}, nil
	}
}

// MakePostAddressEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostAddressEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postAddressRequest)
		e := s.PostAddress(ctx, req.CustomerID, req.Address)
		return postAddressResponse{Err: e}, nil
	}
}

// MakeDeleteAddressEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteAddressEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteAddressRequest)
		e := s.DeleteAddress(ctx, req.CustomerID, req.AddressID)
		return deleteAddressResponse{Err: e}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type postCustomerRequest struct {
	Customer Customer
}

type postCustomerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postCustomerResponse) error() error { return r.Err }

type getCustomerRequest struct {
	ID string
}

type getCustomerResponse struct {
	Customer Customer `json:"customer,omitempty"`
	Err      error    `json:"err,omitempty"`
}

func (r getCustomerResponse) error() error { return r.Err }

type putCustomerRequest struct {
	ID       string
	Customer Customer
}

type putCustomerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r putCustomerResponse) error() error { return nil }

type patchCustomerRequest struct {
	ID       string
	Customer Customer
}

type patchCustomerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r patchCustomerResponse) error() error { return r.Err }

type deleteCustomerRequest struct {
	ID string
}

type deleteCustomerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteCustomerResponse) error() error { return r.Err }

type getAddressesRequest struct {
	CustomerID string
}

type getAddressesResponse struct {
	Addresses []Address `json:"addresses,omitempty"`
	Err       error     `json:"err,omitempty"`
}

func (r getAddressesResponse) error() error { return r.Err }

type getAddressRequest struct {
	CustomerID string
	AddressID  string
}

type getAddressResponse struct {
	Address Address `json:"address,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getAddressResponse) error() error { return r.Err }

type postAddressRequest struct {
	CustomerID string
	Address    Address
}

type postAddressResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postAddressResponse) error() error { return r.Err }

type deleteAddressRequest struct {
	CustomerID string
	AddressID  string
}

type deleteAddressResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteAddressResponse) error() error { return r.Err }
