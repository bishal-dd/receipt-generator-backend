package main

import (
	"fmt"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
)

func InitializeApi()  {
	API_KEY := os.Getenv("CLERK_API")
	fmt.Print(API_KEY)
	clerk.SetKey(API_KEY)
}