package agency

import "context"

// AgencyRepository define o contrato para persistência de Agency.
type AgencyRepository interface {
	Save(ctx context.Context, agency *Agency) error
	FindByID(ctx context.Context, id string) (*Agency, error)
	FindAll(ctx context.Context) ([]*Agency, error)
	Update(ctx context.Context, agency *Agency) error
	Delete(ctx context.Context, id string) error
}