package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/internal/services/getsession"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("useafe: putwall <command>")
		fmt.Println("commands: get_session")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get_session":
		ctx, canel := context.WithTimeout(context.Background(), 10*time.Second)
		defer canel()

		openf1Client := openf1.New()
		service := getsession.New(openf1Client)
		s, err := service.GetSession(ctx, "Belgium", "Sprint", "2023")
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if s == nil {
			fmt.Println("No upcoming session found.")
			return
		}

		fmt.Printf("%s - %s (%s)\n", s.SessionName, s.CircuitName, s.CountryName)
		fmt.Printf("Starts: %s (UTC)\n", s.DateStart.Format(time.RFC1123))
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
