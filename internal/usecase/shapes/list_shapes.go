package shapes

import (
	"context"
	
	domainShape "go-ddd/internal/domain/shapes"
	"go-ddd/internal/dto"
)

type ListShapesUseCase struct {
	shapeRepo domainShape.ShapeRepository
}

func NewListShapesUseCase(shapeRepo domainShape.ShapeRepository) *ListShapesUseCase {
	return &ListShapesUseCase{
		shapeRepo: shapeRepo,
	}
}

func (uc *ListShapesUseCase) Execute(ctx context.Context) (*dto.ListShapesOutputDTO, error) {
	shapes,err := uc.shapeRepo.FindAll(ctx)
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