package main

import (
	"fmt"
	"log"
	"github.com/spf13/viper"
)

func main()  {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error reading config file %s", err)
	}

	value, ok := viper.Get("TEST_KEY").(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	fmt.Println("KEY: ", value)
}