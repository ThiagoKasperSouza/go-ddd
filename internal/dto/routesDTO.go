package dto


type RouteDTO struct {
	RouteID       			string `json:"route_id"`
	AgencyID       			string `json:"agency_id"`
	RouteShortName      	string `json:"route_short_name"`
	RouteLongName			string `json:"route_long_name"`
	RouteDesc     		   	string `json:"route_desc"`
	RouteType    			string `json:"route_type"`
	RouteURL  				string `json:"route_url"`
	RouteColor  			string `json:"route_color"`
	RouteTextColor  		string `json:"route_text_color"`
}

type ListRoutesOutputDTO struct {
	Routes []RouteDTO `json:"routes"`
}