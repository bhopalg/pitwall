package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/internal/services/getsession"
	"github.com/bhopalg/pitwall/internal/services/latest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("useafe: putwall <command>")
		fmt.Println("commands: get_session")
		os.Exit(1)
	}

	getSessionCmd := flag.NewFlagSet("get_session", flag.ExitOnError)

	switch os.Args[1] {
	case "latest":
		ctx, canel := context.WithTimeout(context.Background(), 10*time.Second)
		defer canel()

		openf1Client := openf1.New()
		service := latest.New(openf1Client)
		s, err := service.Next(ctx)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if s == nil {
			fmt.Println("No session found.")
			return
		}

		fmt.Printf("%s - %s (%s)\n", s.SessionName, s.CircuitName, s.CountryName)
		fmt.Printf("Starts: %s (UTC)\n", s.DateStart.Format(time.RFC1123))

	case "get_session":
		ctx, canel := context.WithTimeout(context.Background(), 10*time.Second)
		defer canel()
		country := getSessionCmd.String("country", "Belgium", "country name for session")
		session_type := getSessionCmd.String("type", "Sprint", "session type e.g. Sprint, Race")
		session_year := getSessionCmd.String("year", "2023", "session year")

		getSessionCmd.Parse(os.Args[2:])

		openf1Client := openf1.New()
		service := getsession.New(openf1Client)
		s, err := service.GetSession(ctx, *country, *session_type, *session_year)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if s == nil {
			fmt.Println("No sessions found.")
			return
		}

		fmt.Printf("%s - %s (%s)\n", s.SessionName, s.CircuitName, s.CountryName)
		fmt.Printf("Starts: %s (UTC)\n", s.DateStart.Format(time.RFC1123))
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
