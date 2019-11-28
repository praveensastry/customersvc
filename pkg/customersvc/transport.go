package customersvc

// The customersvc is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a customersvc server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /customers/                          adds another customer
	// GET     /customers/:id                       retrieves the given customer by id
	// PUT     /customers/:id                       post updated customer information about the customer
	// PATCH   /customers/:id                       partial updated customer information
	// DELETE  /customers/:id                       remove the given customer
	// GET     /customers/:id/addresses/            retrieve addresses associated with the customer
	// GET     /customers/:id/addresses/:addressID  retrieve a particular customer address
	// POST    /customers/:id/addresses/            add a new address
	// DELETE  /customers/:id/addresses/:addressID  remove an address

	r.Methods("POST").Path("/customers/").Handler(httptransport.NewServer(
		e.PostCustomerEndpoint,
		decodePostCustomerRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/customers/{id}").Handler(httptransport.NewServer(
		e.GetCustomerEndpoint,
		decodeGetCustomerRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/customers/{id}").Handler(httptransport.NewServer(
		e.PutCustomerEndpoint,
		decodePutCustomerRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/customers/{id}").Handler(httptransport.NewServer(
		e.PatchCustomerEndpoint,
		decodePatchCustomerRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/customers/{id}").Handler(httptransport.NewServer(
		e.DeleteCustomerEndpoint,
		decodeDeleteCustomerRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/customers/{id}/addresses/").Handler(httptransport.NewServer(
		e.GetAddressesEndpoint,
		decodeGetAddressesRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/customers/{id}/addresses/{addressID}").Handler(httptransport.NewServer(
		e.GetAddressEndpoint,
		decodeGetAddressRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/customers/{id}/addresses/").Handler(httptransport.NewServer(
		e.PostAddressEndpoint,
		decodePostAddressRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/customers/{id}/addresses/{addressID}").Handler(httptransport.NewServer(
		e.DeleteAddressEndpoint,
		decodeDeleteAddressRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postCustomerRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Customer); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getCustomerRequest{ID: id}, nil
}

func decodePutCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		return nil, err
	}
	return putCustomerRequest{
		ID:       id,
		Customer: customer,
	}, nil
}

func decodePatchCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		return nil, err
	}
	return patchCustomerRequest{
		ID:       id,
		Customer: customer,
	}, nil
}

func decodeDeleteCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteCustomerRequest{ID: id}, nil
}

func decodeGetAddressesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getAddressesRequest{CustomerID: id}, nil
}

func decodeGetAddressRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	addressID, ok := vars["addressID"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getAddressRequest{
		CustomerID: id,
		AddressID:  addressID,
	}, nil
}

func decodePostAddressRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var address Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		return nil, err
	}
	return postAddressRequest{
		CustomerID: id,
		Address:    address,
	}, nil
}

func decodeDeleteAddressRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	addressID, ok := vars["addressID"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteAddressRequest{
		CustomerID: id,
		AddressID:  addressID,
	}, nil
}

func encodePostCustomerRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/customers/")
	req.URL.Path = "/customers/"
	return encodeRequest(ctx, req, request)
}

func encodeGetCustomerRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/customers/{id}")
	r := request.(getCustomerRequest)
	customerID := url.QueryEscape(r.ID)
	req.URL.Path = "/customers/" + customerID
	return encodeRequest(ctx, req, request)
}

func encodePutCustomerRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/customers/{id}")
	r := request.(putCustomerRequest)
	customerID := url.QueryEscape(r.ID)
	req.URL.Path = "/customers/" + customerID
	return encodeRequest(ctx, req, request)
}

func encodePatchCustomerRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PATCH").Path("/customers/{id}")
	r := request.(patchCustomerRequest)
	customerID := url.QueryEscape(r.ID)
	req.URL.Path = "/customers/" + customerID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteCustomerRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/customers/{id}")
	r := request.(deleteCustomerRequest)
	customerID := url.QueryEscape(r.ID)
	req.URL.Path = "/customers/" + customerID
	return encodeRequest(ctx, req, request)
}

func encodeGetAddressesRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/customers/{id}/addresses/")
	r := request.(getAddressesRequest)
	customerID := url.QueryEscape(r.CustomerID)
	req.URL.Path = "/customers/" + customerID + "/addresses/"
	return encodeRequest(ctx, req, request)
}

func encodeGetAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/customers/{id}/addresses/{addressID}")
	r := request.(getAddressRequest)
	customerID := url.QueryEscape(r.CustomerID)
	addressID := url.QueryEscape(r.AddressID)
	req.URL.Path = "/customers/" + customerID + "/addresses/" + addressID
	return encodeRequest(ctx, req, request)
}

func encodePostAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/customers/{id}/addresses/")
	r := request.(postAddressRequest)
	customerID := url.QueryEscape(r.CustomerID)
	req.URL.Path = "/customers/" + customerID + "/addresses/"
	return encodeRequest(ctx, req, request)
}

func encodeDeleteAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/customers/{id}/addresses/{addressID}")
	r := request.(deleteAddressRequest)
	customerID := url.QueryEscape(r.CustomerID)
	addressID := url.QueryEscape(r.AddressID)
	req.URL.Path = "/customers/" + customerID + "/addresses/" + addressID
	return encodeRequest(ctx, req, request)
}

func decodePostCustomerResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postCustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetCustomerResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getCustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutCustomerResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putCustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePatchCustomerResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response patchCustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteCustomerResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteCustomerResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetAddressesResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getAddressesResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePostAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// customersvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
