package model

import (
	"time"
)

// SourceInfo struct to define source of information; URL and header string
type SourceInfo struct {
	Title string `bson:"title" json:"title" example:"Fall/Winter Academic Timetable"`
	Year  string `bson:"year" json:"year" example:"2020/2021"`
	URL   string `bson:"url" json:"url" example:"https://studentservices.uwo.ca/secure/timetables/mastertt/ttindex.cfm"`
}

// TimeInfo struct to define time
type TimeInfo struct {
	Added time.Time `bson:"added" json:"added"`
}
