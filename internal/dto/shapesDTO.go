package dto


type ShapeDTO struct {
	ShapeID       			string `json:"shape_id"`
	ShapePtLat       		string `json:"shape_pt_lat"`
	ShapePtLon     		  	string `json:"shape_pt_lon"`
	ShapePtSequence			string `json:"shape_pt_sequence"`
}

type ListShapesOutputDTO struct {
	Shapes []ShapeDTO `json:"shapes"`
}