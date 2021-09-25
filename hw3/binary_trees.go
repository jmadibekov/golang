package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func walkDfs(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	walkDfs(t.Left, ch)
	ch <- t.Value
	walkDfs(t.Right, ch)
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	walkDfs(t, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for v1 := range ch1 {
		v2, ok := <-ch2

		if !ok || v1 != v2 {
			return false
		}
	}

	return true
}

func main() {
	ch := make(chan int)
	go Walk(tree.New(5), ch)

	fmt.Println("Testing Walk function")
	for v := range ch {
		fmt.Println(v)
	}

	fmt.Println("Testing Same function")
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
