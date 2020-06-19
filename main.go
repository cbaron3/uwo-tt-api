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
		log.Printf("Error reading config file %s. Using environment variables instead.", err)
		viper.AutomaticEnv()
	}
	
	value, ok := viper.Get("TEST_KEY").(string)

	if !ok {
		log.Fatalf("Key does not exist or has invalid type")
	}

	fmt.Println("KEY: ", value)
	fmt.Println("Docker test")
}