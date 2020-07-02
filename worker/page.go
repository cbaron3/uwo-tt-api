package worker

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"uwo-tt-api/model"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Page defines the website source
type Page struct {
	Header string
	URL    string
	Status string
}

// BuildSourceInfo creates source info based on page information
func (page *Page) BuildSourceInfo() model.SourceInfo {
	sourceInfo := model.SourceInfo{
		Title: page.Header[0 : len(page.Header)-9], // Everything besides last 9 chars defines the Title
		Year:  page.Header[len(page.Header)-9:],    // Last 9 chars define the Year
		URL:   page.URL,
	}

	return sourceInfo
}

// BuildTimeInfo creates time info based on current time
func (page *Page) BuildTimeInfo() model.TimeInfo {
	timeInfo := model.TimeInfo{
		Added: time.Now(),
	}

	return timeInfo
}

// FetchDocument fetches contents of page based on URL
func (page *Page) FetchDocument() (document *goquery.Document, err error) {
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
func (page *Page) PostDocument(data map[string][]string) (document *goquery.Document, err error) {
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

// AddCoursesToDB stores course information that was sourced from the page in mongodb
func (page *Page) AddCoursesToDB(collection *mongo.Collection, courses []model.CourseData) {
	if len(courses) <= 0 {
		fmt.Println("Insufficient course data")
		return
	}

	newCourses := []interface{}{}

	// Convert list of course data into type compatible with mongo insertions (list of empty interface)
	// Add extra metadata to data as well
	for _, course := range courses {
		newCourses = append(newCourses, model.Course{
			Source: page.BuildSourceInfo(),
			Time:   page.BuildTimeInfo(),
			Data:   course,
		})
	}

	// Delete all other courses of the same faculty first to ensure course data is up to date
	deleteCtx, err := collection.DeleteMany(context.TODO(), bson.M{"data.faculty": courses[0].Faculty})

	if err != nil {
		fmt.Println(err)
	}

	// Insert new courses
	insertCtx, err := collection.InsertMany(context.TODO(), newCourses)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Courses Deleted: %d, Courses Inserted: %d\n", deleteCtx.DeletedCount, len(insertCtx.InsertedIDs))
}

// AddOptionsToDB stores option information that was sourced from the page in mongodb
func (page *Page) AddOptionsToDB(db *mongo.Database, collectionName string, options []model.OptionData) {
	// Connect to collection for specific option type
	collection := db.Collection(collectionName)

	newOptions := []interface{}{}

	// Convert list of option data into type compatible with mongo insertions (list of empty interface)
	// Add extra metadata to data as well
	for _, option := range options {
		newOptions = append(newOptions, model.Option{
			Source: page.BuildSourceInfo(),
			Time:   page.BuildTimeInfo(),
			Data:   option,
		})
	}

	// Delete all, keep collection as up to date as possible
	deleteCtx, err := collection.DeleteMany(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println(err)
	}

	// Insert new options
	insertCtx, err := collection.InsertMany(context.TODO(), newOptions)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s Deleted: %d, %s Inserted: %d\n",
		collectionName,
		deleteCtx.DeletedCount,
		collectionName,
		len(insertCtx.InsertedIDs))
}
