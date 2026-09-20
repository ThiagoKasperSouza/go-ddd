package agency

import (
	"context"

	domainAgency "go-ddd/internal/domain/agency"
	"go-ddd/internal/dto"
)

type GetAgencyUseCase struct {
	agencyRepo domainAgency.AgencyRepository
}

func NewGetAgencyUseCase(agencyRepo domainAgency.AgencyRepository) *GetAgencyUseCase {
	return &GetAgencyUseCase{
		agencyRepo: agencyRepo,
	}
}

func (uc *GetAgencyUseCase) Execute(ctx context.Context, id string) (*dto.AgencyDTO, error) {
	agency,err := uc.agencyRepo.FindByID(ctx,id);
	if err != nil {
		return nil,err // domain.ErrAgencyNotFound
	}

	// Mapeia usando os campos corretos do DTO (AgencyID, AgencyName, etc.)
	// e os getters corretos da entidade (ID(), Name(), Url(), etc.)
	lang := agency.Lang()
	phone := agency.Phone()
	fareURL := agency.FareURL()

	return &dto.AgencyDTO{
		AgencyID:       agency.ID(),
		AgencyName:     agency.Name(),
		AgencyURL:      agency.Url(),
		AgencyTimezone: agency.Timezone(),
		AgencyLang:     &lang,
		AgencyPhone:    &phone,
		AgencyFareURL:  &fareURL,
		}, nil
}