package main

import "fmt"

func main() {
	// This is a line comment
	fmt.Println("Hello, World!") // Partial line comment

	/* This is a block comment */
	fmt.Println("Another line")

	/*
		This is a multi-line
		block comment.
	*/

	code := "This is line 13" // code
	
	/* Block start */ fmt.Println("Code between blocks") /* Block end */
}

// Last line comment
