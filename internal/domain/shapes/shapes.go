package shapes

import "errors"

var (
	ErrShapeNotFound = errors.New("Shape not found")
	ErrInvalidShape  = errors.New("invalid Shape data")
)

type Shape struct {
	id       			string 
	shapePtLat       	string 
	shapePtLon     		string 
	shapePtSequence		string
}

// NewAgency garante que uma nova agência seja criada com dados válidos
func NewShape(id,
	shapePtLat,
	shapePtLon,
	shapePtSequence	string) (*Shape, error) {
	if id == "" || shapePtLat == "" || shapePtLon == "" {
		return nil, ErrInvalidShape
	}

	return &Shape{
			id:id,
			shapePtLat:shapePtLat,
			shapePtLon:shapePtLon,
			shapePtSequence:shapePtSequence,
		}, nil
}

func RestoreShape(id,
	shapePtLat,
	shapePtLon,
	shapePtSequence	string) *Shape {
	return &Shape{
			id:id,
			shapePtLat:shapePtLat,
			shapePtLon:shapePtLon,
			shapePtSequence:shapePtSequence,
		}
}

// Métodos de acesso (getters) mantêm o encapsulamento
func (s * Shape) ID() string { return s.id }
func (s * Shape) ShapePtLat() string { return s.shapePtLat }
func (s * Shape) ShapePtLon() string { return s.shapePtLon }
func (s * Shape) ShapePtSequence() string { return s.shapePtSequence }