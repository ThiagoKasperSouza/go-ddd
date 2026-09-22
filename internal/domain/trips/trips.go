package trips

import "errors"

var (
	ErrTripNotFound = errors.New("Trip not found")
	ErrInvalidTrip  = errors.New("invalid Trip data")
)

type Trip struct {
	routeID		   			string
	serviceID       		string
	tripID     		  		string
	tripHeadsign			string 
	tripShortName      		string 
	directionID       		string
	blockID     		  	string 
	shapeID					string
	wheelchairAccessible	string
}

// NewAgency garante que uma nova agência seja criada com dados válidos
func NewTrip(routeID,
	serviceID,
	tripID,
	tripHeadsign, 
	tripShortName, 
	directionID,
	blockID, 
	shapeID,
	wheelchairAccessible string) (*Trip, error) {
	if tripID == "" || serviceID == "" || shapeID == "" {
		return nil, ErrInvalidTrip
	}

	return &Trip{
			routeID:routeID,			
			serviceID:serviceID,
			tripID:tripID,
			tripHeadsign:tripHeadsign,
			tripShortName:tripShortName,
			directionID:directionID,
			blockID: blockID,
			shapeID:shapeID,
			wheelchairAccessible:wheelchairAccessible,
		}, nil
}

func RestoreTrip(routeID,
	serviceID,
	tripID,
	tripHeadsign, 
	tripShortName, 
	directionID,
	blockID, 
	shapeID,
	wheelchairAccessible string) *Trip {
	return &Trip{
			routeID:routeID,			
			serviceID:serviceID,
			tripID:tripID,
			tripHeadsign:tripHeadsign,
			tripShortName:tripShortName,
			directionID:directionID,
			blockID: blockID,
			shapeID:shapeID,
			wheelchairAccessible:wheelchairAccessible,
		}
}

// Métodos de acesso (getters) mantêm o encapsulamento
func (t * Trip) RouteID() string { return t.routeID }
func (t * Trip) ServiceID() string { return t.serviceID }
func (t * Trip) TripID() string { return t.tripID }
func (t * Trip) TripHeadsign() string { return t.tripHeadsign }
func (t * Trip) TripShortName() string { return t.tripShortName }
func (t * Trip) DirectionID() string { return t.directionID }
func (t * Trip) BlockID() string { return t.blockID }
func (t * Trip) ShapeID() string { return t.shapeID }
func (t * Trip) WheelchairAccessible() string { return t.wheelchairAccessible }
