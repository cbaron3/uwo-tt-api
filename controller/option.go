package controller

import (
	"net/http"
	_ "uwo-tt-api/model" // Placeholder; will be required soon
)

// ListSubjects Handler for listing subjects
func (c *Controller) ListSubjects(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Subjects")
}

// ListSuffixes Handler for listing suffixes
func (c *Controller) ListSuffixes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Suffixes")
}

// ListDeliveryTypes Handler for listing delivery
func (c *Controller) ListDeliveryTypes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Delivery Types")
}

// ListComponents Handler for listing components
func (c *Controller) ListComponents(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Components")
}

// ListStartTimes Handler for listing start times
func (c *Controller) ListStartTimes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Start Times")
}

// ListEndTimes Handler for listing end times
func (c *Controller) ListEndTimes(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("End Times")
}

// ListCampuses Handler for listing campuses
func (c *Controller) ListCampuses(w http.ResponseWriter, r *http.Request) {
	hitEndpoint("Campuses")
}

// Course number and days are a given; simply provided usage description in documentation
