package model

// SectionData struct
type SectionData struct {
	Number      string `bson:"number" 		json:"number" 		example:"001"`
	Component   string `bson:"component" 	json:"component" 	example:"LEC"`
	ClassNumber string `bson:"classNumber" 	json:"classNumber" 	example:"5000"`
	Days        string `bson:"days" 		json:"days" 		example:"M W F"`
	StartTime   string `bson:"startTime" 	json:"startTime" 	example:"8:30 AM"`
	EndTime     string `bson:"endTime" 		json:"endTime" 		example:"1:30 PM"`
	Location    string `bson:"location" 	json:"location" 	example:"NS 145"`
	Instructor  string `bson:"instructor" 	json:"instructor" 	example:"Haffie"`
	Reqs        string `bson:"requisites" 	json:"requisites" 	example:"REQUISITES:..."`
	Status      string `bson:"status" 		json:"status" 		example:"Full"`
	Campus      string `bson:"campus" 		json:"campus" 		example:"Main"`
	Delivery    string `bson:"delivery" 	json:"delivery" 	example:"Distance Studies/Online"`
}

// CourseData struct
type CourseData struct {
	Faculty     string        `bson:"faculty" 		json:"faculty" 		example:"CLASSICS"`
	Code        string        `bson:"code" 			json:"code" 		example:"2053B"`
	Name        string        `bson:"name" 			json:"name" 		example:"MATH FOR FINANCIAL ANALYSIS"`
	Description string        `bson:"description" 	json:"description" 	example:"Course description"`
	Sections    []SectionData `bson:"sections" 		json:"sections"`
}

// Course struct to define a course to be stored in database
type Course struct {
	Source SourceInfo `bson:"source" json:"source"`
	Time   TimeInfo   `bson:"time" json:"time"`
	Data   CourseData `bson:"data" json:"data"`
}

// AddSection blank
func (course *CourseData) AddSection(item SectionData) []SectionData {
	course.Sections = append(course.Sections, item)
	return course.Sections
}
