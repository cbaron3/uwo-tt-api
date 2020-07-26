package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"uwo-tt-api/model"
	_ "uwo-tt-api/model" // Placeholder; will be required soon
)

func (c *Controller) optionsEndpoint(collectionName string, w http.ResponseWriter, r *http.Request) {
	HitEndpoint("courses")

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Connect to collection
	collection := c.DB.Collection(collectionName)

	// Check if url can be parsed
	if err := r.ParseForm(); err != nil {
		w = NewError(w, http.StatusBadRequest, err, "Failed to parse option query parameters")
		return
	}

	// Extract find filters
	findFilter, err := ExtractOptFilter(r)
	if err != nil {
		w = NewError(w, http.StatusBadRequest, err, "Failed to extract option filters")
		return
	}

	findOptions, err := ExtractOptParams(r)
	if err != nil {
		w = NewError(w, http.StatusBadRequest, err, "Failed to extract option options")
		return
	}

	cur, err := collection.Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		w = NewError(w, http.StatusBadRequest, err, "DB query failed; malformed filter or option")
		return
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			w = NewError(w, http.StatusBadRequest, err, "Failed to decode db result")
			return
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		w = NewError(w, http.StatusBadRequest, err, "Failed to iterate over db results")
		return
	}

	//Close the cursor once finished
	fmt.Printf("Found %d documents in %s", len(options), collectionName)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// ListSubjects godoc
// @Summary List all course subjects
// @Description Get list of subjects where each subject is a string pair representing form text and form value
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Router /subjects [get]
func (c *Controller) ListSubjects(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("subjects", w, r)
}

// ListSuffixes godoc
// @Summary List all course suffixes
// @Description Get list of suffixes where each suffix is a string pair representing form text and form value. The suffix represents course weight and session.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Router /suffixes [get]
func (c *Controller) ListSuffixes(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("suffixes", w, r)
}

// ListDeliveryTypes godoc
// @Summary List all course delivery types
// @Description Get list of delivery types where each delivery type is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Router /delivery_types [get]
func (c *Controller) ListDeliveryTypes(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("delivery_types", w, r)
}

// ListComponents godoc
// @Summary List all course components types
// @Description Get list of course component types where each component type is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Router /components [get]
func (c *Controller) ListComponents(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("components", w, r)
}

// ListStartTimes godoc
// @Summary List all course start times
// @Description Get list of possible course start times where each time element is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Router /start_times [get]
func (c *Controller) ListStartTimes(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("start_times", w, r)
}

// ListEndTimes godoc
// @Summary List all course end times
// @Description Get list of possible course end times where each time element is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Router /end_times [get]
func (c *Controller) ListEndTimes(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("end_times", w, r)
}

// ListCampuses godoc
// @Summary List all possible campuses for courses
// @Description Get list of possible campuses where each campus is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Router /campuses [get]
func (c *Controller) ListCampuses(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("campuses", w, r)
}
