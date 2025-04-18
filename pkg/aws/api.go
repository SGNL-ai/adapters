// Copyright 2025 SGNL.ai, Inc.
package aws

import (
	"context"
)

// EntityGetter is an interface for entities that require a Get operation.
type EntityGetter[T any] interface {
	// The Get method retrieves the specified entity.
	//
	// Parameters:
	//
	//   ctx: Context for cancellation and deadlines.
	//   entity: Entity of type T to be retrieved.
	//
	// Returns:
	//
	//   T: Retrieved entity detail of type T.
	//   error: Error, if any.
	Get(ctx context.Context, entity T) (T, error)
}

// EntityLister is an interface for entities that require a List operation.
type EntityLister[T any] interface {
	// The List method retrieves a list of entities with pagination support.
	//
	// Parameters:
	//
	//   ctx: Context for cancellation and deadlines.
	//   opts: Options to control the List operation, such as filtering and pagination.
	//
	// Returns:
	//
	//   []T: Slice of entities of type T.
	//   *string: Token for the next page.
	//   error: Error, if any.
	List(ctx context.Context, opts *Options) ([]T, *string, error)
}
