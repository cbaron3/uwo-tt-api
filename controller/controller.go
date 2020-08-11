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
	Inclusive bool `json:"inclusive" schema:"inclusive"`

	SectionNumber      []string `json:"section-number" schema:"section-number" example:"gte:001"`
	SectionComponent   []string `json:"section-component" schema:"section-component" example:"exact:TUT"`
	SectionClassNumber []string `json:"section-class-number" schema:"section-class-number" example:"lt:1000"`
	SectionLocation    []string `json:"section-location" schema:"section-location" example:"exact:NS-145"`
	SecionInstructor   []string `json:"section-instructor" schema:"section-instructor" example:"exact:Rahman"`
	SectionReqs        []string `json:"section-reqs" schema:"section-reqs"`
	SectionStatus      []string `json:"section-status" schema:"section-status" example:"exact:Full"`
	SectionCampus      []string `json:"section-campus" schema:"section-campus"	example:"exact:Main"`
	SectionDelivery    []string `json:"section-delivery" schema:"section-delivery" example:"exact:Distance Studies/Online"`
	SectionDay         []string `json:"section-time-day" schema:"section-time-day" example:"exact:M"`
	SectionStartTime   []string `json:"section-time-start-time" schema:"section-time-start-time" example:"except:8:30 AM"`
	SectionEndTime     []string `json:"section-time-end-time" schema:"section-time-end-time" example:"lte:7:00 PM"`

	ClassFaculty     []string `json:"course-faculty" schema:"course-faculty" example:"exact:PSYCH"`
	ClassNumber      []string `json:"course-number" schema:"course-number" example:"gte:3000"`
	ClassSuffix      []string `json:"course-suffix" schema:"course-suffix" example:"exact:F"`
	ClassName        []string `json:"course-name" schema:"course-name" example:"exact:INTRODUCTION TO PSYCHOLOGY"`
	ClassDescription []string `json:"course-description" schema:"course-description"`

	SortBy string `json:"sortby" schema:"sortby" example:"sortby=course-number"`
	Dec    bool   `json:"dec" schema:"dec" example:"true"`

	Offset int `json:"offset" schema:"offset" example:"10"`

	Limit int `json:"limit" schema:"limit" example:"5"`
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
		fmt.Println(err)
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
			return bson.M{}, fmt.Errorf("Invalid section course number  %d", op)
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

		sort := strings.Split(params.SortBy, "-")

		if len(sort) < 2 || len(sort) > 4 {
			return options.Find(), errors.New("Invalid sort criteria")
		}

		sortParam := ""

		// TODO: FIX TIME SORTING
		// Instead of model vs controller, create each as their own resource. better code organization

		if sort[0] == "course" {
			sortParam += "courseData."
		} else if sort[0] == "section" {
			sortParam += "sectionData."
		} else {
			return options.Find(), errors.New("Invalid sort criteria")
		}

		index := 1

		if sort[1] == "time" {
			sortParam += "times."
			index++
		}

		for i := index; i < len(sort); i++ {
			if i == index {
				sortParam += sort[i]
			} else {
				sortParam += strings.Title(strings.ToLower(sort[i]))
			}
		}

		fmt.Println(sortParam)

		// By default, sort ascending unless descending is specfied
		if params.Dec == true {
			result.SetSort(bson.D{{sortParam, -1}})
		} else {
			result.SetSort(bson.D{{sortParam, 1}})
		}
	}

	// Determine pagination parameters if they exist
	if params.Limit != 0 {
		result.SetLimit(int64(params.Limit))

		// Can only create a skip in records if the limit is known
		if params.Offset != 0 {
			result.SetSkip(int64(params.Offset - 1))
		}
	}

	return result, nil
}

// OptionQueryParams for decoding (gorilla) query params into a struct for handling
type OptionQueryParams struct {
	Inclusive bool     `json:"inclusive" schema:"inclusive"`
	Value     []string `json:"value" schema:"value" example:"exact:Main"`
	Text      []string `json:"text" schema:"text" example:"gte:ACTURSCI"`

	SortBy string `json:"sortby" schema:"sortby" example:"sortby=value"`
	Dec    bool   `json:"dec" schema:"dec" example:"true"`

	Offset int `json:"offset" schema:"offset" example:"10"`
	Limit  int `json:"limit" schema:"limit" example:"5"`
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
			result.SetSkip(int64(params.Offset - 1))
		}
	}

	return result, nil
}
