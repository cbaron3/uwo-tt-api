package worker

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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

// PageScraper defines the context for the page to be scraped and the location of scrape resulst
type PageScraper struct {
	Header string
	URL    string
	Status string
	DB     *mongo.Database
	Form   *goquery.Selection
}

// PageResult encompasses data that is passed into channel to be parsed
type PageResult struct {
	Name string
	Doc  *goquery.Document
}

// BuildSourceInfo creates source info based on page information
func (page *PageScraper) BuildSourceInfo() model.SourceInfo {
	sourceInfo := model.SourceInfo{
		Title: page.Header[0 : len(page.Header)-9], // Everything besides last 9 chars defines the Title
		Year:  page.Header[len(page.Header)-9:],    // Last 9 chars define the Year
		URL:   page.URL,
	}

	return sourceInfo
}

// BuildTimeInfo creates time info based on current time
func (page *PageScraper) BuildTimeInfo() model.TimeInfo {
	timeInfo := model.TimeInfo{
		Added: time.Now(),
	}

	return timeInfo
}

// FetchDocument fetches contents of page based on URL
func (page *PageScraper) FetchDocument() (document *goquery.Document, err error) {
	// GET request for website contents
	resp, err := http.Get(page.URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// goquery for parsing
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Grab page status and header information
	page.Status = resp.Status

	header := doc.Find(".page-header h1 small")
	page.Header = Trim(header.Text())

	return doc, nil
}

// PostDocument fetches content from page URL after submitting post request with form data
func (page *PageScraper) PostDocument(data map[string][]string) (document *goquery.Document, err error) {
	// POST request for website contents
	resp, err := http.PostForm(page.URL, data)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// goquery for content parsing
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Grab page status
	page.Status = resp.Status

	return doc, nil
}

// ScrapeOptToDB temp
func (page *PageScraper) ScrapeOptToDB(collectionName string, selector string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Connect to temporary collection
	tempCollection := page.DB.Collection(collectionName + "_temp")

	// Delete all existing documents if recovering from crash or shutdown
	res, err := tempCollection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("DeleteMany() ERROR:", err)
		return
	}

	fmt.Printf("Cleared %d documents from temporary %s collection\n", res.DeletedCount, collectionName)

	// Scraped options, empty map slice
	opts := []interface{}{}

	// Find all instances of selector and append it to result
	page.Form.Find(selector + " option").Each(func(i int, elem *goquery.Selection) {
		text := elem.Text()
		value, _ := elem.Attr("value")

		opts = append(opts, model.Option{
			Source: page.BuildSourceInfo(),
			Time:   page.BuildTimeInfo(),
			Data: model.OptionData{
				Value: Trim(value),
				Text:  Trim(text),
			},
		})
	})

	// Insert data into a new temporary collection to not affect main collection during scraping
	insertCtx, err := tempCollection.InsertMany(context.TODO(), opts)
	if err != nil {
		fmt.Println(err)
	}

	// Create aggregation pipeline were first, all data is matched, then all data is written to collectionName
	// $out replaces data in collection
	pipeline := bson.A{
		bson.M{"$match": bson.M{}},
		bson.M{"$out": collectionName},
	}

	// Use aggregation to make sure main collection is never empty
	startTime := time.Now()
	tempCollection.Aggregate(context.TODO(), pipeline)

	fmt.Printf("%s Inserted: %d. Aggregation time: %s\n",
		collectionName,
		len(insertCtx.InsertedIDs),
		time.Since(startTime).String())

	// Drop temporary collection
	err = tempCollection.Drop(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

}

// extract information course specific information from goquery selection
func extractCourseInfo(courseList *goquery.Selection, courseIndex int) model.CourseComponent {
	// Grab header based on course index
	header := courseList.ChildrenFiltered("h4").Eq(courseIndex).Text()

	// Grab description based on course index
	desc := courseList.ChildrenFiltered("p").Eq(courseIndex).Text()

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

	num, err := strconv.Atoi(Trim(number))
	if err != nil {
		// handle error
		fmt.Println(err)
		num = 0
	}

	return model.CourseComponent{
		Faculty:     Trim(faculty),
		Number:      num,
		Suffix:      Trim(suffix),
		Name:        Trim(name),
		Description: Trim(desc)}
}

// extract section specific information from goquery selection
func extractSectionInfo(section *goquery.Selection) model.SectionComponent {
	var s model.SectionComponent

	// Filter section into each individual section column component
	var start string
	var end string
	var days []string

	section.ChildrenFiltered("td").Each(func(k int, elem *goquery.Selection) {

		// k represents index of the table heading; column number
		switch k {
		case NumberCol:
			i, err := strconv.Atoi(Trim(elem.Text()))
			if err != nil {
				// handle error
				fmt.Println(err)
			} else {
				s.Number = i
			}
		case ComponentCol:
			s.Component = Trim(elem.Text())
		case ClassNbrCol:
			i, err := strconv.Atoi(Trim(elem.Text()))
			if err != nil {
				// handle error
				fmt.Println(err)
			} else {
				s.ClassNumber = i
			}
		case DaysCol:
			// Find all valid table elements that represent days the class is scheduled for
			elem.Find("td").Each(func(d int, day *goquery.Selection) {
				if day.Text() != "&nbsp;" {
					days = append(days, Trim(day.Text()))
				}
			})
			// s.Days = Trim(s.Days)
		case StartTimeCol:
			start = Trim(elem.Text())
		case EndTimeCol:
			end = Trim(elem.Text())
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

	// Collect days and times
	for _, day := range days {
		if day != "" {
			s.Times = append(s.Times, model.TimeComponent{day, start, end})
		}
	}

	return s
}

// ScrapeCoursesToDB scrapes course information from pages incoming into channel and store info in database
func (page *PageScraper) ScrapeCoursesToDB(c chan PageResult, size int) {

	// Connect to temporary collection
	tempCollection := page.DB.Collection("courses_temp")

	// Delete all existing documents if recovering from crash
	res, err := tempCollection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("DeleteMany() ERROR:", err)
		return
	}

	fmt.Printf("Cleared %d documents from temporary courses collection\n", res.DeletedCount)

	counter := 1
	// Iterate over documents in the channel as long as the channel is open
	for doc := range c {
		courses := doc.Doc.Find(".span12")

		// Size-1 to account for the ANY element
		fmt.Printf("Scraping - #%d: %s - %.2f%%\n", counter, doc.Name, float32((float32(counter)/float32(size-1))*100.00))
		counter++

		// Filter course list into each individual course table
		courses.ChildrenFiltered("table").Each(func(i int, course *goquery.Selection) {

			// Course info and section info are not grouped into a div so need to match table to header/p with index
			courseData := extractCourseInfo(courses, i)

			// Filter course into each individual course section
			course.ChildrenFiltered("tbody").ChildrenFiltered("tr").Each(func(_ int, section *goquery.Selection) {

				sectionData := extractSectionInfo(section)

				courseSection := model.Section{
					Source:      page.BuildSourceInfo(),
					Time:        page.BuildTimeInfo(),
					CourseData:  courseData,
					SectionData: sectionData,
				}

				// Update all components with the section information and component by adding new time information
				query := bson.M{"$and": bson.A{
					bson.M{"courseData": courseData},
					bson.M{"sectionData.number": sectionData.Number},
					bson.M{"sectionData.component": sectionData.Component}}}

				changes := bson.M{"$push": bson.M{"sectionData.times": bson.M{"$each": sectionData.Times}}}

				// Update
				updateResult, err := tempCollection.UpdateMany(context.TODO(), query, changes)
				if err != nil {
					if len(sectionData.Times) != 0 {
						fmt.Println(err)
					}
				}

				// If no modifications were made, insert the document
				if updateResult.ModifiedCount == 0 {
					_, insertErr := tempCollection.InsertOne(context.TODO(), courseSection)
					if insertErr != nil {
						fmt.Println(insertErr)
					}
				}
			})
		})
	}

	// TODO: Refractor out because used twice
	// Create aggregation pipeline where first, all data is matched, then all data is written to collectionName
	pipeline := bson.A{
		bson.M{"$match": bson.M{}},
		bson.M{"$out": "courses"},
	}

	// Use aggregation to make sure main collection is never empty
	startTime := time.Now()
	tempCollection.Aggregate(context.TODO(), pipeline)

	fmt.Printf("Course aggregation time: %s\n", time.Since(startTime).String())

	// Drop temporary collection
	err = tempCollection.Drop(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
}
