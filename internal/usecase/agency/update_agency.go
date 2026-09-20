package agency

import (
	"context"

	domainAgency "go-ddd/internal/domain/agency"
	"go-ddd/internal/dto"
)

type UpdateAgencyUseCase struct {
	agencyRepo domainAgency.AgencyRepository
}

func NewUpdateAgencyUseCase(agencyRepo domainAgency.AgencyRepository) *UpdateAgencyUseCase {
	return &UpdateAgencyUseCase{
		agencyRepo: agencyRepo,
	}
}

func (uc *UpdateAgencyUseCase) Execute(ctx context.Context, id string, input dto.UpdateAgencyDTO) (*dto.AgencyDTO, error) {
	// 1. Opcional: Verifica se a agência existe antes de atualizar
	_, err := uc.agencyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err // Retorna domain.ErrAgencyNotFound se nao existir
	}

	// 2. Reconstitui a entidade com o ID existente e novos valores
	updatedAgency, err := domainAgency.NewAgency(
		id,
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

	// 3. Persiste a alteração no banco
	if err := uc.agencyRepo.Update(ctx, updatedAgency); err != nil {
		return nil, err
	}

	// 4. Mapeia para o DTO de retorno
	lang := updatedAgency.Lang()
	phone := updatedAgency.Phone()
	fareURL := updatedAgency.FareURL()

	return &dto.AgencyDTO{
		AgencyID:       updatedAgency.ID(),
		AgencyName:     updatedAgency.Name(),
		AgencyURL:      updatedAgency.Url(),
		AgencyTimezone: updatedAgency.Timezone(),
		AgencyLang:     &lang,
		AgencyPhone:    &phone,
		AgencyFareURL:  &fareURL,
	}, nil
}

// Função auxiliar (pode reaproveitar do create_agency.go se estiver no mesmo pacote)
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}