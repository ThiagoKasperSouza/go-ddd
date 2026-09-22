package trips

import "context"

// TripRepository define o contrato para persistência de Trip.
type TripRepository interface {
	FindByID(ctx context.Context, id string) (*Trip, error)
	FindAll(ctx context.Context) ([]*Trip, error)
}