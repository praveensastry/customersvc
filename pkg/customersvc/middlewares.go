package customersvc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostCustomer(ctx context.Context, p Customer) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostCustomer", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostCustomer(ctx, p)
}

func (mw loggingMiddleware) GetCustomer(ctx context.Context, id string) (p Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetCustomer", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetCustomer(ctx, id)
}

func (mw loggingMiddleware) PutCustomer(ctx context.Context, id string, p Customer) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutCustomer", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutCustomer(ctx, id, p)
}

func (mw loggingMiddleware) PatchCustomer(ctx context.Context, id string, p Customer) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchCustomer", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchCustomer(ctx, id, p)
}

func (mw loggingMiddleware) DeleteCustomer(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteCustomer", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteCustomer(ctx, id)
}

func (mw loggingMiddleware) GetAddresses(ctx context.Context, customerID string) (addresses []Address, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAddresses", "customerID", customerID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAddresses(ctx, customerID)
}

func (mw loggingMiddleware) GetAddress(ctx context.Context, customerID string, addressID string) (a Address, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAddress", "customerID", customerID, "addressID", addressID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAddress(ctx, customerID, addressID)
}

func (mw loggingMiddleware) PostAddress(ctx context.Context, customerID string, a Address) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostAddress", "customerID", customerID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostAddress(ctx, customerID, a)
}

func (mw loggingMiddleware) DeleteAddress(ctx context.Context, customerID string, addressID string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteAddress", "customerID", customerID, "addressID", addressID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteAddress(ctx, customerID, addressID)
}
