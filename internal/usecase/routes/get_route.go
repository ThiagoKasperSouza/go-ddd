package routes

import (
	"context"

	domainRoute "go-ddd/internal/domain/routes"
	"go-ddd/internal/dto"
)

type GetRouteUseCase struct {
	routeRepo domainRoute.RouteRepository
}

func NewGetRouteUseCase(routeRepo domainRoute.RouteRepository) *GetRouteUseCase {
	return &GetRouteUseCase{
		routeRepo: routeRepo,
	}
}

func (uc *GetRouteUseCase) Execute(ctx context.Context, id string) (*dto.RouteDTO, error) {
	route,err := uc.routeRepo.FindByID(ctx,id);
	if err != nil {
		return nil,err // domain.ErrRouteNotFound
	}

	return &dto.RouteDTO{
		RouteID: route.ID(),
		AgencyID:route.AgencyID(),
		RouteShortName: route.ShortName(),
		RouteLongName: route.LongName(),
		RouteDesc: route.Desc(),
		RouteType: route.Type(),
		RouteURL: route.RouteURL(),
		RouteColor: route.RouteColor(),
		RouteTextColor: route.RouteTextColor(),
		}, nil
}