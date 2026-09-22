package trips

import (
	"context"

	domainTrip "go-ddd/internal/domain/trips"
	"go-ddd/internal/dto"
)

type GetTripUseCase struct {
	tripRepo domainTrip.TripRepository
}

func NewGetTripUseCase(tripRepo domainTrip.TripRepository) *GetTripUseCase {
	return &GetTripUseCase{
		tripRepo: tripRepo,
	}
}

func (uc *GetTripUseCase) Execute(ctx context.Context, id string) (*dto.TripDTO, error) {
	trip,err := uc.tripRepo.FindByID(ctx,id);
	if err != nil {
		return nil,err // domain.ErrTripNotFound
	}

	return &dto.TripDTO{
		RouteID: trip.RouteID(),
		ServiceID: trip.ServiceID(),
		TripID: trip.TripID(),
		TripHeadsign: trip.TripHeadsign(),
		TripShortName: trip.TripShortName(),
		DirectionID: trip.DirectionID(),
		BlockID: trip.BlockID(),
		ShapeID: trip.ShapeID(),
		WheelchairAccessible: trip.WheelchairAccessible(),
		}, nil
}