package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	domainShape "go-ddd/internal/domain/shapes" // substitua pelo caminho correto do seu módulo
)

// Esta linha valida em tempo de compilação se *AShapeRepository implementa a interface domain.ShapeRepository
var _ domainShape.ShapeRepository = (*ShapeRepository)(nil)

type ShapeRepository struct {
	db *sql.DB
}

// NewShapeRepository cria uma nova instância da implementação Postgres
func NewShapeRepository(db *sql.DB) *ShapeRepository {
	return &ShapeRepository{
		db: db,
	}
}

// FindByID busca uma agência pelo seu ID
func (r *ShapeRepository) FindByID(ctx context.Context, id string) (*domainShape.Shape, error) {
	query := `
		SELECT 
		shape_id,
		shape_pt_lat,
		shape_pt_lon,
		shape_pt_sequence
		FROM shapes
		WHERE shape_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var (
		ShapeID       		string 
		ShapePtLat       	string 
		ShapePtLon     		string 
		ShapePtSequence		string
	)

	err := row.Scan(&ShapeID,&ShapePtLat, &ShapePtLon, &ShapePtSequence)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainShape.ErrShapeNotFound // Mapeia erro de infraestrutura para erro de domínio
		}
		return nil, fmt.Errorf("failed to find shape by id: %w", err)
	}

	return domainShape.RestoreShape(ShapeID,ShapePtLat, ShapePtLon, ShapePtSequence), nil
}

// FindAll retorna todas as agências cadastradas
func (r *ShapeRepository) FindAll(ctx context.Context) ([]*domainShape.Shape, error) {
	query := `
		SELECT 
		shape_id,
		shape_pt_lat,
		shape_pt_lon,
		shape_pt_sequence
		FROM shapes
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query Shapes: %w", err)
	}
	defer rows.Close()

	var Shapes []*domainShape.Shape

	for rows.Next() {
		var (
			ShapeID       		string 
			ShapePtLat       	string 
			ShapePtLon     		string 
			ShapePtSequence		string
		)

		if err := rows.Scan(&ShapeID,&ShapePtLat, &ShapePtLon, &ShapePtSequence); err != nil {
			return nil, fmt.Errorf("failed to scan Shapes row: %w", err)
		}

		Shapes = append(Shapes,domainShape.RestoreShape(ShapeID,ShapePtLat, ShapePtLon, ShapePtSequence))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during Shape iteration: %w", err)
	}

	return Shapes, nil
}
