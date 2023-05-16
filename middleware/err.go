package middleware

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
)

func ErrResponseMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response interface{}) error {
		err := next(ctx, request, response)
		return err
	}
}
