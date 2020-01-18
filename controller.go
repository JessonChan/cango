package cango

type URI interface {
}

type Controller struct {
	URI  `value:"/cinema/{cinemaId}/movie/{movieId}"`
	Name string
}

type SugarController struct {
	*Controller
}

func (s *SugarController) Get(param struct {
	URI      `value:"/people/{peopleId}.json"`
	CinemaId int `nil:"false"`
	MovieId  int
	PeopleId int
}) {
	if s.URI == nil {
		return
	}
	if param.URI == nil {
		return
	}
	if param.CinemaId == 0 {
		return
	}
	if param.MovieId == 0 {
	}
	return
}
