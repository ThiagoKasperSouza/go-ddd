package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	domainRoute "go-ddd/internal/domain/routes" // substitua pelo caminho correto do seu módulo
)

// Esta linha valida em tempo de compilação se *ARouteRepository implementa a interface domain.RouteRepository
var _ domainRoute.RouteRepository = (*RouteRepository)(nil)

type RouteRepository struct {
	db *sql.DB
}

// NewRouteRepository cria uma nova instância da implementação Postgres
func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{
		db: db,
	}
}

// FindByID busca uma agência pelo seu ID
func (r *RouteRepository) FindByID(ctx context.Context, id string) (*domainRoute.Route, error) {
	query := `
		SELECT 
		route_id,
		agency_id,
		route_short_name,
		route_long_name,
		COALESCE(route_desc, '') AS route_desc,
		route_type,
		route_url,
		route_color,
		route_text_color
		FROM routes
		WHERE route_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var (
		routeId      			string 
		agencyId	    		string 
		routeShortName      	string 
		routeLongName			string 
		routeDesc     		   	string 
		routeType    			string 
		routeURL  				string 
		routeColor  			string
		routeTextColor  		string
	)

	err := row.Scan(&routeId,&agencyId,&routeShortName,&routeLongName, &routeDesc, &routeType, &routeURL, &routeColor, &routeTextColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainRoute.ErrRouteNotFound // Mapeia erro de infraestrutura para erro de domínio
		}
		return nil, fmt.Errorf("failed to find route by id: %w", err)
	}

	return domainRoute.RestoreRoute(routeId,agencyId,routeShortName,routeLongName, routeDesc, routeType, routeURL, routeColor, routeTextColor), nil
}

// FindAll retorna todas as agências cadastradas
func (r *RouteRepository) FindAll(ctx context.Context) ([]*domainRoute.Route, error) {
	query := `
		SELECT 
		route_id,
		agency_id,
		route_short_name,
		route_long_name,
		COALESCE(route_desc, '') AS route_desc,
		route_type,
		route_url,
		route_color,
		route_text_color
		FROM routes
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query routes: %w", err)
	}
	defer rows.Close()

	var routes []*domainRoute.Route

	for rows.Next() {
		var (
			routeId      			string 
			agencyId	    		string 
			routeShortName      	string 
			routeLongName			string 
			routeDesc     		   	string 
			routeType    			string 
			routeURL  				string 
			routeColor  			string
			routeTextColor  		string
		)

		if err := rows.Scan(&routeId,&agencyId,&routeShortName,&routeLongName, &routeDesc, &routeType, &routeURL, &routeColor, &routeTextColor); err != nil {
			return nil, fmt.Errorf("failed to scan routes row: %w", err)
		}

		routes = append(routes,domainRoute.RestoreRoute(routeId,agencyId,routeShortName,routeLongName, routeDesc, routeType, routeURL, routeColor, routeTextColor))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during route iteration: %w", err)
	}

	return routes, nil
}
