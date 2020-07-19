package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"uwo-tt-api/model"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	go page.ScrapeCoursesToDB(c)

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

	fmt.Println("Course scraping:", time.Since(startTime))

}
