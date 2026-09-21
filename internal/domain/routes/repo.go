package routes

import "context"

// RouteRepository define o contrato para persistência de Route.
type RouteRepository interface {
	FindByID(ctx context.Context, id string) (*Route, error)
	FindAll(ctx context.Context) ([]*Route, error)
}