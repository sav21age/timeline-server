package domain

import (
	"time"
)

type Season struct {
	ID              uint16            `json:"id" bson:"_id"`
	StartDate       time.Time         `json:"startDate" bson:"startDate"`
	EndDate         time.Time         `json:"endDate" bson:"endDate"`
	CurrentMatchday uint8             `json:"currentMatchday" bson:"currentMatchday"`
	Winner          *WinnerSeason     `json:"winner" bson:"winner"`
	Competition     CompetitionSeason `json:"competition" bson:"competition"`
	Available       bool              `json:"available" bson:"available,omitempty"`
}

type CompetitionSeason struct {
	ID     uint16     `json:"id" bson:"_id,omitempty"`
	Name   string     `json:"name" bson:"name"`
	Code   string     `json:"code" bson:"code"`
	Emblem string     `json:"emblem" bson:"emblem"`
	Area   AreaSeason `json:"area" bson:"area"`
}

type AreaSeason struct {
	ID          uint16 `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Ensign      string `json:"ensign" bson:"ensign"`
}

type WinnerSeason struct {
	ID        uint16 `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string `json:"name,omitempty" bson:"name"`
	ShortName string `json:"shortName,omitempty" bson:"shortName"`
	Tla       string `json:"tla,omitempty" bson:"tla"`
	Crest     string `json:"crest,omitempty" bson:"crest"`
}

//--

type GetSeasonsQueryParams struct {
	Pagination
}
