package redis

import (
	"context"
	"encoding/json"
	"log"
	"mushturd/pkg/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func FetchMushersFromCache(key string) ([]models.Musher, error) {
	data, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist in the cache
		log.Printf("Key '%s' does not exist in the cache", key)
		return nil, nil
	} else if err != nil {
		// Error retrieving data from cache
		log.Printf("Error retrieving data from cache for key '%s': %v", key, err)
		return nil, err
	}
	var mushers []models.Musher
	err = json.Unmarshal([]byte(data), &mushers)
	if err != nil {
		// Error deserializing data
		log.Printf("Error unmarshalling data '%v'", err)
		return nil, err
	}
	log.Printf("Data retrieved from cache for key '%s'", key)
	return mushers, nil
}

func GetMushersFromCacheOrAPI(cacheKey string, fetchDataFunc func() []models.Musher) ([]models.Musher, error) {
	ctx := context.Background()

	data, err := redisClient.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		// Key does not exist in cache
		log.Printf("Key '%s' does not exist in the cache", cacheKey)

		// Fetch data from API
		mushers := fetchDataFunc()
		// should add error handling to scraper

		// store fetched data in cache
		serialized, err := json.Marshal(mushers)
		if err != nil {
			log.Printf("Error marshalling data '%v'", err)
			return nil, err
		}

		err = redisClient.Set(ctx, cacheKey, serialized, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Error storing data in cache for key '%s': %v", cacheKey, err)
			return nil, err
		}
		// return queried data
		return mushers, nil

	} else if err != nil {
		log.Printf("Error retrieving data from cache for key '%s': %v", cacheKey, err)
		return nil, err
	}

	// Data found in cache
	var mushers []models.Musher
	err = json.Unmarshal([]byte(data), &mushers)
	if err != nil {
		log.Printf("Error unmarshalling data '%v'", err)
		return nil, err
	}

	log.Printf("Data retrieved from cache for key '%s'", cacheKey)
	return mushers, nil
}
