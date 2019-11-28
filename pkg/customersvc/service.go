package customersvc

import (
	"context"
	"errors"
	"sync"
)

// Service is a simple CRUD interface for user customers.
type Service interface {
	PostCustomer(ctx context.Context, p Customer) error
	GetCustomer(ctx context.Context, id string) (Customer, error)
	PutCustomer(ctx context.Context, id string, p Customer) error
	PatchCustomer(ctx context.Context, id string, p Customer) error
	DeleteCustomer(ctx context.Context, id string) error
	GetAddresses(ctx context.Context, customerID string) ([]Address, error)
	GetAddress(ctx context.Context, customerID string, addressID string) (Address, error)
	PostAddress(ctx context.Context, customerID string, a Address) error
	DeleteAddress(ctx context.Context, customerID string, addressID string) error
}

// Customer represents a single user customer.
// ID should be globally unique.
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
}

// Address is a field of a user customer.
// ID should be unique within the customer (at a minimum).
type Address struct {
	ID       string `json:"id"`
	Location string `json:"location,omitempty"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Customer
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]Customer{},
	}
}

func (s *inmemService) PostCustomer(ctx context.Context, p Customer) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[p.ID]; ok {
		return ErrAlreadyExists // POST = create, don't overwrite
	}
	s.m[p.ID] = p
	return nil
}

func (s *inmemService) GetCustomer(ctx context.Context, id string) (Customer, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return Customer{}, ErrNotFound
	}
	return p, nil
}

func (s *inmemService) PutCustomer(ctx context.Context, id string, p Customer) error {
	if id != p.ID {
		return ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[id] = p // PUT = create or update
	return nil
}

func (s *inmemService) PatchCustomer(ctx context.Context, id string, p Customer) error {
	if p.ID != "" && id != p.ID {
		return ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	existing, ok := s.m[id]
	if !ok {
		return ErrNotFound // PATCH = update existing, don't create
	}

	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the Customer definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if p.Name != "" {
		existing.Name = p.Name
	}
	if len(p.Addresses) > 0 {
		existing.Addresses = p.Addresses
	}
	s.m[id] = existing
	return nil
}

func (s *inmemService) DeleteCustomer(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}

func (s *inmemService) GetAddresses(ctx context.Context, customerID string) ([]Address, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[customerID]
	if !ok {
		return []Address{}, ErrNotFound
	}
	return p.Addresses, nil
}

func (s *inmemService) GetAddress(ctx context.Context, customerID string, addressID string) (Address, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[customerID]
	if !ok {
		return Address{}, ErrNotFound
	}
	for _, address := range p.Addresses {
		if address.ID == addressID {
			return address, nil
		}
	}
	return Address{}, ErrNotFound
}

func (s *inmemService) PostAddress(ctx context.Context, customerID string, a Address) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[customerID]
	if !ok {
		return ErrNotFound
	}
	for _, address := range p.Addresses {
		if address.ID == a.ID {
			return ErrAlreadyExists
		}
	}
	p.Addresses = append(p.Addresses, a)
	s.m[customerID] = p
	return nil
}

func (s *inmemService) DeleteAddress(ctx context.Context, customerID string, addressID string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[customerID]
	if !ok {
		return ErrNotFound
	}
	newAddresses := make([]Address, 0, len(p.Addresses))
	for _, address := range p.Addresses {
		if address.ID == addressID {
			continue // delete
		}
		newAddresses = append(newAddresses, address)
	}
	if len(newAddresses) == len(p.Addresses) {
		return ErrNotFound
	}
	p.Addresses = newAddresses
	s.m[customerID] = p
	return nil
}
