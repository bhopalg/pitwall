package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bhopalg/pitwall/domain"
	"github.com/bhopalg/pitwall/internal/cache"
	"github.com/bhopalg/pitwall/internal/openf1"
	"github.com/bhopalg/pitwall/internal/services/getsession"
	"github.com/bhopalg/pitwall/internal/services/latest"
	"github.com/bhopalg/pitwall/internal/services/remind"
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
	fileCache := &cache.FileCache{Dir: "../../.pitwall_cache"}

	getSessionCmd := flag.NewFlagSet("get_session", flag.ExitOnError)

	ctx, canel := context.WithTimeout(context.Background(), 10*time.Second)
	defer canel()

	openf1Client := openf1.New()

	switch os.Args[1] {
	case "remind":
		remindCmd := flag.NewFlagSet("remind", flag.ExitOnError)
		threshold := remindCmd.Int("minutes", 30, "minutes threshold for reminder")
		quiet := remindCmd.Bool("quiet", false, "suppress output if no reminder")

		remindCmd.Parse(os.Args[2:])

		service := latest.New(openf1Client, fileCache)
		res, err := service.Next(ctx)

		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(2)
		}

		if res.Session == nil {
			if !*quiet {
				fmt.Println("No upcoming sessions found.")
			}
			os.Exit(1)
		}

		trigger, diff := remind.ShouldRemind(now, res.Session.DateStart, *threshold)

		if trigger {
			fmt.Printf("REMIND: %s starts in %v!\n", res.Session.SessionName, diff.Round(time.Minute))
			os.Exit(0)
		}

		if !*quiet {
			fmt.Printf("No reminder needed. Next session (%s) is in %v.\n",
				res.Session.SessionName, diff.Round(time.Minute))
		}
		os.Exit(1)
	case "cache":
		if len(os.Args) < 3 {
			fmt.Println("usage: pitwall cache <info|clear>")
			return
		}

		subCommand := os.Args[2]

		switch subCommand {
		case "info":
			entries, path, err := fileCache.Info()
			if err != nil {
				fmt.Printf("Error reading cache: %v\n", err)
				return
			}

			fmt.Printf("Cache Location: %s\n", path)
			fmt.Printf("Total Entries:  %d\n\n", len(entries))

			if len(entries) > 0 {
				fmt.Printf("%-30s %-20s %-10s %-10s\n", "KEY", "CREATED AT", "STALE", "SIZE")
				for _, e := range entries {
					staleStr := "no"
					if e.IsStale {
						staleStr = "YES"
					}
					fmt.Printf("%-30s %-20s %-10s %-10d B\n",
						e.Key,
						e.CreatedAt.Format("02 Jan 15:04"),
						staleStr,
						e.Size,
					)
				}
			}
		case "clear":
			count, err := fileCache.Clear()
			if err != nil {
				fmt.Printf("Error clearing cache: %v\n", err)
				return
			}
			if count == 0 {
				fmt.Println("Nothing to clear.")
			} else {
				fmt.Printf("Successfully removed %d cache entries.\n", count)
			}
		default:
			fmt.Printf("unknown cache command: %s\n", subCommand)
		}

	case "weekend":
		country := getSessionCmd.String("country", "Belgium", "country name for session")
		session_year := getSessionCmd.String("year", "2023", "session year")

		getSessionCmd.Parse(os.Args[2:])

		service := weekend.New(openf1Client, fileCache)
		sessions, err := service.Weekend(ctx, *country, *session_year)

		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if sessions.Sessions != nil && sessions.Warning != "" {
			fmt.Println(sessions.Warning)
		}

		if sessions.Sessions == nil || len(*sessions.Sessions) == 0 {
			fmt.Println("No sessions found.")
			return
		}

		firstSession := (*sessions.Sessions)[0]
		fmt.Printf("%s Grand Prix - %s\n\n", firstSession.CountryName, firstSession.CircuitName)

		groupSessions := createWeekendGroup(*sessions.Sessions)

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
		service := latest.New(openf1Client, fileCache)
		s, err := service.Next(ctx)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if s.Session == nil {
			fmt.Println("No session found.")
			return
		}

		fmt.Printf("%s - %s (%s)\n", s.Session.SessionName, s.Session.CircuitName, s.Session.CountryName)
		utils.PrintSessionStatus(s.Session, now)

	case "get_session":
		country := getSessionCmd.String("country", "Belgium", "country name for session")
		session_type := getSessionCmd.String("type", "Sprint", "session type e.g. Sprint, Race")
		session_year := getSessionCmd.String("year", "2023", "session year")

		getSessionCmd.Parse(os.Args[2:])

		service := getsession.New(openf1Client, fileCache)
		s, err := service.GetSession(ctx, *country, *session_type, *session_year)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if s.Session != nil && s.Warning != "" {
			fmt.Println(s.Warning)
		}

		if s.Session == nil {
			fmt.Println("No sessions found.")
			return
		}

		if s.Session != nil && s.Warning != "" {
			fmt.Println(s.Warning)
		}

		fmt.Printf("%s - %s (%s)\n", s.Session.SessionName, s.Session.CircuitName, s.Session.CountryName)
		fmt.Printf("Starts: %s (UTC)\n", s.Session.DateStart.Format(time.RFC1123))
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
