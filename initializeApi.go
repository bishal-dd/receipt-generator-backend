package main

import (
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
)

func InitializeApi()  {
	API_KEY := os.Getenv("CLERK_API")
	clerk.SetKey(API_KEY)
}