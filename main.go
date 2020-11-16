package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strconv"

	"github.com/deepesh15/Avatar-Me/generator"
)

func createHash(userInput string) []byte {
	hash := sha1.New()
	hash.Write([]byte(userInput))
	createdHash := hash.Sum(nil)
	return createdHash
}

func main() {
	args := os.Args

	if len(args[1:]) != 3 {
		fmt.Println("BAD INPUT\n" + "run avatar <word> <size greater than 128> <output filename>")
		return
	}
	wordToHash := args[1]
	sideStr := args[2]
	fileName := args[3]
	side, err := strconv.Atoi(sideStr)
	if err != nil || side < 128 {
		fmt.Println("BAD INPUT\n + Make sure the side value is greater than 128")
		return
	}
	icon := generator.New(side)
	errCreating := icon.Create(createHash(wordToHash), fileName)
	if errCreating != nil {
		fmt.Printf("Error occurred: %v", err)
	}
	fmt.Println("Operation successfull !!")
}
