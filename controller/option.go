package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"uwo-tt-api/model"
	_ "uwo-tt-api/model" // Placeholder; will be required soon

	"go.mongodb.org/mongo-driver/bson"
)

// ListSubjects godoc
// @Summary List all course subjects
// @Description Get list of subjects where each subject is a string pair representing form text and form value
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /subjects [get]
func (c *Controller) ListSubjects(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Subjects")

	collection := c.DB.Collection("subjects")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in subjects", len(options))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// ListSuffixes godoc
// @Summary List all course suffixes
// @Description Get list of suffixes where each suffix is a string pair representing form text and form value. The suffix represents course weight and session.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /suffixes [get]
func (c *Controller) ListSuffixes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Suffixes")

	collection := c.DB.Collection("suffixes")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in suffixes", len(options))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// ListDeliveryTypes godoc
// @Summary List all course delivery types
// @Description Get list of delivery types where each delivery type is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /delivery_types [get]
func (c *Controller) ListDeliveryTypes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Delivery Types")

	collection := c.DB.Collection("delivery_types")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in delivery types", len(options))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// ListComponents godoc
// @Summary List all course components types
// @Description Get list of course component types where each component type is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /components [get]
func (c *Controller) ListComponents(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Components")

	collection := c.DB.Collection("components")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in components", len(options))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// ListStartTimes godoc
// @Summary List all course start times
// @Description Get list of possible course start times where each time element is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /start_times [get]
func (c *Controller) ListStartTimes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Start Times")

	collection := c.DB.Collection("start_times")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in start times", len(options))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// ListEndTimes godoc
// @Summary List all course end times
// @Description Get list of possible course end times where each time element is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /end_times [get]
func (c *Controller) ListEndTimes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("End Times")

	collection := c.DB.Collection("end_times")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in end times", len(options))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// ListCampuses godoc
// @Summary List all possible campuses for courses
// @Description Get list of possible campuses where each campus is a string pair representing form text and form value.
// @Tags options
// @Produce json
// @Success 200 {array} model.Option
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /campuses [get]
func (c *Controller) ListCampuses(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Campuses")

	collection := c.DB.Collection("campuses")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	//Define an array in which you can store the decoded documents
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found %d documents in campuses", len(options))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(options)
}

// Course number and days are a given; simply provided usage description in documentation
