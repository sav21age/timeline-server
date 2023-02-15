package domain

import (
	"time"
)

type MatchDTO struct {
	ID          uint                  `json:"id" bson:"_id"`
	Competition CompetitionMatchShort `json:"competition" bson:"competition"`
	Season      SeasonMatchShort      `json:"season" bson:"season"`
	UtcDate     time.Time             `json:"utcDate" bson:"utcDate"`
	Status      string                `json:"status" bson:"status"`
	Matchday    uint8                 `json:"matchday" bson:"matchday"`
	Score       ScoreMatchShort       `json:"score" bson:"score"`
	HomeTeam    TeamShortMatch        `json:"homeTeam" bson:"homeTeam"`
	AwayTeam    TeamShortMatch        `json:"awayTeam" bson:"awayTeam"`
}

type Match struct {
	ID          uint             `json:"id" bson:"_id,omitempty"`
	Head2Head   Head2Head        `json:"head2head" bson:"head2head"`
	Competition CompetitionMatch `json:"competition,omitempty" bson:"competition"`
	Season      SeasonMatch      `json:"season,omitempty" bson:"season"`
	UtcDate     time.Time        `json:"utcDate" bson:"utcDate"`
	Status      string           `json:"status,omitempty" bson:"status"`
	Venue       string           `json:"venue,omitempty" bson:"venue"`
	Matchday    uint8            `json:"matchday,omitempty" bson:"matchday"`
	Stage       string           `json:"stage,omitempty" bson:"stage"`
	Group       string           `json:"group,omitempty" bson:"group"`
	Score       ScoreMatch       `json:"score,omitempty" bson:"score"`
	HomeTeam    TeamShortMatch   `json:"homeTeam,omitempty" bson:"homeTeam"`
	AwayTeam    TeamShortMatch   `json:"awayTeam,omitempty" bson:"awayTeam"`
	Referees    *[]RefereeMatch  `json:"referees,omitempty" bson:"referees,omitempty"`
}

type Head2Head struct {
	NumberOfMatches uint16    `json:"numberOfMatches" bson:"numberOfMatches"`
	TotalGoals      uint8     `json:"totalGoals,omitempty" bson:"totalGoals"`
	HomeTeam        TeamMatch `json:"homeTeam,omitempty" bson:"homeTeam"`
	AwayTeam        TeamMatch `json:"awayTeam,omitempty" bson:"awayTeam"`
}

type TeamMatch struct {
	ID     uint   `json:"id" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Wins   uint8  `json:"wins" bson:"wins"`
	Draws  uint8  `json:"draws" bson:"draws"`
	Losses uint8  `json:"losses" bson:"losses"`
}

type CompetitionMatchShort struct {
	ID   uint16    `json:"id" bson:"_id,omitempty"`
	Name string    `json:"name" bson:"name"`
	Area AreaMatch `json:"area" bson:"area"`
}

type CompetitionMatch struct {
	ID   uint16    `json:"id" bson:"_id,omitempty"`
	Name string    `json:"name" bson:"name"`
	Area AreaMatch `json:"area" bson:"area"`
}

type AreaMatch struct {
	Name        string `json:"name" bson:"name"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Ensign      string `json:"ensign" bson:"ensign"`
}

type SeasonMatchShort struct {
	ID        uint16    `json:"id" bson:"_id,omitempty"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
}

type SeasonMatch struct {
	ID              uint16       `json:"id" bson:"_id,omitempty"`
	StartDate       time.Time    `json:"startDate" bson:"startDate"`
	EndDate         time.Time    `json:"endDate" bson:"endDate"`
	CurrentMatchday uint8        `json:"currentMatchday" bson:"currentMatchday"`
	Winner          *WinnerMatch `json:"winner" bson:"winner"`
}

type WinnerMatch struct {
	ID        uint16 `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string `json:"name,omitempty" bson:"name"`
	ShortName string `json:"shortName,omitempty" bson:"shortName"`
	Tla       string `json:"tla,omitempty" bson:"tla"`
	Crest     string `json:"crest,omitempty" bson:"crest"`
}

type ScoreMatchShort struct {
	Winner   string     `json:"winner,omitempty" bson:"winner"`
	FullTime MatchScore `json:"fullTime,omitempty" bson:"fullTime"`
}

type ScoreMatch struct {
	Winner    string     `json:"winner,omitempty" bson:"winner"`
	Duration  string     `json:"duration,omitempty" bson:"duration"`
	FullTime  MatchScore `json:"fullTime,omitempty" bson:"fullTime"`
	HalfTime  MatchScore `json:"halfTime,omitempty" bson:"halfTime"`
	ExtraTime MatchScore `json:"extraTime,omitempty" bson:"extraTime"`
	Penalties MatchScore `json:"penalties,omitempty" bson:"penalties"`
}

type MatchScore struct {
	HomeTeam *uint8 `json:"homeTeam,omitempty" bson:"homeTeam"`
	AwayTeam *uint8 `json:"awayTeam,omitempty" bson:"awayTeam"`
}

type TeamShortMatch struct {
	ID   uint16 `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name"`
}

type RefereeMatch struct {
	ID          uint   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string `json:"name,omitempty" bson:"name"`
	Role        string `json:"role" bson:"role"`
	Nationality string `json:"nationality" bson:"nationality"`
}

//--

type GetMatchesQueryParams struct {
	Datination  Datination
	ClubId      int
	MatchStatus string
}

type GetDatesQueryParams struct {
	ClubId      int
	MatchStatus string
}
