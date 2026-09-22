package dto


type TripDTO struct {
	RouteID       			string `json:"route_id"`
	ServiceID       		string `json:"service_id"`
	TripID     		  		string `json:"trip_id"`
	TripHeadsign			string `json:"trip_headsign"`
	TripShortName      		string `json:"trip_shortname"`
	DirectionID       		string `json:"direction_id"`
	BlockID     		  	string `json:"block_id"`
	ShapeID					string `json:"shape_id"`
	WheelchairAccessible	string `json:"wheelchair_accessible"`
}


type ListTripsOutputDTO struct {
	Trips []TripDTO `json:"trips"`
}