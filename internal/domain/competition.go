package domain

import (
	"time"
)

type Competition struct {
	ID     uint16            `json:"id" bson:"_id"`
	Name   string            `json:"name" bson:"name"`
	Code   string            `json:"code" bson:"code"`
	Emblem string            `json:"emblem" bson:"emblem"`
	Area   AreaCompetition   `json:"area" bson:"area"`
	Season SeasonCompetition `json:"season" bson:"season"`
}

type AreaCompetition struct {
	ID          uint16 `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Ensign      string `json:"ensign" bson:"ensign"`
}

type SeasonCompetition struct {
	ID              uint16             `json:"id" bson:"_id,omitempty"`
	StartDate       time.Time          `json:"startDate" bson:"startDate"`
	EndDate         time.Time          `json:"endDate" bson:"endDate"`
	CurrentMatchday uint8              `json:"currentMatchday" bson:"currentMatchday"`
	Winner          *WinnerCompetition `json:"winner" bson:"winner"`
}

type WinnerCompetition struct {
	ID        uint16 `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string `json:"name,omitempty" bson:"name"`
	ShortName string `json:"shortName,omitempty" bson:"shortName"`
	Tla       string `json:"tla,omitempty" bson:"tla"`
	Crest     string `json:"crest,omitempty" bson:"crest"`
}
