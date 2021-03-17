# Get current weather data for several cities

REST API for getting weather forecast using 3rd party API. The request to my API can contain several cities, each city should be recognized as separate concurrent request to 3rd party. If during the last 60 seconds the city is repeated in next request, it shold be taken from cache which is implemented using Redis. The response is parsed from Json and returned as simple string.

# How to use it

Run in command line

```bash
$ mkdir myApp
$ cd myApp
$ git clone https://github.com/kravcs/weather_v2.git .
$ docker-compose up
```

After application built, run the sample call

```bash
$ curl -X GET "http://localhost:8000/weather?city=Vienna,Moscow"
```
