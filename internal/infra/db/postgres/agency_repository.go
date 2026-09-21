package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	domainAgency "go-ddd/internal/domain/agency" // substitua pelo caminho correto do seu módulo
)

// Esta linha valida em tempo de compilação se *AgencyRepository implementa a interface domain.AgencyRepository
var _ domainAgency.AgencyRepository = (*AgencyRepository)(nil)

type AgencyRepository struct {
	db *sql.DB
}

// NewAgencyRepository cria uma nova instância da implementação Postgres
func NewAgencyRepository(db *sql.DB) *AgencyRepository {
	return &AgencyRepository{
		db: db,
	}
}

// Save insere uma nova agência no banco de dados
func (r *AgencyRepository) Save(ctx context.Context, agency *domainAgency.Agency) error {
	query := `
		INSERT INTO agency (agency_name, agency_url, agency_timezone, agency_lang, agency_phone, agency_fare_url)
		VALUES ($1,$2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		agency.Name(),
		agency.Url(),
		agency.Timezone(),
		agency.Lang(),
		agency.Phone(),
		agency.FareURL(),
	)

	if err != nil {
		return fmt.Errorf("failed to save agency: %w", err)
	}

	return nil
}

// FindByID busca uma agência pelo seu ID
func (r *AgencyRepository) FindByID(ctx context.Context, id string) (*domainAgency.Agency, error) {
	query := `
		SELECT agency_id, agency_name, agency_url, agency_timezone, agency_lang, agency_phone, agency_fare_url
		FROM agency
		WHERE agency_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var (
		agencyID string
		name     string
		url      string
		timezone string
		lang     string
		phone    string
		fareURL  string
	)

	err := row.Scan(&agencyID, &name, &url, &timezone, &lang, &phone, &fareURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainAgency.ErrAgencyNotFound // Mapeia erro de infraestrutura para erro de domínio
		}
		return nil, fmt.Errorf("failed to find agency by id: %w", err)
	}

	return domainAgency.RestoreAgency(agencyID, name, url, timezone, lang, phone, fareURL), nil
}

// FindAll retorna todas as agências cadastradas
func (r *AgencyRepository) FindAll(ctx context.Context) ([]*domainAgency.Agency, error) {
	query := `
		SELECT agency_id, agency_name, agency_url, agency_timezone, agency_lang, agency_phone, agency_fare_url
		FROM agency
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query agencies: %w", err)
	}
	defer rows.Close()

	var agencies []*domainAgency.Agency

	for rows.Next() {
		var (
			id       string
			name     string
			url      string
			timezone string
			lang     string
			phone    string
			fareURL  string
		)

		if err := rows.Scan(&id, &name, &url, &timezone, &lang, &phone, &fareURL); err != nil {
			return nil, fmt.Errorf("failed to scan agency row: %w", err)
		}

		agencies = append(agencies, domainAgency.RestoreAgency(id, name, url, timezone, lang, phone, fareURL))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during agency iteration: %w", err)
	}

	return agencies, nil
}

// Update atualiza os dados de uma agência existente
func (r *AgencyRepository) Update(ctx context.Context, agency *domainAgency.Agency) error {
	query := `
		UPDATE agency
		SET agency_name = $1, agency_url = $2, agency_timezone = $3, agency_lang = $4, agency_phone = $5, agency_fare_url = $6
		WHERE agency_id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		agency.Name(),
		agency.Url(),
		agency.Timezone(),
		agency.Lang(),
		agency.Phone(),
		agency.FareURL(),
		agency.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domainAgency.ErrAgencyNotFound
	}

	return nil
}

// Delete remove uma agência pelo ID
func (r *AgencyRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM agency WHERE agency_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete agency: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domainAgency.ErrAgencyNotFound
	}

	return nil
}