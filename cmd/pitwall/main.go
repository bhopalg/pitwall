package main

import (
	"context"
	"fmt"
	"os"
	"time"

	getsession "github.com/bhopalg/pitwall/internal/services/get_session"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("useafe: putwall <command>")
		fmt.Println("commands: next")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get_session":
		ctx, canel := context.WithTimeout(context.Background(), 10*time.Second)
		defer canel()

		s, err := getsession.GetSession(ctx)
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
