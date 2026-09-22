package trips

import (
	"context"
	
	domainTrip "go-ddd/internal/domain/trips"
	"go-ddd/internal/dto"
)

type ListTripsUseCase struct {
	tripRepo domainTrip.TripRepository
}

func NewListTripsUseCase(tripRepo domainTrip.TripRepository) *ListTripsUseCase {
	return &ListTripsUseCase{
		tripRepo: tripRepo,
	}
}

func (uc *ListTripsUseCase) Execute(ctx context.Context) (*dto.ListTripsOutputDTO, error) {
	trips,err := uc.tripRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	//transformacao em lista
	var outputTrips []dto.TripDTO
	for _,trip := range trips {


		outputTrips = append(outputTrips,dto.TripDTO{
		RouteID: trip.RouteID(),
		ServiceID: trip.ServiceID(),
		TripID: trip.TripID(),
		TripHeadsign: trip.TripHeadsign(),
		TripShortName: trip.TripShortName(),
		DirectionID: trip.DirectionID(),
		BlockID: trip.BlockID(),
		ShapeID: trip.ShapeID(),
		WheelchairAccessible: trip.WheelchairAccessible(),
		})
	}
	return &dto.ListTripsOutputDTO{
		Trips: outputTrips,
	},nil
}