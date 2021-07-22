package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/csv-loader/services"
	"github.com/csv-loader/services/csvbuilder"
	"github.com/csv-loader/services/csvextractor"
	echo "github.com/labstack/echo/v4"

	"github.com/csv-loader/services/loaders/redisloader"

	redis "github.com/go-redis/redis/v8"

	csvHandler "github.com/csv-loader/http"
	csvGetter "github.com/csv-loader/services/csvgetter"
	store "github.com/csv-loader/stores/csvstore"
)

func main() {
	file, err := os.Open("./csv-data/data.csv")
	if err != nil {
		log.Fatalf("Error while opening the csv file, error: %v", err)
	}

	csvReader := csv.NewReader(file)
	csvExtractor := csvextractor.New(csvReader)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASS"),
	})

	w, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		w = 8
	}

	chunkSize, err := strconv.Atoi(os.Getenv("CHUNK_SIZE"))
	if err != nil {
		chunkSize = 8000
	}

	builder := csvbuilder.New(csvExtractor, map[string]services.CSVLoader{
		"redis": redisloader.New(redisClient),
	}, w, chunkSize)

	// Load the CSV file in datastore
	t1 := time.Now()
	builder.BuildCSV()
	log.Printf("Time taken to load the entries: %v", time.Since(t1))

	redisStore := store.New(redisClient)
	csvGetterSvc := csvGetter.New(redisStore)
	handler := csvHandler.New(csvGetterSvc)

	e := echo.New()

	e.GET("csv/:key", handler.Read)

	e.Logger.Fatal(e.Start(":8000"))
}
