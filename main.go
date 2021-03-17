package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"

	cache "github.com//bekjeera/k8s/tree/master/myApp/cache"
	redis "github.com/bekjeera/k8s/tree/master/myApp/cache/redis"
	cfg "github.com//bekjeera/k8s/tree/master/myApp/config"
	handler "github.com/bekjeera/k8s/tree/master/myApp/handler"
)

var (
	config  cfg.Config
	storage cache.Storage
)

func main() {
	validate := validator.New()

	err := cfg.LoadConfig(validate, &config)
	if err != nil {
		log.Fatalf("Error while loading configuration:\n %v", err.Error())
	}

	storageStart()
	serverStart()
}

func storageStart() {
	var err error
	connectionInfo := fmt.Sprintf(
		"%s://%s:%d",
		config.Database.Driver,
		config.Database.Host,
		config.Database.Port,
	)
	if storage, err = redis.NewStorage(connectionInfo); err != nil {
		log.Fatalf("Redis error: %v", err)
	}
	fmt.Println("Redis started!")
}

func serverStart() {

	router := mux.NewRouter()

	wh := &handler.WeatherHandler{
		APIEnpoint:    config.API.Endpoint,
		APIKey:        config.API.Apikey,
		CacheDuration: config.Cache.Duration,
		Storage:       storage,
	}
	router.Handle("/weather", handler.ErrorHandler(wh.GetWeatherHandler)).Methods("GET")

	server := http.Server{
		Addr:        ":" + strconv.Itoa(config.Server.Port),
		Handler:     router,
		ReadTimeout: 10 * time.Second,
	}

	fmt.Println("Server started!")
	log.Fatal(server.ListenAndServe())
}
