package model

type Cab struct {
	Id          int
	Location    Location
	isAvailable bool
}

func NewCab(x, y int) Cab {
	c := Cab{Location: Location{X: x, Y: y}, isAvailable: true, Id: -1}
	return c
}

func (c *Cab) GetLocation() Location {
	return c.Location
}

func (c *Cab) UpdateLocation(x, y int) {
	c.Location = Location{x, y}
}

func (c *Cab) IsAvailable() bool {
	return c.isAvailable
}

func (c *Cab) SetIsAvailable(val bool) {
	c.isAvailable = true
}

func (c *Cab) SetId(id int) {
	c.Id = id
}
