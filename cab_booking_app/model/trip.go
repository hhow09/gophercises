package model

type TRIP_STATUS int

const (
	IN_PROGRESS TRIP_STATUS = iota + 1
	FINISHED
)

type Trip struct {
	status      TRIP_STATUS
	rider       Rider
	cab         Cab
	destination Location
}

func (t *Trip) SetStatus(s TRIP_STATUS) {
	t.status = s
}

func (t *Trip) GetRider() Rider {
	return t.rider
}

func NewTrip(rider Rider, cab Cab, destination Location) Trip {
	t := Trip{status: IN_PROGRESS, rider: rider, cab: cab, destination: destination}
	return t
}
