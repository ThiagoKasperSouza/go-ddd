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

func (uc *GetShapeUseCase) Execute(ctx context.Context, id string) (*dto.ListShapesOutputDTO, error) {
	shapes,err := uc.shapeRepo.FindByID(ctx,id)
	if err != nil {
		return nil, err
	}

	//transformacao em lista
	var outputShapes []dto.ShapeDTO
	for _,shape := range shapes {


		outputShapes = append(outputShapes,dto.ShapeDTO{
		ShapeID: shape.ID(),      			
		ShapePtLat: shape.ShapePtLat(),
		ShapePtLon: shape.ShapePtLon(),
		ShapePtSequence: shape.ShapePtSequence(),
		})
	}
	return &dto.ListShapesOutputDTO{
		Shapes: outputShapes,
	},nil
}