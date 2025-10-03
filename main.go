package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	inlineMode   = flag.Bool("inline", false, "run in inline mode (disable fullscreen TUI)")
	inlineModeS  = flag.Bool("i", false, "run in inline mode (shorthand for -inline)")
	showVersion  = flag.Bool("version", false, "display version information")
	showVersionS = flag.Bool("v", false, "display version information (shorthand for -version)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "timer - minimal tui countdown/timer app under 5mb memory usage \n\n")
	fmt.Fprintf(os.Stderr, "Usage: timer [options] [<duration>]\n\n")
	fmt.Fprintf(os.Stderr, "Duration: number with unit (5s, 2m, 1h). No unit defaults to seconds.\n")
	fmt.Fprintf(os.Stderr, "          If omitted, runs as a counter (stopwatch) counting up from 00:00.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  timer            # counter/stopwatch (counts up)\n")
	fmt.Fprintf(os.Stderr, "  timer 5          # 5 seconds countdown (fullscreen)\n")
	fmt.Fprintf(os.Stderr, "  timer 2m         # 2 minutes countdown (fullscreen)\n")
	fmt.Fprintf(os.Stderr, "  timer -i 30s     # inline mode countdown\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Handle version
	if *showVersion || *showVersionS {
		fmt.Println("timer version", version)
		return
	}

	// Accept 0 or 1 positional arg
	args := flag.Args()
	if len(args) > 1 {
		usage()
		os.Exit(1)
	}

	// Parse duration (0 means counter mode)
	var duration time.Duration
	if len(args) == 0 {
		// Counter mode - use 0 duration as signal
		duration = 0
	} else {
		durStr := args[0]
		addSuffixIfArgIsNumber(&durStr, "s")
		var err error
		duration, err = time.ParseDuration(durStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid duration %q\n", args[0])
			os.Exit(1)
		}
	}

	// Merge short/long flags - fullscreen is default, inline disables it
	useInline := *inlineMode || *inlineModeS

	// Run timer (fullscreen unless inline flag is set)
	if err := runTimer(duration, !useInline); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
