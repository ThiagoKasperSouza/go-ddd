package shapes

import (
	"context"

	domainShape "go-ddd/internal/domain/shapes"
	"go-ddd/internal/dto"
)

type GetShapeUseCase struct {
	shapeRepo domainShape.ShapeRepository
}

func NewGetShapeUseCase(shapeRepo domainShape.ShapeRepository) *GetShapeUseCase {
	return &GetShapeUseCase{
		shapeRepo: shapeRepo,
	}
}

func (uc *GetShapeUseCase) Execute(ctx context.Context, id string) (*dto.ShapeDTO, error) {
	shape,err := uc.shapeRepo.FindByID(ctx,id);
	if err != nil {
		return nil,err // domain.ErrShapeNotFound
	}

	return &dto.ShapeDTO{
		ShapeID: shape.ID(),      			
		ShapePtLat: shape.ShapePtLat(),
		ShapePtLon: shape.ShapePtLon(),
		ShapePtSequence: shape.ShapePtSequence(),
		}, nil
}