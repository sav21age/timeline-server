package domain

import (
	"time"
)

type ClubDTO struct {
	ID        uint16 `json:"id" bson:"_id"`
	ShortName string `json:"shortName" bson:"shortName"`

	Name  string   `json:"name" bson:"name"`
	Crest string   `json:"crest" bson:"crest"`
	Area  AreaClub `json:"area" bson:"area"`
	// RunningCompetitions []RunningCompetitions `json:"runningCompetitions" bson:"runningCompetitions"`
}

type Club struct {
	ID        uint16 `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name"`
	ShortName string `json:"shortName" bson:"shortName"`
	Tla       string `json:"tla" bson:"tla"`
	Crest     string `json:"crest" bson:"crest"`
	Address   string `json:"address" bson:"address"`
	// Phone               string                `json:"phone" bson:"phone"`
	Website string `json:"website" bson:"website"`
	// Email               string                `json:"email" bson:"email"`
	Founded             uint16                `json:"founded" bson:"founded"`
	ClubColors          string                `json:"clubColors" bson:"clubColors"`
	Venue               string                `json:"venue" bson:"venue"`
	Area                AreaClub              `json:"area" bson:"area"`
	Coach               Coach                 `json:"coach" bson:"coach"`
	RunningCompetitions []RunningCompetitions `json:"runningCompetitions" bson:"runningCompetitions"`
	Squad               []Squad               `json:"squad" bson:"squad"`
}

type RunningCompetitions struct {
	ID     uint16 `json:"id" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Code   string `json:"code" bson:"code"`
	Type   string `json:"type" bson:"type"`
	Emblem string `json:"emblem" bson:"emblem"`

	// Area AreaClub `json:"area" bson:"area"`
}

type AreaClub struct {
	ID   uint16 `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
	Flag string `json:"flag" bson:"flag"`
}

type Coach struct {
	ID          uint32     `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName   string     `json:"firstName,omitempty" bson:"firstName"`
	LastName    string     `json:"lastName,omitempty" bson:"lastName"`
	Name        string     `json:"name,omitempty" bson:"name"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty" bson:"dateOfBirth"`
	Nationality string     `json:"nationality,omitempty" bson:"nationality"`
}

type Squad struct {
	ID          uint32     `json:"id" bson:"_id,omitempty"`
	Name        string     `json:"name" bson:"name"`
	Position    string     `json:"position" bson:"position"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty" bson:"dateOfBirth"`
	// CountryOfBirth string    `json:"countryOfBirth" bson:"countryOfBirth"`
	Nationality string `json:"nationality" bson:"nationality"`
	// ShirtNumber    uint8     `json:"shirtNumber" bson:"shirtNumber"`
	// Role           string    `json:"role" bson:"role"`
}

//--

type GetClubsQueryParams struct {
	Pagination
	AreaId int
	SortBy string
}
