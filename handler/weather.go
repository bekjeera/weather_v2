package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/kravcs/weather_v2/cache"
	"github.com/kravcs/weather_v2/model"
)

type Data struct {
	Result model.CityWeather
	Err    error
}

type WeatherHandler struct {
	APIEnpoint    string
	APIKey        string
	CacheDuration int
	Storage       cache.Storage
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// GetWeatherHandler provides weather forecast for one or several cities
func (wh WeatherHandler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	log.Println("handler started")
	defer timeTrack(time.Now(), "weather handler")

	cityParam := r.FormValue("city")
	if cityParam == "" {
		return NewStatusError(http.StatusBadRequest, "Query parameter \"city\" is missing")
	}
	citiesToForecast := strings.Split(cityParam, ",")

	cityChannel := make(chan Data, len(citiesToForecast))

	for _, city := range citiesToForecast {
		log.Printf("go routine for %s \n", city)
		go wh.getForecast(ctx, city, cityChannel)
	}

	var cities []model.CityWeather

	for {
		select {
		case data := <-cityChannel:
			if data.Err != nil {
				return NewStatusError(http.StatusBadRequest, data.Err.Error())
			}

			cities = append(cities, data.Result)

			log.Println("got city from goroutine \n")

			if len(cities) == len(citiesToForecast) {
				log.Println("Got all the cities \n")
				log.Println(http.StatusOK)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-type", "application/json")
				json.NewEncoder(w).Encode(cities)

				return nil
			}

		case <-ctx.Done():
			log.Println(ctx.Err().Error())
			return NewStatusError(http.StatusBadRequest, ctx.Err().Error())
		}
	}
}

func (wh WeatherHandler) getForecast(ctx context.Context, city string, cityChannel chan<- Data) {

	select {
	case <-ctx.Done():
		log.Println("(fail): Something happen")
		return
	default:
		var cityObj model.CityWeather

		log.Printf("try to get  %s from cache \n", city)
		// try to take data from chache
		cityForecast := wh.Storage.Get(city)
		if len(cityForecast) != 0 {
			json.Unmarshal(cityForecast, &cityObj)
			cityChannel <- Data{Result: cityObj}

			log.Printf("I've got the city from cahce %s \n", city)

			return
		}

		log.Printf("try to get  %s from api \n", city)

		// try to take data from 3rd api
		err := wh.search(ctx, city, &cityObj)
		if err != nil {
			cityChannel <- Data{Err: err}
			return
		}

		// try to save content in cache
		content, _ := json.Marshal(cityObj)
		wh.Storage.Set(city, content, time.Duration(60*time.Second))
		cityChannel <- Data{Result: cityObj}
		return
	}
}

func (wh WeatherHandler) search(ctx context.Context, city string, this *model.CityWeather) error {

	// waiting for cancelation
	select {
	case <-ctx.Done():
		return nil
	default:
		log.Printf("try to get  %s from api search \n", city)

		url := wh.APIEnpoint + "?APPID=" + wh.APIKey + "&q=" + city

		res, err := http.Get(url)
		if err != nil {
			return errors.Wrap(err, "Error while GET url response from 3rd api")
		}

		body, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return errors.New("Error StatusCode from 3rd api: " + string(body))
		}

		err = json.Unmarshal(body, this)
		if err != nil {
			return errors.Wrap(err, "Error response body from 3rd api")
		}

		log.Printf("city %s from api call is ok \n", city)
		return nil
	}
}
