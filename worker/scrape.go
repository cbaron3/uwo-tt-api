package worker

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"uwo-tt-api/model"

	"github.com/PuerkitoBio/goquery"
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

// scrape data from options field define by selector string
func scrapeOptions(form *goquery.Selection, selector string) []model.OptionData {
	var options []model.OptionData

	// Find all instances of selector and append it to result
	form.Find(selector + " option").Each(func(i int, elem *goquery.Selection) {
		text := elem.Text()
		value, _ := elem.Attr("value")

		options = append(options, model.OptionData{Value: Trim(value), Text: Trim(text)})
	})

	return options
}

// scrape data from time option; either start or stop
func scrapeTimeOptions(form *goquery.Selection, timeType string) []model.OptionData {
	var options []model.OptionData

	// Find all instances of time selector
	form.Find("#inputTime").Each(func(i int, elem *goquery.Selection) {
		// Grab name of selector
		name, _ := elem.Attr("name")

		if name == "time" && timeType == "start" {
			// If the name of the element is TIME and we want START time types, collect all options
			options = scrapeOptions(elem, "")

		} else if name == "end_time" && timeType == "end" {
			// If the name of the element is TIME and we want START time types, collect all options
			options = scrapeOptions(elem, "")
		}
	})

	return options
}

// scrape data from course list div
func scrapeCourses(courseList *goquery.Selection) []model.CourseData {
	var courses []model.CourseData

	// Grab each course in div; limited to direct children to avoid parsing tables further down the element list
	courseList.ChildrenFiltered("table").Each(func(i int, course *goquery.Selection) {

		// Get header and description based on course index
		headerElement := courseList.ChildrenFiltered("h4").Eq(i).Text()
		desc := courseList.ChildrenFiltered("p").Eq(i).Text()

		headerSplit := strings.Split(headerElement, "-")
		courseTitle := strings.Split(headerSplit[0], " ")

		// Define course properties
		c := model.CourseData{
			Faculty:     Trim(courseTitle[0]),
			Code:        Trim(courseTitle[1]),
			Name:        Trim(headerSplit[1]),
			Description: Trim(desc),
		}

		// Grab each section in course; limited to direct children to avoid parsing tbodys and trs further down the element list
		course.ChildrenFiltered("tbody").ChildrenFiltered("tr").Each(func(_ int, section *goquery.Selection) {

			var s model.SectionData

			// For each table heading of a section; limited to direct children to avoid parsing td further down the element list
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

			// Add section to course
			c.AddSection(s)
		})

		// Add course to list of courses for subject
		courses = append(courses, c)
	})

	return courses
}

// ScrapeTimeTable scrapes Western TimeTable for form options and course data
func ScrapeTimeTable(db *mongo.Database) {
	// Capture start time
	startCapture := time.Now()

	// Create page to be scraped
	page := Page{
		URL: "https://studentservices.uwo.ca/secure/timetables/mastertt/ttindex.cfm/",
	}

	doc, err := page.FetchDocument()
	if err != nil {
		fmt.Println("Error fetching document")
	}

	fmt.Println(page.Header, ": ", page.Status)

	// Scrape options
	form := doc.Find("#searchForm")

	// Scrape options for each selector
	subjects := scrapeOptions(form, "#inputSubject")
	suffixes := scrapeOptions(form, "#inputDesignation")
	courseTypes := scrapeOptions(form, "#inputCourseType")
	components := scrapeOptions(form, "#inputComponent")
	campuses := scrapeOptions(form, "#inputCampus")

	// Time is scraped slightly differently
	startTimes := scrapeTimeOptions(form, "start")
	endTimes := scrapeTimeOptions(form, "end")

	// Store options in DB
	page.AddOptionsToDB(db, "subjects", subjects)
	page.AddOptionsToDB(db, "suffixes", suffixes)
	page.AddOptionsToDB(db, "course_types", courseTypes)
	page.AddOptionsToDB(db, "components", components)
	page.AddOptionsToDB(db, "campuses", campuses)
	page.AddOptionsToDB(db, "start_times", startTimes)
	page.AddOptionsToDB(db, "end_times", endTimes)

	// Improve randomness
	rand.Seed(time.Now().UnixNano())

	collection := db.Collection("courses")

	// Post site form for each available subject
	for _, subject := range subjects {
		if subject.Value == "" {
			continue
		}

		// Sleep for a random time between min and max to avoid getting blocked from site
		SleepRandom(10, 20)

		// Post to page with subject
		data := CreateData(subject.Value)
		doc, err := page.PostDocument(data)
		if err != nil {
			fmt.Println("Error fetching document")
		}

		fmt.Println(subject.Value, ": ", page.Status)

		// Scrape courses from relevant div
		courseList := doc.Find(".span12")
		courses := scrapeCourses(courseList)

		// Store courses from page into database
		page.AddCoursesToDB(collection, courses)
	}

	fmt.Println(time.Since(startCapture))
}
