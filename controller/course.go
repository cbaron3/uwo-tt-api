package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"uwo-tt-api/model"
)

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
