package worker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"uwo-tt-api/model"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NumberCol section number column index
const NumberCol = 0

// ComponentCol course component column index
const ComponentCol = 1

// ClassNbrCol section class number column index
const ClassNbrCol = 2

// DaysCol section scheduled days column index
const DaysCol = 3

// StartTimeCol section start time column index
const StartTimeCol = 4

// EndTimeCol section end time column index
const EndTimeCol = 5

// LocationCol section location column index
const LocationCol = 6

// InstructorCol section instructor column index
const InstructorCol = 7

// RequisitesCol section prerequisities column index
const RequisitesCol = 8

// StatusCol section status column index
const StatusCol = 9

// CampusCol section campus column index
const CampusCol = 10

// DeliveryCol section delivery type column index
const DeliveryCol = 11

type CourseComponent struct {
	Faculty     string
	Number      string
	Suffix      string
	Name        string
	Description string
}

type SectionComponent struct {
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

type CourseSection struct {
	Course  CourseComponent
	Section SectionComponent
}

// ScrapeCourses tbd
func ScrapeCourses(c chan *goquery.Document) {
	for doc := range c {
		courses := doc.Find(".span12")

		fmt.Println("New Doc")

		// Filter course list into each individual course table
		courses.ChildrenFiltered("table").Each(func(i int, course *goquery.Selection) {
			// Grab header based on course index
			header := courses.ChildrenFiltered("h4").Eq(i).Text()

			// Grab description based on course index
			desc := courses.ChildrenFiltered("p").Eq(i).Text()

			// Course name
			name := strings.Split(header, "-")[1]

			// Course faculty
			faculty := strings.Split(strings.Split(header, "-")[0], " ")[0]

			// Course number and suffix
			// 1000A -> "1000" "A"
			// 1000 -> "1000" ""
			suffix := ""
			number := strings.Split(strings.Split(header, "-")[0], " ")[1]
			if len(number) == 5 {
				suffix = string(number[4])
				number = number[:4]
			}

			fmt.Println(number, suffix)

			cc := CourseComponent{
				Faculty:     Trim(faculty),
				Number:      Trim(number),
				Suffix:      Trim(suffix),
				Name:        Trim(name),
				Description: Trim(desc)}

			// Filter course into each individual course section
			course.ChildrenFiltered("tbody").ChildrenFiltered("tr").Each(func(_ int, section *goquery.Selection) {

				var s SectionComponent

				// Filter section into each individual section column component
				section.ChildrenFiltered("td").Each(func(k int, elem *goquery.Selection) {
					// k represents index of the table heading; column number
					switch k {
					case NumberCol:
						s.Number = Trim(elem.Text())
					case ComponentCol:
						s.Component = Trim(elem.Text())
					case ClassNbrCol:
						s.ClassNumber = Trim(elem.Text())
					case DaysCol:
						// Find all valid table elements that represent days the class is scheduled for
						elem.Find("td").Each(func(d int, day *goquery.Selection) {
							if day.Text() != "&nbsp;" {
								s.Days += " " + Trim(day.Text())
							}
						})
						s.Days = Trim(s.Days)
					case StartTimeCol:
						s.StartTime = Trim(elem.Text())
					case EndTimeCol:
						s.EndTime = Trim(elem.Text())
					case LocationCol:
						s.Location = Trim(elem.Text())
					case InstructorCol:
						s.Instructor = Trim(elem.Text())
					case RequisitesCol:
						s.Reqs = Trim(elem.Text())
					case StatusCol:
						s.Status = Trim(elem.Text())
					case CampusCol:
						s.Campus = Trim(elem.Text())
					case DeliveryCol:
						s.Delivery = Trim(elem.Text())
					}
				})

				type Section struct {
					// ....
					Times []Time
					// ....
				}

				type Time struct {
					Day   string
					Start string
					End   string
				}
				// Check to see if section number and section component combination already exist
				// If exists, add day, start-time, end-time to days https://stackoverflow.com/questions/29817535/mongodb-how-to-insert-additional-object-into-object-collection-in-golang
				// If not exists, add

				// Find with query. If Find with query returns result, call Update with update filter that include $push. If find with query returns no result, insert document
			})
		})
	}
}

// // scrape data from course list div
// func scrapeCourses(courseList *goquery.Selection) []model.CourseData {
// 	var courses []model.CourseData

// 	// Grab each course in div; limited to direct children to avoid parsing tables further down the element list
// 	courseList.ChildrenFiltered("table").Each(func(i int, course *goquery.Selection) {

// 		// Get header and description based on course index
// 		headerElement := courseList.ChildrenFiltered("h4").Eq(i).Text()
// 		desc := courseList.ChildrenFiltered("p").Eq(i).Text()

// 		headerSplit := strings.Split(headerElement, "-")
// 		courseTitle := strings.Split(headerSplit[0], " ")

// 		// Define course properties
// 		c := model.CourseData{
// 			Faculty:     Trim(courseTitle[0]),
// 			Code:        Trim(courseTitle[1]),
// 			Name:        Trim(headerSplit[1]),
// 			Description: Trim(desc),
// 		}

// 		// Grab each section in course; limited to direct children to avoid parsing tbodys and trs further down the element list
// 		course.ChildrenFiltered("tbody").ChildrenFiltered("tr").Each(func(_ int, section *goquery.Selection) {

// 			var s model.SectionData

// 			// For each table heading of a section; limited to direct children to avoid parsing td further down the element list
// 			section.ChildrenFiltered("td").Each(func(k int, elem *goquery.Selection) {

// 				// k represents index of the table heading; column number
// 				switch k {
// 				case NumberCol:
// 					s.Number = Trim(elem.Text())
// 				case ComponentCol:
// 					s.Component = Trim(elem.Text())
// 				case ClassNbrCol:
// 					s.ClassNumber = Trim(elem.Text())
// 				case DaysCol:
// 					// Find all valid table elements that represent days the class is scheduled for
// 					elem.Find("td").Each(func(d int, day *goquery.Selection) {
// 						if day.Text() != "&nbsp;" {
// 							s.Days += " " + Trim(day.Text())
// 						}
// 					})
// 					s.Days = Trim(s.Days)
// 				case StartTimeCol:
// 					s.StartTime = Trim(elem.Text())
// 				case EndTimeCol:
// 					s.EndTime = Trim(elem.Text())
// 				case LocationCol:
// 					s.Location = Trim(elem.Text())
// 				case InstructorCol:
// 					s.Instructor = Trim(elem.Text())
// 				case RequisitesCol:
// 					s.Reqs = Trim(elem.Text())
// 				case StatusCol:
// 					s.Status = Trim(elem.Text())
// 				case CampusCol:
// 					s.Campus = Trim(elem.Text())
// 				case DeliveryCol:
// 					s.Delivery = Trim(elem.Text())
// 				}
// 			})

// 			// Add section to course
// 			c.AddSection(s)
// 		})

// 		// Add course to list of courses for subject
// 		courses = append(courses, c)
// 	})

// 	return courses
// }

// ScrapeTimeTable scraper
func ScrapeTimeTable(db *mongo.Database) {
	// Capture start time for metrics
	startTime := time.Now()

	// Create page to be scraped
	page := PageScraper{
		URL: "https://studentservices.uwo.ca/secure/timetables/mastertt/ttindex.cfm/",
		DB:  db,
	}

	// Fetch document synchronously
	doc, err := page.FetchDocument()
	if err != nil {
		fmt.Println("Error fetching document")
	}

	fmt.Println(page.Header, ": ", page.Status)

	// Find base search form
	page.Form = doc.Find("#searchForm")

	// Create a mapping of DB collections to html selectors where the scrape results of each selector will end up in the corresponding collection
	collectionToSelector := map[string]string{
		"subjects":     "#inputSubject",
		"suffixes":     "#inputDesignation",
		"course_types": "#inputCourseType",
		"components":   "#inputComponent",
		"campuses":     "#inputCampus",
		"start_times":  "[name=time]",
		"end_times":    "[name=end_time]",
	}

	fmt.Println("Options scraping:", time.Since(startTime))

	var wg sync.WaitGroup

	// Order is not preserved when looping over a map but order is not required in this case
	for key := range collectionToSelector {
		wg.Add(1)
		go page.ScrapeOptToDB(key, collectionToSelector[key], &wg)
	}

	// Wait for options to finish scraping to determine time improvments (~140ms -> ~20ms)
	wg.Wait()
	fmt.Println("Options scraping:", time.Since(startTime))

	startTime = time.Now()

	// COURSES SECTION
	// Grab available subjects from DB
	collection := page.DB.Collection("subjects")
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var subjects []model.Option
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		subjects = append(subjects, elem)
	}

	c := make(chan *goquery.Document)
	go ScrapeCourses(c)

	// Iterate over all subjects
	for _, subject := range subjects {
		if subject.Data.Value == "" {
			continue
		}

		// SleepRandom(10, 20)
		time.Sleep(time.Duration(10) * time.Second)

		// Post to page with subject. Needs to be done synchronously to keep time requirements
		data := CreateData(subject.Data.Value)
		doc, err := page.PostDocument(data)
		if err != nil {
			fmt.Println("Error fetching document")
		}

		c <- doc
	}

	close(c)
}

// // ScrapeTimeTable scrapes Western TimeTable for form options and course data
// func ScrapeTimeTable(db *mongo.Database) {
// 	// Capture start time
// 	startCapture := time.Now()

// 	// Create page to be scraped
// 	page := Page{
// 		URL: "https://studentservices.uwo.ca/secure/timetables/mastertt/ttindex.cfm/",
// 	}

// 	doc, err := page.FetchDocument()
// 	if err != nil {
// 		fmt.Println("Error fetching document")
// 	}

// 	fmt.Println(page.Header, ": ", page.Status)

// 	// Scrape options
// 	form := doc.Find("#searchForm")

// 	// Scrape options for each selector
// 	subjects := scrapeOptions(form, "#inputSubject")
// 	suffixes := scrapeOptions(form, "#inputDesignation")
// 	courseTypes := scrapeOptions(form, "#inputCourseType")
// 	components := scrapeOptions(form, "#inputComponent")
// 	campuses := scrapeOptions(form, "#inputCampus")

// 	// Time is scraped slightly differently
// 	startTimes := scrapeTimeOptions(form, "start")
// 	endTimes := scrapeTimeOptions(form, "end")

// 	// Store options in DB
// 	page.AddOptionsToDB(db, "subjects", subjects)
// 	page.AddOptionsToDB(db, "suffixes", suffixes)
// 	page.AddOptionsToDB(db, "course_types", courseTypes)
// 	page.AddOptionsToDB(db, "components", components)
// 	page.AddOptionsToDB(db, "campuses", campuses)
// 	page.AddOptionsToDB(db, "start_times", startTimes)
// 	page.AddOptionsToDB(db, "end_times", endTimes)

// 	// Improve randomness
// 	rand.Seed(time.Now().UnixNano())

// 	collection := db.Collection("courses")

// 	// Post site form for each available subject
// 	for _, subject := range subjects {
// 		if subject.Value == "" {
// 			continue
// 		}

// 		// Sleep for a random time between min and max to avoid getting blocked from site
// 		SleepRandom(10, 20)

// 		// Post to page with subject
// 		data := CreateData(subject.Value)
// 		doc, err := page.PostDocument(data)
// 		if err != nil {
// 			fmt.Println("Error fetching document")
// 		}

// 		fmt.Println(subject.Value, ": ", page.Status)

// 		// Scrape courses from relevant div
// 		courseList := doc.Find(".span12")
// 		courses := scrapeCourses(courseList)

// 		// Store courses from page into database
// 		page.AddCoursesToDB(collection, courses)
// 	}

// 	fmt.Println(time.Since(startCapture))
// }
