package agency

import (
	"context"

	domainAgency "go-ddd/internal/domain/agency"
	"go-ddd/internal/dto"
)

type CreateAgencyUseCase struct {
	agencyRepo domainAgency.AgencyRepository
}

func NewCreateAgencyUseCase(agencyRepo domainAgency.AgencyRepository) *CreateAgencyUseCase {
	return &CreateAgencyUseCase{
		agencyRepo: agencyRepo,
	}
}

func (uc *CreateAgencyUseCase) Execute(ctx context.Context, input dto.CreateAgencyDTO) (*dto.AgencyDTO, error) {
	// Cria a entidade de domínio tratando os campos ponteiros (*string)
	agency, err := domainAgency.NewAgency(
		input.AgencyID,
		input.AgencyName,
		input.AgencyURL,
		input.AgencyTimezone,
		derefString(input.AgencyLang),
		derefString(input.AgencyPhone),
		derefString(input.AgencyFareURL),
	)
	if err != nil {
		return nil, err
	}

	// Persiste a entidade usando o repositório
	if err := uc.agencyRepo.Save(ctx, agency); err != nil {
		return nil, err
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

// Função auxiliar para desreferenciar *string de forma segura
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}