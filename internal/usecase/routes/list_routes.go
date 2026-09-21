package routes

import (
	"context"
	
	domainRoute "go-ddd/internal/domain/routes"
	"go-ddd/internal/dto"
)

type ListRoutesUseCase struct {
	routeRepo domainRoute.RouteRepository
}

func NewListRoutesUseCase(routeRepo domainRoute.RouteRepository) *ListRoutesUseCase {
	return &ListRoutesUseCase{
		routeRepo: routeRepo,
	}
}

func (uc *ListRoutesUseCase) Execute(ctx context.Context) (*dto.ListRoutesOutputDTO, error) {
	routes,err := uc.routeRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	//transformacao em lista
	var outputRoutes []dto.RouteDTO
	for _,route := range routes {


		outputRoutes = append(outputRoutes, dto.RouteDTO{
			RouteID: route.ID(),
			AgencyID:route.AgencyID(),
			RouteShortName: route.ShortName(),
			RouteLongName: route.LongName(),
			RouteDesc: route.Desc(),
			RouteType: route.Type(),
			RouteURL: route.RouteURL(),
			RouteColor: route.RouteColor(),
			RouteTextColor: route.RouteTextColor(),
		})
	}
	return &dto.ListRoutesOutputDTO{
		Routes: outputRoutes,
	},nil
}