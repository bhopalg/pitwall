package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/internal/services/getsession"
	"github.com/bhopalg/pitwall/internal/services/latest"
	"github.com/bhopalg/pitwall/internal/services/weekend"
	"github.com/bhopalg/pitwall/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("useafe: putwall <command>")
		fmt.Println("commands: get_session")
		os.Exit(1)
	}

	now := time.Now().UTC()

	getSessionCmd := flag.NewFlagSet("get_session", flag.ExitOnError)

	ctx, canel := context.WithTimeout(context.Background(), 10*time.Second)
	defer canel()

	openf1Client := openf1.New()

	switch os.Args[1] {
	case "weekend":
		country := getSessionCmd.String("country", "Belgium", "country name for session")
		session_year := getSessionCmd.String("year", "2023", "session year")

		getSessionCmd.Parse(os.Args[2:])

		service := weekend.New(openf1Client)
		sessions, err := service.Weekend(ctx, *country, *session_year)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if sessions == nil || len(*sessions) == 0 {
			fmt.Println("No sessions found.")
			return
		}

		firstSession := (*sessions)[0]
		fmt.Printf("%s Grand Prix - %s\n\n", firstSession.CountryName, firstSession.CircuitName)

		groupSessions := createWeekendGroup(*sessions)

		orderedDays := []string{"Fri", "Sat", "Sun"}

		for _, day := range orderedDays {
			if s, ok := groupSessions[day]; ok {
				fmt.Printf("%s\n", day)
				for _, session := range s {
					fmt.Printf("\t%s\t%s\n", session.SessionName, session.DateStart.Format("15:04"))
				}
				fmt.Println()
			}
		}

	case "latest":
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
		utils.PrintSessionStatus(s, now)

	case "get_session":
		country := getSessionCmd.String("country", "Belgium", "country name for session")
		session_type := getSessionCmd.String("type", "Sprint", "session type e.g. Sprint, Race")
		session_year := getSessionCmd.String("year", "2023", "session year")

		getSessionCmd.Parse(os.Args[2:])

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

func createWeekendGroup(session []domain.Session) map[string][]domain.Session {
	sessionsByDay := make(map[string][]domain.Session)

	for _, s := range session {
		dayKey := s.DateStart.Format("Mon")
		sessionsByDay[dayKey] = append(sessionsByDay[dayKey], s)
	}

	return sessionsByDay
}
