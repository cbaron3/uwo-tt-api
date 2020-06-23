package controller

import "fmt"

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

func hitEndpoint(name string) {
	fmt.Printf("*** ENDPOINT RESOURCE HIT --> %s\n", name)
}