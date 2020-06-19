package main

import (
	"fmt"
	"log"
	"time"
	"github.com/spf13/viper"

	"github.com/go-co-op/gocron"

	"net/http"
)

func task() {
    fmt.Println("I am running task.")
}

func home(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}


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

	// defines a new scheduler that schedules and runs jobs
    s1 := gocron.NewScheduler(time.UTC)

    s1.Every(5).Seconds().StartImmediately().Do(task)

    // scheduler starts running jobs and current thread continues to execute
	s1.StartAsync()
	
	// s1.StartBlocking()

	http.HandleFunc("/", home)
    log.Fatal(http.ListenAndServe(":8080", nil))
}