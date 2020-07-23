package controller

import (
	"net/http"
	_ "uwo-tt-api/httputil"
	_ "uwo-tt-api/model"
)

// CourseParameters filter for course query
type CourseParameters struct {
	SubjectFilters []string `bson:"subjectFilters" json:"subjectFilters" example:"ACTURSCI,MSE,Psychology"`
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (c* Controller) ListSections(w http.ResponseWriter, r *http.Request) {
	// Leaves sections as is
}


func (c *Controller) ListCourses(w http.ResponseWriter, r *http.Request) {
	// Combines sections
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
