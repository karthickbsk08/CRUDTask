// Package catching provides utility functions for interacting with a Redis cache,
// including creating a Redis client, and performing basic cache operations such as
// getting, setting, and deleting cache entries.
package catching

import (
	"context"
	"tasks/helpers"
	"time"

	"github.com/redis/go-redis/v9"
)

var GRedisClient *redis.Client

// CreateCacheClient initializes a Redis client using the provided Redis address.
// It logs the connection status using the provided HelperStruct for debugging.
// If the connection fails, it logs the error and returns without creating the client.
//
// Parameters:
//   - pDebug: Pointer to a HelperStruct used for logging.
//   - pRedisAddr: Address of the Redis server (e.g., "localhost:6379").
//
// Note: This function does not return the Redis client; it is intended for connection testing and logging.
func CreateCacheClient(pDebug *helpers.HelperStruct, pRedisAddr string) {
	pDebug.Log(helpers.Statement, "CreateCacheClient(+), pRedisAddr : ", pRedisAddr)

	var ctx = context.Background()

	// Create a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: pRedisAddr, // Adjust the address as needed
	})
	GRedisClient = rdb

	//Check connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		pDebug.Log(helpers.Statement, "Could not connect to Redis:", err)
		return
	}

	pDebug.Log(helpers.Statement, "Connected to Redis successfully!")
	pDebug.Log(helpers.Statement, "CreateCacheClient(-)")
	// You can now use rdb to interact with Redis
}

// GetFromCache retrieves the value associated with the given key from the Redis cache.
// It logs the operation and returns the value as a string. If the key does not exist,
// it returns an empty string and no error. If another error occurs, it is logged and returned.
//
// Parameters:
//   - pDebug: Pointer to a HelperStruct used for logging.
//   - rdb: Redis client instance.
//   - key: The key to retrieve from the cache.
//
// Returns:
//   - string: The value associated with the key, or an empty string if the key does not exist.
//   - error: An error if the operation fails, otherwise nil.
func GetFromCache(pDebug *helpers.HelperStruct, rdb *redis.Client, ctx context.Context, key string) (string, error) {
	pDebug.Log(helpers.Statement, "GetFromCache(+), key : ", key)
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		pDebug.Log(helpers.Statement, "Key does not exist in cache")
		return "", nil
	} else if err != nil {
		pDebug.Log(helpers.Statement, "Error fetching from cache:", err)
		return "", err
	}

	pDebug.Log(helpers.Statement, "GetFromCache(-), value : ", val)
	return val, nil
}

// SetToCache sets a key-value pair in the Redis cache with the specified expiration duration.
// It logs the operation and returns an error if the operation fails.
//
// Parameters:
//   - pDebug: Pointer to a HelperStruct used for logging.
//   - rdb: Redis client instance.
//   - key: The key to set in the cache.
//   - value: The value to associate with the key.
//   - expiration: The expiration duration for the cache entry.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func SetToCache(pDebug *helpers.HelperStruct, rdb *redis.Client, ctx context.Context, key string, value string, expiration time.Duration) error {
	pDebug.Log(helpers.Statement, "SetToCache(+), key : ", key, ", value : ", value, ", expiration : ", expiration)

	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		pDebug.Log(helpers.Statement, "Error setting cache:", err)
		return err
	}

	pDebug.Log(helpers.Statement, "SetToCache(-)")
	return nil
}

// DeleteFromCache removes the specified key from the Redis cache.
// It logs the operation and returns an error if the operation fails.
//
// Parameters:
//   - pDebug: Pointer to a HelperStruct used for logging.
//   - rdb: Redis client instance.
//   - key: The key to delete from the cache.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func DeleteFromCache(pDebug *helpers.HelperStruct, rdb *redis.Client, key string) error {
	pDebug.Log(helpers.Statement, "DeleteFromCache(+), key : ", key)
	var ctx = context.Background()
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		pDebug.Log(helpers.Statement, "Error deleting from cache:", err)
		return err
	}

	pDebug.Log(helpers.Statement, "DeleteFromCache(-)")
	return nil
}
