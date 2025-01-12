package main

import (
	"fmt"
	"time"

	"github.com/go-jedi/go-test/pkg/cache"
)

type testData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// setAndGetToCache set and get data from cache.
func setAndGetToCache() {
	defaultTTL := 10
	td := testData{
		Name: "test",
		Age:  10,
	}

	c := cache.NewCache(time.Duration(defaultTTL) * time.Second)

	c.Set("test", td)

	var res testData
	if !c.Get("test", &res) {
		fmt.Println("data from cache by key is not found or expired")
		return
	}

	fmt.Println("res:", res)
}

// deleteDataFromCache delete data by key from cache.
func deleteDataFromCache() {
	defaultTTL := 10
	td := testData{
		Name: "test",
		Age:  10,
	}

	c := cache.NewCache(time.Duration(defaultTTL) * time.Second)

	c.Set("test", td)

	c.Delete("test")

	var res testData
	if !c.Get("test", &res) {
		fmt.Println("data is not found by key")
		return
	}

	fmt.Println("res:", res)
}

// addCleanupInCache add cleanup in cache.
func addCleanupInCache() {
	defaultTTL := 10
	td := testData{
		Name: "test",
		Age:  10,
	}

	c := cache.NewCache(time.Duration(defaultTTL) * time.Second)

	c.StartCleanup(5 * time.Second)

	c.Set("test", td, 4*time.Second)

	var res testData
	if !c.Get("test", &res) {
		fmt.Println("data is not found by key")
		return
	}

	time.Sleep(6 * time.Second)

	var resTwo testData
	if !c.Get("test", &resTwo) {
		fmt.Println("data is not found by key")
		return
	}

	fmt.Println("resTwo:", resTwo)
}

// checkExpiredInCache check expired in cache.
func checkExpiredInCache() {
	defaultTTL := 10
	td := testData{
		Name: "test",
		Age:  10,
	}

	c := cache.NewCache(time.Duration(defaultTTL) * time.Second)

	c.Set("test", td)

	time.Sleep(11 * time.Second)

	var res testData
	if !c.Get("test", &res) {
		fmt.Println("data is not found by key")
		return
	}

	fmt.Println("res:", res)
}

func main() {
	setAndGetToCache()
	deleteDataFromCache()
	addCleanupInCache()
	checkExpiredInCache()
}
