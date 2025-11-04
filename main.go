package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	inlineMode   = flag.Bool("inline", false, "run in inline mode (disable fullscreen TUI)")
	inlineModeS  = flag.Bool("i", false, "run in inline mode (shorthand for -inline)")
	showVersion  = flag.Bool("version", false, "display version information")
	showVersionS = flag.Bool("v", false, "display version information (shorthand for -version)")
	pausedMode   = flag.Bool("paused", false, "start timer in paused state")
	pausedModeS  = flag.Bool("p", false, "start timer in paused state (shorthand for -paused)")
	timerName    = flag.String("name", "", "name for the timer")
)

func usage() {
	fmt.Fprintf(os.Stderr, "timer - minimal tui countdown/timer app under 5mb memory usage \n\n")
	fmt.Fprintf(os.Stderr, "Usage: timer [options] [<duration>]\n\n")
	fmt.Fprintf(os.Stderr, "Duration: number with unit (5s, 2m, 1h). No unit defaults to seconds.\n")
	fmt.Fprintf(os.Stderr, "          If omitted, runs as a counter (stopwatch) counting up from 00:00.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  timer                    # counter/stopwatch (counts up)\n")
	fmt.Fprintf(os.Stderr, "  timer 5                  # 5 seconds countdown (fullscreen)\n")
	fmt.Fprintf(os.Stderr, "  timer 2m                 # 2 minutes countdown (fullscreen)\n")
	fmt.Fprintf(os.Stderr, "  timer -i 30s             # inline mode countdown\n")
	fmt.Fprintf(os.Stderr, "  timer -p 5m              # 5 minutes countdown starting paused\n")
	fmt.Fprintf(os.Stderr, "  timer -name \"Pomodoro\" 25m  # named timer with notification\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Get args after initial flag parse
	args := flag.Args()

	// Separate flags and positional from remaining args
	var positional []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if arg == "-i" || arg == "-inline" {
				*inlineMode = true
			} else if arg == "-p" || arg == "-paused" {
				*pausedMode = true
			} else if arg == "-v" || arg == "-version" {
				*showVersion = true
			} else {
				usage()
				os.Exit(1)
			}
		} else {
			positional = append(positional, arg)
		}
	}

	// Handle version
	if *showVersion || *showVersionS {
		fmt.Println("timer version", version)
		return
	}

	// Accept 0 or 1 positional arg
	if len(positional) > 1 {
		usage()
		os.Exit(1)
	}

	// Parse duration (0 means counter mode)
	var duration time.Duration
	if len(positional) == 0 {
		// Counter mode - use 0 duration as signal
		duration = 0
	} else {
		durStr := positional[0]
		addSuffixIfArgIsNumber(&durStr, "s")
		var err error
		duration, err = time.ParseDuration(durStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid duration %q\n", positional[0])
			os.Exit(1)
		}
	}

	// Merge short/long flags - fullscreen is default, inline disables it
	useInline := *inlineMode || *inlineModeS
	initialPaused := *pausedMode || *pausedModeS

	// Channel for timer summary
	summaryCh := make(chan TimerSummary, 1)

	// Run timer (fullscreen unless inline flag is set)
	if err := runTimer(duration, !useInline, initialPaused, *timerName, summaryCh); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Receive and print summary
	summary := <-summaryCh
	if summary.Name != "" {
		fmt.Printf("Name: %s\n", summary.Name)
	}
	fmt.Printf("Start: %s\n", summary.Start.Format("2006-01-02 15:04:05"))
	fmt.Printf("End: %s\n", summary.End.Format("2006-01-02 15:04:05"))
	fmt.Printf("Duration: %s\n", summary.Duration)
	fmt.Printf("Mode: %s\n", summary.Mode)
	fmt.Printf("Finished: %t\n", summary.Finished)
}
