package routes

import "errors"

var (
	ErrRouteNotFound = errors.New("Route not found")
	ErrInvalidRoute  = errors.New("invalid Route data")
)

type Route struct {
	id      				string 
	agencyId	    		string 
	routeShortName      	string 
	routeLongName			string 
	routeDesc     		   	string 
	routeType    			string 
	routeURL  				string 
	routeColor  			string
	routeTextColor  		string
}

// NewAgency garante que uma nova agência seja criada com dados válidos
func NewRoute(id,
	agencyId,
	routeShortName,
	routeLongName, 
	routeDesc,     		  
	routeType,     
	routeURL,  	 
	routeColor,  
	routeTextColor string ) (*Route, error) {
	if id == "" || routeShortName == "" || routeLongName == "" {
		return nil, ErrInvalidRoute
	}

	return &Route{
		id: id,     		 
		agencyId:agencyId,	     
		routeShortName: routeShortName,     
		routeLongName: routeLongName,	 
		routeDesc: routeDesc,    		   
		routeType: routeType,   	 
		routeURL:  routeURL,		 
		routeColor: routeColor, 	
		routeTextColor: routeTextColor, 
	}, nil
}

func RestoreRoute(id,
	agencyId,
	routeShortName,
	routeLongName, 
	routeDesc,     		  
	routeType,     
	routeURL,  	 
	routeColor,  
	routeTextColor string ) *Route {
	return &Route{
		id: id,     		 
		agencyId:agencyId,	     
		routeShortName: routeShortName,     
		routeLongName: routeLongName,	 
		routeDesc: routeDesc,    		   
		routeType: routeType,   	 
		routeURL:  routeURL,		 
		routeColor: routeColor, 	
		routeTextColor: routeTextColor, 
	}
}

// Métodos de acesso (getters) mantêm o encapsulamento
func (r * Route) ID() string { return r.id }
func (r * Route) AgencyID() string { return r.agencyId }
func (r * Route) ShortName() string { return r.routeShortName }
func (r * Route) LongName() string { return r.routeLongName }
func (r * Route) Desc() string { return r.routeDesc }
func (r * Route) Type() string { return r.routeType }
func (r * Route) RouteURL() string { return r.routeURL }
func (r * Route) RouteColor() string { return r.routeColor }
func (r * Route) RouteTextColor() string { return r.routeTextColor }