package controller

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
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
