package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Controller struct which acts as base for all endpoint methods
type Controller struct {
	DB *mongo.Database
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// HitEndpoint simple helper to log when an endpoint was hit
func HitEndpoint(name string) {
	fmt.Printf("*** ENDPOINT RESOURCE HIT --> %s\n", name)
}

// FilterToDBOp lookup table for query parameter filter commands to mongo command
var FilterToDBOp = map[string]string{
	"exact":  "$eq",
	"except": "$ne",
	"gt":     "$gt",
	"gte":    "$gte",
	"lt":     "$lt",
	"lte":    "$lte",
}

// CourseQueryParams for decoding (gorilla) query params into a struct for handling
type CourseQueryParams struct {
	// Filtering
	Inclusive          bool     `schema:"inclusive"`
	SectionNumber      []string `schema:"section-number"`
	SectionComponent   []string `schema:"section-component"`
	SectionClassNumber []string `schema:"section-class-number"`
	SectionLocation    []string `schema:"section-location"`
	SecionInstructor   []string `schema:"section-instructor"`
	SectionReqs        []string `schema:"section-reqs"`
	SectionStatus      []string `schema:"section-status"`
	SectionCampus      []string `schema:"section-campus"`
	SectionDelivery    []string `schema:"section-delivery"`
	SectionDay         []string `schema:"section-day"`
	SectionStartTime   []string `schema:"section-start-time"`
	SectionEndTime     []string `schema:"section-end-time"`
	ClassFaculty       []string `schema:"class-faculty"`
	ClassNumber        []string `schema:"class-number"`
	ClassSuffix        []string `schema:"class-suffix"`
	ClassName          []string `schema:"class-name"`
	ClassDescription   []string `schema:"class-desc"`

	// Sorting
	SortBy string `schema:"sortby"`
	Dec    bool   `schema:"dec"`

	// Pagination
	Offset int `schema:"offset"`
	Limit  int `schema:"limit"`
}

// ExtractCourseFilter extracts course filters from request
func ExtractCourseFilter(r *http.Request) (bson.M, error) {

	if r == nil {
		return bson.M{}, errors.New("Request object is nil")
	}

	// Create struct to decode params into
	params := new(CourseQueryParams)

	if err := schema.NewDecoder().Decode(params, r.Form); err != nil {
		fmt.Println("ExtractCourseFilter failed to decode request form into struct")
		return bson.M{}, errors.New("Course query filters failed to decode")
	}

	// Capture array of filters
	filters := bson.A{}

	// TODO: Let us ignore this monstrosity of code
	for _, value := range params.SectionNumber {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		num, err := strconv.Atoi(opValue)
		if err != nil {
			return bson.M{}, errors.New("Section number value failed to parse to integer")
		}

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.number": bson.M{val: num}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section number command %d", op)
		}
	}

	for _, value := range params.SectionComponent {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.component": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section component command %s", op)
		}
	}

	for _, value := range params.SectionClassNumber {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		num, err := strconv.Atoi(opValue)
		if err != nil {
			// handle error
			fmt.Println(err)
			num = 0
		}

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.classNumber": bson.M{val: num}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section class number  %d", op)
		}
	}

	for _, value := range params.SectionLocation {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.location": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section location command %s", op)
		}
	}

	for _, value := range params.SecionInstructor {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.instructor": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section instructor command %s", op)
		}
	}

	for _, value := range params.SectionReqs {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.requisites": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section requisites command %s", op)
		}
	}

	for _, value := range params.SectionStatus {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.status": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section status command %s", op)
		}
	}

	for _, value := range params.SectionCampus {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.campus": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section campus command %s", op)
		}
	}

	for _, value := range params.SectionDelivery {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.delivery": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section delivery command %s", op)
		}
	}

	for _, value := range params.SectionDay {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.times.days": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section days command %s", op)
		}
	}

	for _, value := range params.SectionStartTime {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.times.startTime": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section start time command %s", op)
		}
	}

	for _, value := range params.SectionEndTime {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"sectionData.times.endTime": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid section end time command %s", op)
		}
	}

	for _, value := range params.ClassFaculty {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"courseData.faculty": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid course faculty command %s", op)
		}
	}

	for _, value := range params.ClassNumber {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		num, err := strconv.Atoi(opValue)
		if err != nil {
			// handle error
			fmt.Println(err)
			num = 0
		}

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"courseData.number": bson.M{val: num}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid course number command %s", op)
		}
	}

	for _, value := range params.ClassSuffix {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"courseData.suffix": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid course suffix command %s", op)
		}
	}

	for _, value := range params.ClassName {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"courseData.name": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid course name command %s", op)
		}
	}

	for _, value := range params.ClassDescription {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"courseData.description": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid course description command %s", op)
		}
	}

	if len(filters) > 0 {
		if params.Inclusive == true {
			return bson.M{"$or": filters}, nil
		} else {
			return bson.M{"$and": filters}, nil
		}
	}

	return bson.M{}, nil
}

// ExtractCourseParams extract extra params from request besides filters into a set of find options
func ExtractCourseParams(r *http.Request) (*options.FindOptions, error) {

	if r == nil {
		return options.Find(), errors.New("Request object is nil")
	}

	// Create struct to decode params into
	params := new(CourseQueryParams)

	if err := schema.NewDecoder().Decode(params, r.Form); err != nil {
		return options.Find(), errors.New("Course query options failed to decode")
	}

	// Capture find options
	result := options.Find()

	// Determine sort parameters if they exist
	if params.SortBy != "" {
		// By default, sort ascending unless descending is specfied
		if params.Dec == true {
			result.SetSort(bson.D{{"data." + params.SortBy, -1}})
		} else {
			result.SetSort(bson.D{{"data." + params.SortBy, 1}})
		}
	}

	// Determine pagination parameters if they exist
	if params.Limit != 0 {
		result.SetLimit(int64(params.Limit))

		// Can only create a skip in records if the limit is known
		if params.Offset != 0 {
			result.SetSkip(int64(params.Limit * (params.Offset - 1)))
		}
	}

	return result, nil
}

// OptionQueryParams for decoding (gorilla) query params into a struct for handling
type OptionQueryParams struct {
	// Filtering
	Inclusive bool     `schema:"inclusive"`
	Value     []string `schema:"value"`
	Text      []string `schema:"text"`

	// Sorting
	SortBy string `schema:"sortby"`
	Dec    bool   `schema:"dec"`

	// Pagination
	Offset int `schema:"offset"`
	Limit  int `schema:"limit"`
}

// ExtractOptFilter extracts option filters from request
func ExtractOptFilter(r *http.Request) (bson.M, error) {

	if r == nil {
		return bson.M{}, errors.New("Request object is nil")
	}

	// Create struct to decode params into
	params := new(OptionQueryParams)

	if err := schema.NewDecoder().Decode(params, r.Form); err != nil {
		fmt.Println("ExtractOptFilter failed to decode request form into struct")
		return bson.M{}, errors.New("Option query filters failed to decode")
	}

	// Capture array of filters
	filters := bson.A{}

	for _, value := range params.Value {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"data.value": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid filter command %s", op)
		}
	}

	for _, value := range params.Text {
		f := strings.Split(value, ":")

		op := f[0]
		opValue := f[1]

		if val, ok := FilterToDBOp[op]; ok {
			filters = append(filters, bson.M{"data.text": bson.M{val: opValue}})
		} else {
			return bson.M{}, fmt.Errorf("Invalid filter command %s", op)
		}
	}

	if len(filters) > 0 {
		if params.Inclusive == true {
			return bson.M{"$or": filters}, nil
		} else {
			return bson.M{"$and": filters}, nil
		}
	}

	return bson.M{}, nil
}

// ExtractOptParams extract extra params from request besides filters into a set of find options
func ExtractOptParams(r *http.Request) (*options.FindOptions, error) {

	if r == nil {
		return options.Find(), errors.New("Request object is nil")
	}

	// Create struct to decode params into
	params := new(OptionQueryParams)

	if err := schema.NewDecoder().Decode(params, r.Form); err != nil {
		return options.Find(), errors.New("Course query options failed to decode")
	}

	// Capture find options
	result := options.Find()

	// Determine sort parameters if they exist
	if params.SortBy != "" {
		// By default, sort ascending unless descending is specfied
		if params.Dec == true {
			result.SetSort(bson.D{{"data." + params.SortBy, -1}})
		} else {
			result.SetSort(bson.D{{"data." + params.SortBy, 1}})
		}
	}

	// Determine pagination parameters if they exist
	if params.Limit != 0 {
		result.SetLimit(int64(params.Limit))

		// Can only create a skip in records if the limit is known
		if params.Offset != 0 {
			result.SetSkip(int64(params.Limit * (params.Offset - 1)))
		}
	}

	return result, nil
}
