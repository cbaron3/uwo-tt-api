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
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/schema"

	"strings"
	"errors"
	
)

type OptQueryParams struct {
	// Filtering
	Inclusive bool `schema:"inclusive"`  
	Value []string `schema:"value"`
	Text []string `schema:"text"`

	// Sorting
	SortBy string `schema:"sortby"`
	Dec    bool   `schema:"dec"`

	// Pagination 
    Offset int  `schema:"offset"`
    Limit  int  `schema:"limit"`

}

func ExtractOptFilter(r *http.Request) (bson.M, error) {

	cmds := map[string]string{
		"exact": "$eq",
		"gt":   "$gt",
		"gte": "$gte",
		"lt": "$lt",
		"lte": "$lte",
	}


	filter := new(OptQueryParams)
	
	result := bson.M{}
	arrfilter := bson.A{}

    if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		fmt.Println("Opt decoding failed")
	} else {
		

		for _, value := range filter.Value {
			f := strings.Split(value, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"data.value" : bson.M{val: f[1]}})
			} else {
				return bson.M{}, errors.New(fmt.Sprintf("Invalid filter command %s", f[0]))
			}
		}

		for _, text := range filter.Text {
			f := strings.Split(text, ":")

			if val, ok := cmds[f[0]]; ok {
				arrfilter = append(arrfilter, bson.M{"data.text" : bson.M{val: f[1]}})
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
	filter := new(OptQueryParams)
	
    if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		return opts, errors.New("Option decoding failed")
	} else {
		fmt.Println(filter)

		// Determine sort parameters if they exist
		if filter.SortBy != "" {
			// By default, sort ascending unless descending is specfied
			if filter.Dec == true {
				opts.SetSort(bson.D{{"data."+filter.SortBy, -1}})
			} else {
				opts.SetSort(bson.D{{"data."+filter.SortBy, 1}})
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

func (c *Controller) optionsEndpoint(collectionName string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hitEndpoint(collectionName)

	collection := c.DB.Collection(collectionName)

	// Check if url can be parsed
	if err := r.ParseForm(); err != nil {
		fmt.Println("Form failed to parse")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Extract find filters
	findFilter, err := ExtractOptFilter(r)
	if err != nil {
		fmt.Println("Filters failed to extract")
		// Handle error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	findOptions, err := ExtractOptParams(r)
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
	var options []model.Option

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.Option
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Failed to decode")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		options = append(options, elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Failed to iterate through collection")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
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
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /campuses [get]
func (c *Controller) ListCampuses(w http.ResponseWriter, r *http.Request) {
	c.optionsEndpoint("campuses", w, r)
}

// Course number and days are a given; simply provided usage description in documentation
