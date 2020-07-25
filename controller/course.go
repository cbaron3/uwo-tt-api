package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	_ "uwo-tt-api/httputil"
	"uwo-tt-api/model"
	_ "uwo-tt-api/model"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CourseParameters filter for course query
type CourseParameters struct {
	SubjectFilters []string `bson:"subjectFilters" json:"subjectFilters" example:"ACTURSCI,MSE,Psychology"`
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

// For integers db.MyCollection.find({ $expr: { $lte: [ { $toDouble: "$Price" }, 1000.0 ] } })
// https://stackoverflow.com/questions/47915936/mongo-cast-string-to-number-for-query
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ExtractCourseFilter(r *http.Request) (bson.M, error) {

	cmds := map[string]string{
		"exact":  "$eq",
		"except": "$ne",
		"gt":     "$gt",
		"gte":    "$gte",
		"lt":     "$lt",
		"lte":    "$lte",
	}

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

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.number": bson.M{val: num}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section number command %d", f[0]))
			}
		}

		for _, value := range filter.SectionComponent {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
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

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.classNumber": bson.M{val: num}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section class number  %d", f[0]))
			}
		}

		for _, value := range filter.SectionLocation {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.location": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section location command %s", f[0]))
			}
		}

		for _, value := range filter.SecionInstructor {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.instructor": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section instructor command %s", f[0]))
			}
		}

		for _, value := range filter.SectionReqs {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.requisites": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section requisites command %s", f[0]))
			}
		}

		for _, value := range filter.SectionStatus {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.status": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section status command %s", f[0]))
			}
		}

		for _, value := range filter.SectionCampus {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.campus": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section campus command %s", f[0]))
			}
		}

		for _, value := range filter.SectionDelivery {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.delivery": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section delivery command %s", f[0]))
			}
		}

		for _, value := range filter.SectionDay {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.times.days": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section days command %s", f[0]))
			}
		}

		for _, value := range filter.SectionStartTime {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.times.startTime": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section start time command %s", f[0]))
			}
		}

		for _, value := range filter.SectionEndTime {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"sectionData.times.endTime": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid section end time command %s", f[0]))
			}
		}

		for _, value := range filter.ClassFaculty {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
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

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.number": bson.M{val: num}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course number command %s", f[0]))
			}
		}

		for _, value := range filter.ClassSuffix {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.suffix": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course suffix command %s", f[0]))
			}
		}

		for _, value := range filter.ClassName {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"courseData.name": bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid course name command %s", f[0]))
			}
		}

		for _, value := range filter.ClassDescription {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
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

func (c *Controller) ListSections(w http.ResponseWriter, r *http.Request) {
	// Leaves sections as is
	w.Header().Set("Content-Type", "application/json")
	hitEndpoint("courses")

	collection := c.DB.Collection("courses")

	// Check if url can be parsed
	if err := r.ParseForm(); err != nil {
		fmt.Println("Form failed to parse")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Extract find filters
	findFilter, err := ExtractCourseFilter(r)
	if err != nil {
		fmt.Println("Filters failed to extract")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	findOptions, err := ExtractCourseParams(r)
	if err != nil {
		fmt.Println("Options failed to extract")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cur, err := collection.Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		fmt.Println("DB query failed; malformed filter or option")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//Define an array in which you can store the decoded documents
	var sections []model.Section

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Section
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Failed to decode")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		sections = append(sections, elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Failed to iterate through collection")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in %s", len(sections), "courses")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sections)

}

func (c *Controller) ListCourses(w http.ResponseWriter, r *http.Request) {
	// Leaves sections as is
	w.Header().Set("Content-Type", "application/json")
	hitEndpoint("courses")

	collection := c.DB.Collection("courses")

	// Check if url can be parsed
	if err := r.ParseForm(); err != nil {
		fmt.Println("Form failed to parse")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Extract find filters
	findFilter, err := ExtractCourseFilter(r)
	if err != nil {
		fmt.Println("Filters failed to extract")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	findOptions, err := ExtractCourseParams(r)
	if err != nil {
		fmt.Println("Options failed to extract")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cur, err := collection.Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		fmt.Println("DB query failed; malformed filter or option")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//Define an array in which you can store the decoded documents
	var sections []model.Section

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Section
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Failed to decode")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		sections = append(sections, elem)
	}

	// Create list of courses
	var courses []model.Course

	for _, section := range sections {
		if len(courses) == 0 {
			// If the first course, simply wrap into course and append
			var course model.Course
			course.Source = section.Source
			course.Time = section.Time
			course.CourseData = section.CourseData
			course.SectionData = append(course.SectionData, section.SectionData)

			courses = append(courses, course)
		} else {
			last := courses[len(courses)-1]

			if last.CourseData == section.CourseData {
				courses[len(courses)-1].SectionData = append(courses[len(courses)-1].SectionData, section.SectionData)
			} else {
				var course model.Course
				course.Source = section.Source
				course.Time = section.Time
				course.CourseData = section.CourseData
				course.SectionData = append(course.SectionData, section.SectionData)

				courses = append(courses, course)
			}

		}
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Failed to iterate through collection")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in %s", len(courses), "courses")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}

// ListCourses godoc
// @Summary List courses
// @Description List all courses and the sections in each course that match query filters
// @Tags courses
// @Accept  json
// @Produce  json
// @Param subjects body CourseParameters false "Subjects"
// @Success 200 {array} model.Course
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /courses [get]
//func (c *Controller) ListCourses(w http.ResponseWriter, r *http.Request) {
// hitEndpoint("Courses")

// fmt.Println("GET params were:", r.URL.Query())

// collection := c.DB.Collection("courses")
// // Create empty search filter
// filter := bson.A{}

// subject := r.URL.Query().Get("subject")
// if subject != "" {
// 	// TODO: Validate if subject is valid or not
// 	filter = append(filter, bson.M{"data.faculty": subject})
// }

// // TODO: Suffix is currently not easily supported
// // suffix := r.URL.Query().Get("suffix")
// // if suffix != "" {
// // 	// TODO: Validate if suffix is valid or not
// // 	filter = append(filter, bson.M{"data.TBD": suffix})
// // }

// // TODO: Suffix is currently not easily supported
// // number := r.URL.Query().Get("number")
// // if number != "" {
// // 	// TODO: Validate if number is valid or not
// // 	filter = append(filter, bson.M{"data.TBD": number})
// // }

// deliveryType := r.URL.Query().Get("delivery")
// deliveryTypeValid := true
// if deliveryType != "" {
// 	// TODO: Validate if delivery_type is valid or not
// 	// deliveryTypeValid = false

// 	filter = append(filter, bson.M{"data.sections.delivery": deliveryType})

// }

// campus := r.URL.Query().Get("campus")
// campusValid := true
// if campus != "" {
// 	// TODO: Validate if campus is valid or not
// 	// campusValid = false

// 	filter = append(filter, bson.M{"data.sections.campus": campus})
// }

// instructor := r.URL.Query().Get("instructor")
// instructorValid := true
// if instructor != "" {
// 	// TODO: Validate if campus is valid or not
// 	// campusValid = false

// 	filter = append(filter, bson.M{"data.sections.instructor": instructor})
// }

// // open := r.URL.Query().Get("open")

// // component := r.URL.Query().Get("component")
// // start_time := r.URL.Query().Get("start_time")
// // end_time := r.URL.Query().Get("end_time")

// days := r.URL.Query()["day"]
// daysValid := true
// if len(days) > 0 {
// 	daysFilter := bson.A{}

// 	for _, day := range days {
// 		daysFilter = append(daysFilter, bson.M{"data.sections.days": day})
// 	}

// 	filter = append(filter, bson.M{"$or": daysFilter})
// }

// // campus := r.URL.Query().Get("campus")

// cur, err := collection.Find(context.TODO(), bson.M{"$and": filter})
// if err != nil {
// 	log.Fatal(err)
// }

// //Define an array in which you can store the decoded documents
// var courses []model.Course

// for cur.Next(context.TODO()) {
// 	//Create a value into which the single document can be decoded
// 	var elem model.Course
// 	err := cur.Decode(&elem)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// IMPROVE: Current database modelling technique results in all sections being returned when one section matches a query filter
// 	// Could redesign database but easier solution is to just only add sections that match the filter
// 	var validSections []model.SectionData
// 	for _, section := range elem.Data.Sections {
// 		canAdd := true

// 		if (!campusValid || section.Campus != campus) && campus != "" {
// 			canAdd = canAdd && false
// 			fmt.Println("campus clears")
// 		} else if (!deliveryTypeValid || section.Delivery != deliveryType) && deliveryType != "" {
// 			canAdd = canAdd && false
// 			fmt.Println("delivery clears")
// 		} else if (!instructorValid || section.Instructor != instructor) && instructor != "" {
// 			canAdd = canAdd && false
// 			fmt.Println("instructor clears")
// 		} else if (!daysValid || !stringInSlice(section.Days, days)) && len(days) > 0 {
// 			canAdd = canAdd && false
// 			fmt.Println("days clears")
// 		}

// 		if canAdd {
// 			validSections = append(validSections, section)
// 		}
// 	}

// 	elem.Data.Sections = validSections
// 	courses = append(courses, elem)
// }

// if err := cur.Err(); err != nil {
// 	log.Fatal(err)
// }

// //Close the cursor once finished
// cur.Close(context.TODO())

// fmt.Printf("Found %d documents in courses with filters", len(courses))

// w.Header().Set("Content-Type", "application/json")
// w.WriteHeader(http.StatusOK)
// json.NewEncoder(w).Encode(courses)
//}

// CourseValue to CourseText
// CourseText to CourseValue
// OptionText to OptionValue
// OptionValue to OptionText
