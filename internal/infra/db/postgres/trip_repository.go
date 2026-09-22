package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	domainTrip "go-ddd/internal/domain/trips" // substitua pelo caminho correto do seu módulo
)

// Esta linha valida em tempo de compilação se *ATripRepository implementa a interface domain.TripRepository
var _ domainTrip.TripRepository = (*TripRepository)(nil)

type TripRepository struct {
	db *sql.DB
}

// NewTripRepository cria uma nova instância da implementação Postgres
func NewTripRepository(db *sql.DB) *TripRepository {
	return &TripRepository{
		db: db,
	}
}

// FindByID busca uma agência pelo seu ID
func (r *TripRepository) FindByID(ctx context.Context, id string) (*domainTrip.Trip, error) {
	query := `
		SELECT 
		route_id,
		service_id,
		trip_id,
		trip_headsign,
		trip_short_name,
		direction_id INT,
		COALESCE(block_id,'') AS block_id,
		shape_id,
		wheelchair_accessible
		FROM trips
		WHERE trip_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var (
		RouteID		   			string
		ServiceID       		string
		TripID     		  		string
		TripHeadsign			string 
		TripShortName      		string 
		DirectionID       		string
		BlockID     		  	string 
		ShapeID					string
		WheelchairAccessible	string
	)

	err := row.Scan(&RouteID,&ServiceID,&TripID,&TripHeadsign,&TripShortName, &DirectionID, &BlockID, &ShapeID, &WheelchairAccessible)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainTrip.ErrTripNotFound // Mapeia erro de infraestrutura para erro de domínio
		}
		return nil, fmt.Errorf("failed to find Trip by id: %w", err)
	}

	return domainTrip.RestoreTrip(RouteID,ServiceID,TripID,TripHeadsign,TripShortName, DirectionID, BlockID, ShapeID, WheelchairAccessible), nil
}

// FindAll retorna todas as agências cadastradas
func (r *TripRepository) FindAll(ctx context.Context) ([]*domainTrip.Trip, error) {
	query := `
		SELECT 
		route_id,
		service_id,
		trip_id,
		trip_headsign,
		trip_short_name,
		direction_id,
		COALESCE(block_id,'') AS block_id,
		shape_id,
		wheelchair_accessible
		FROM trips
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query Trips: %w", err)
	}
	defer rows.Close()

	var Trips []*domainTrip.Trip

	for rows.Next() {
		var (
			RouteID		   			string
			ServiceID       		string
			TripID     		  		string
			TripHeadsign			string 
			TripShortName      		string 
			DirectionID       		string
			BlockID     		  	string 
			ShapeID					string
			WheelchairAccessible	string
		)

		if err := rows.Scan(&RouteID,&ServiceID,&TripID,&TripHeadsign,&TripShortName, &DirectionID, &BlockID, &ShapeID, &WheelchairAccessible); err != nil {
			return nil, fmt.Errorf("failed to scan Trips row: %w", err)
		}

		Trips = append(Trips,domainTrip.RestoreTrip(RouteID,ServiceID,TripID,TripHeadsign,TripShortName, DirectionID, BlockID, ShapeID, WheelchairAccessible))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during Trip iteration: %w", err)
	}

	return Trips, nil
}
