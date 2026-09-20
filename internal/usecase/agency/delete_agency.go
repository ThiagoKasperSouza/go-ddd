package agency

import (
	"context"

	domainAgency "go-ddd/internal/domain/agency"
)

type DeleteAgencyUseCase struct {
	agencyRepo domainAgency.AgencyRepository
}

func NewDeleteAgencyUseCase(agencyRepo domainAgency.AgencyRepository) *DeleteAgencyUseCase {
	return &DeleteAgencyUseCase{
		agencyRepo: agencyRepo,
	}
}

func (uc *DeleteAgencyUseCase) Execute(ctx context.Context, id string) error {
	// Executa a exclusão no repositório.
	// Se a agência não existir, o repositório deve retornar domain.ErrAgencyNotFound
	err := uc.agencyRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}