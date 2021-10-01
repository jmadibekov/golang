package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func main() {
	// binary_trees.go
	ch := make(chan int)
	go Walk(tree.New(5), ch)

	fmt.Println("Testing Walk function")
	for v := range ch {
		fmt.Println(v)
	}

	fmt.Println("Testing Same function")
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))

	// web_crawler.go
	fetched := FetchedUrls{urls: make(map[string]bool)}
	fmt.Println("Starting to crawl")

	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher, &fetched)
	wg.Wait()
}
