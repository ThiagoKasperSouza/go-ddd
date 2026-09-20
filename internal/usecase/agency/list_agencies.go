package agency

import (
	"context"
	
	domainAgency "go-ddd/internal/domain/agency"
	"go-ddd/internal/dto"
)

type ListAgenciesUseCase struct {
	agencyRepo domainAgency.AgencyRepository
}

func NewListAgenciesUseCase(agencyRepo domainAgency.AgencyRepository) *ListAgenciesUseCase {
	return &ListAgenciesUseCase{
		agencyRepo: agencyRepo,
	}
}

func (uc *ListAgenciesUseCase) Execute(ctx context.Context) (*dto.ListAgenciesOutputDTO, error) {
	agencies,err := uc.agencyRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	//transformacao em lista
	var outputAgencies []dto.AgencyDTO
	for _,agency := range agencies {

		// Mapeia usando os campos corretos do DTO (AgencyID, AgencyName, etc.)
		// e os getters corretos da entidade (ID(), Name(), Url(), etc.)
		lang := agency.Lang()
		phone := agency.Phone()
		fareURL := agency.FareURL()

		outputAgencies = append(outputAgencies, dto.AgencyDTO{
			AgencyID:       agency.ID(),
			AgencyName:     agency.Name(),
			AgencyURL:      agency.Url(),
			AgencyTimezone: agency.Timezone(),
			AgencyLang:     &lang,
			AgencyPhone:    &phone,
			AgencyFareURL:  &fareURL,
		})
	}
	return &dto.ListAgenciesOutputDTO{
		Agencies: outputAgencies,
	},nil
}