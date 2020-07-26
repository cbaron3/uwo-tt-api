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

// Controller example
type Controller struct {
	DB *mongo.Database
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

func hitEndpoint(name string) {
	fmt.Printf("*** ENDPOINT RESOURCE HIT --> %s\n", name)
}

var FilterToDB = map[string]string{
	"exact":  "$eq",
	"except": "$ne",
	"gt":     "$gt",
	"gte":    "$gte",
	"lt":     "$lt",
	"lte":    "$lte",
}

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

func ExtractCourseFilter(r *http.Request) (bson.M, error) {

	filter := new(CourseQueryParams)

	result := bson.M{}
	arrfilter := bson.A{}

	if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		fmt.Println("Opt decoding failed")
	} else {

		fmt.Println(filter)

		for _, value := range filter.SectionNumber {
			f := strings.Split(value, ":")

			num, err := strconv.Atoi(f[1])
			if err != nil {
				// handle error
				fmt.Println(err)
				num = 0
			}

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.number": bson.M{val: num}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section number command %d", f[0]))
			}
		}

		for _, value := range filter.SectionComponent {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.component": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section component command %s", f[0]))
			}
		}

		for _, value := range filter.SectionClassNumber {
			f := strings.Split(value, ":")

			num, err := strconv.Atoi(f[1])
			if err != nil {
				// handle error
				fmt.Println(err)
				num = 0
			}

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.classNumber": bson.M{val: num}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section class number  %d", f[0]))
			}
		}

		for _, value := range filter.SectionLocation {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.location": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section location command %s", f[0]))
			}
		}

		for _, value := range filter.SecionInstructor {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.instructor": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section instructor command %s", f[0]))
			}
		}

		for _, value := range filter.SectionReqs {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.requisites": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section requisites command %s", f[0]))
			}
		}

		for _, value := range filter.SectionStatus {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.status": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section status command %s", f[0]))
			}
		}

		for _, value := range filter.SectionCampus {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.campus": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section campus command %s", f[0]))
			}
		}

		for _, value := range filter.SectionDelivery {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.delivery": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section delivery command %s", f[0]))
			}
		}

		for _, value := range filter.SectionDay {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.times.days": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section days command %s", f[0]))
			}
		}

		for _, value := range filter.SectionStartTime {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.times.startTime": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section start time command %s", f[0]))
			}
		}

		for _, value := range filter.SectionEndTime {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.times.endTime": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section end time command %s", f[0]))
			}
		}

		for _, value := range filter.ClassFaculty {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.faculty": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course faculty command %s", f[0]))
			}
		}

		for _, value := range filter.ClassNumber {
			f := strings.Split(value, ":")

			num, err := strconv.Atoi(f[1])
			if err != nil {
				// handle error
				fmt.Println(err)
				num = 0
			}

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.number": bson.M{val: num}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course number command %s", f[0]))
			}
		}

		for _, value := range filter.ClassSuffix {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.suffix": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course suffix command %s", f[0]))
			}
		}

		for _, value := range filter.ClassName {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.name": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course name command %s", f[0]))
			}
		}

		for _, value := range filter.ClassDescription {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.description": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course description command %s", f[0]))
			}
		}

		if len(arrfilter) > 0 {
			if filter.Inclusive == true {
				result = bson.M{"$or": arrfilter}
			} else {
				result = bson.M{"$and": arrfilter}
			}
		}

	}

	return result, nil
}

// No need for filtering with options
func ExtractCourseParams(r *http.Request) (*options.FindOptions, error) {

	opts := options.Find()
	filter := new(CourseQueryParams)

	if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		return opts, errors.New("Option decoding failed")
	} else {
		fmt.Println(filter)

		// Determine sort parameters if they exist
		if filter.SortBy != "" {
			// By default, sort ascending unless descending is specfied
			if filter.Dec == true {
				opts.SetSort(bson.D{{"data." + filter.SortBy, -1}})
			} else {
				opts.SetSort(bson.D{{"data." + filter.SortBy, 1}})
			}
		}

		// Determine pagination parameters if they exist
		if filter.Limit != 0 {
			opts.SetLimit(int64(filter.Limit))

			// Can only create a skip in records if the limit is known
			if filter.Offset != 0 {
				opts.SetSkip(int64(filter.Limit * (filter.Offset - 1)))
			}
		}
	}

	return opts, nil
}

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

func ExtractOptFilter(r *http.Request) (bson.M, error) {

	filter := new(OptionQueryParams)

	result := bson.M{}
	arrfilter := bson.A{}

	if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		fmt.Println("Opt decoding failed")
	} else {

		for _, value := range filter.Value {
			f := strings.Split(value, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"data.value": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid filter command %s", f[0]))
			}
		}

		for _, text := range filter.Text {
			f := strings.Split(text, ":")

			if val, ok := FilterToDB[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"data.text": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid filter command %s", f[0]))
			}
		}

		if len(arrfilter) > 0 {
			if filter.Inclusive == true {
				result = bson.M{"$or": arrfilter}
			} else {
				result = bson.M{"$and": arrfilter}
			}
		}

	}

	return result, nil
}

// No need for filtering with options
func ExtractOptParams(r *http.Request) (*options.FindOptions, error) {

	opts := options.Find()
	filter := new(OptionQueryParams)

	if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		return opts, errors.New("Option decoding failed")
	} else {
		fmt.Println(filter)

		// Determine sort parameters if they exist
		if filter.SortBy != "" {
			// By default, sort ascending unless descending is specfied
			if filter.Dec == true {
				opts.SetSort(bson.D{{"data." + filter.SortBy, -1}})
			} else {
				opts.SetSort(bson.D{{"data." + filter.SortBy, 1}})
			}
		}

		// Determine pagination parameters if they exist
		if filter.Limit != 0 {
			opts.SetLimit(int64(filter.Limit))

			// Can only create a skip in records if the limit is known
			if filter.Offset != 0 {
				opts.SetSkip(int64(filter.Limit * (filter.Offset - 1)))
			}
		}
	}

	return opts, nil
}
