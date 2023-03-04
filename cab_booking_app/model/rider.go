package model

type Rider struct {
	Id       int
	Location Location
}

func NewRider(x, y int) Rider {
	r := Rider{Location: Location{X: x, Y: y}, Id: -1}
	return r
}

func (r *Rider) GetLocation() Location {
	return r.Location
}

func (r *Rider) SetId(id int) {
	r.Id = id
}
