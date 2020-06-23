package model

// Section struct
type Section struct {
	Number      string ``
	Component   string ``
	ClassNumber string ``
	Days        string ``
	StartTime   string ``
	EndTime     string ``
	Location    string ``
	Instructor  string ``
	Reqs        string ``
	Status      string ``
	Campus      string ``
	Delivery    string ``
	Open        bool   ``
}

// Course struct
type Course struct {
	Faculty     string    ``
	Code        string    ``
	Name        string    ``
	Description string    ``
	sections    []Section ``
}
