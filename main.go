package main

import "fmt"

func main() {
	fmt.Println(greet("Go"))
}

func greet(name string) string {
	return fmt.Sprintf("Hello from %s!", name)
}
