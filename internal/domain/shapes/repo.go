package shapes

import "context"

// ShapeRepository define o contrato para persistência de Shape.
type ShapeRepository interface {
	FindByID(ctx context.Context, id string) (*Shape, error)
	FindAll(ctx context.Context) ([]*Shape, error)
}