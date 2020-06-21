package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/mux"
	
	"github.com/spf13/viper"

	"github.com/go-co-op/gocron"

	"net/http"

	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/gorilla/docs"
)

func task() {
    fmt.Println("I am running task.")
}

func home(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func subjectHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Subjects")
	fmt.Println("Endpoint Hit: Subjects")
}

func suffixHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Suffixes")
	fmt.Println("Endpoint Hit: Suffixes")
}

func deliveryHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Delivery Types")
	fmt.Println("Endpoint Hit: Delivery Types")
}

func componentHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Components")
	fmt.Println("Endpoint Hit: Components")
}

func startTimeHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Start Times")
	fmt.Println("Endpoint Hit: Start Times")
}

func endTimeHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "End Times")
	fmt.Println("Endpoint Hit: End Times")
}

func dayHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Days")
	fmt.Println("Endpoint Hit: Days")
}

func campusHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Campuses")
	fmt.Println("Endpoint Hit: Campuses")
}

func courseHandler(w http.ResponseWriter, r *http.Request){
	// Multiple subject filters
	subjects := r.URL.Query()["subject"]
	if len(subjects) > 0 {
		fmt.Println(subjects)
	} else {
		fmt.Println("No subject data in URL")
	}

	// Multiple suffix filters
	suffixes := r.URL.Query()["suffix"]
	if len(suffixes) > 0 {
		fmt.Println(suffixes)
	} else {
		fmt.Println("No suffix data in URL")
	}

	// Multiple number filters
	// TODO: Add wildcarding idea. So every course with number 3xxx or 4xxx
	numbers := r.URL.Query()["number"]
	if len(numbers) > 0 {
		fmt.Println(numbers)
	} else {
		fmt.Println("No number data in URL")
	}

	// Multiple delivery filters
	deliveryTypes := r.URL.Query()["delivery"]
	if len(deliveryTypes) > 0 {
		fmt.Println(deliveryTypes)
	} else {
		fmt.Println("No delivery data in URL")
	}


	// No other options besides True OR False
	// Returns empty string if no value
	open := r.URL.Query().Get("open")
	if open != "" {
		fmt.Println(open)
	} else {
		fmt.Println("No course registration open data in URL")
	}

	// Multiple components possible
	component := r.URL.Query()["component"]
	if len(component) > 0 {
		fmt.Println(component)
	} else {
		fmt.Println("No component data in URL")
	}

	// Upper time bound
	startTime := r.URL.Query().Get("start_time")
	if startTime != "" {
		fmt.Println(startTime)
	} else {
		fmt.Println("No start time data in URL")
	}

	// Lower time bound
	endTime := r.URL.Query().Get("end_time")
	if endTime != "" {
		fmt.Println(endTime)
	} else {
		fmt.Println("No end time data in URL")
	}

	// Multiple days possible
	days := r.URL.Query()["day"]
	if len(days) > 0 {
		fmt.Println(days)
	} else {
		fmt.Println("No day data in URL")
	}

	// Multiple campuses possible
	campuses := r.URL.Query()["campus"]
	if len(campuses) > 0 {
		fmt.Println(campuses)
	} else {
		fmt.Println("No campus data in URL")
	}

	fmt.Fprintf(w, "Courses")
	fmt.Println("Endpoint Hit: Courses")
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
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
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", home)
	
	myRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	
	myRouter.HandleFunc("/subjects", subjectHandler)
	myRouter.HandleFunc("/suffixes", suffixHandler)
	myRouter.HandleFunc("/delivery_types", deliveryHandler)
	myRouter.HandleFunc("/components", componentHandler)
	myRouter.HandleFunc("/start_times", startTimeHandler)
	myRouter.HandleFunc("/end_times", endTimeHandler)
	myRouter.HandleFunc("/days", dayHandler)
	myRouter.HandleFunc("/campuses", campusHandler)
	myRouter.HandleFunc("/courses", courseHandler)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}