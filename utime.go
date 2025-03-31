// To build: go build -o utime_go utime_go
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Function to display usage information
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-epoch] [-beat] [@<beat_time> | <epoch_time>]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Displays current UTC and local time information by default.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -epoch       Print the current Unix epoch time.\n")
	fmt.Fprintf(os.Stderr, "  -beat        Print the current Swatch Internet Time (@beats).\n")
	fmt.Fprintf(os.Stderr, "  <epoch_time> Convert the given Unix epoch time (integer seconds) to local time.\n")
	fmt.Fprintf(os.Stderr, "  @<beat_time> Convert the given Swatch Internet Time (0-999) to a local time range for the current day.\n")
	fmt.Fprintf(os.Stderr, "  -h, --help   Show this help message.\n") // Added for clarity, flag pkg handles it
	flag.PrintDefaults()                                               // Print flag package default help info
	os.Exit(1)
}

// Function to print current UTC and Local time info
func printCurrentTimes() {
	// Get the current time in UTC
	utcNow := time.Now().UTC()

	// Get the current time in the system's local timezone
	localNow := time.Now() // time.Now() returns time in local zone
	localLoc := localNow.Location()
	zoneName, zoneOffset := localNow.Zone()
	offsetHours := zoneOffset / 3600

	// Check DST - Go handles this implicitly via the Zone offset/name.
	// We can show the offset for clarity.

	fmt.Println("Current Time Information:")
	fmt.Println("==============================")
	fmt.Printf("UTC Time      : %s\n", utcNow.Format("2006-01-02 15:04:05 MST"))   // MST will be UTC here
	fmt.Printf("Local Time    : %s\n", localNow.Format("2006-01-02 15:04:05 MST")) // MST will be local zone abbr.
	fmt.Printf("Time Zone     : %s (UTC%+d)\n", zoneName, offsetHours)
	// Note: Go's time doesn't expose a simple boolean DST flag like Python's pytz/datetime.
	// The zone name (e.g., EDT vs EST) and offset reflect DST status.
	// For a more explicit check, one might compare zone offsets at different times of year,
	// but relying on the zone name/offset from time.Now() is standard practice.
	fmt.Printf("Location      : %s\n", localLoc.String())
	fmt.Println("==============================")
}

// Function to calculate and print Swatch Internet Time (Beat Time)
func printBeatTime() {
	utcSeconds := time.Now().UTC().Unix()
	// BMT is UTC+1
	secondsSinceMidnightBMT := (utcSeconds + 3600) % 86400            // Seconds since midnight BMT
	beatTime := (float64(secondsSinceMidnightBMT) * 1000.0) / 86400.0 // Calculate beats (float for precision)
	fmt.Printf("@%06.2f\n", beatTime)                                 // Format like @###.##
}

// Function to print the current epoch time
func printEpochTime() {
	fmt.Println(time.Now().Unix())
}

// Function to convert epoch to local time
func convertEpochToLocal(epochStr string) {
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid epoch time '%s'. Please provide an integer.\n", epochStr)
		usage()
		return // Keep linters happy
	}

	// Create time object from epoch (assumed UTC)
	t := time.Unix(epoch, 0).UTC()

	// Convert to local time
	localTime := t.Local() // Converts using system's local timezone

	fmt.Println(localTime.Format("2006-01-02 15:04:05 MST")) // Format includes timezone abbr.
}

// Function to convert beat time to a local time range for the current day
func convertBeatToLocalRange(beatStr string) {
	beatInput := strings.TrimPrefix(beatStr, "@")
	beat, err := strconv.Atoi(beatInput)
	if err != nil || beat < 0 || beat > 999 {
		fmt.Fprintf(os.Stderr, "Error: Invalid beat time '%s'. Must be '@' followed by 0-999.\n", beatStr)
		usage()
		return // Keep linters happy
	}

	// Define BMT location (UTC+1)
	bmtLoc := time.FixedZone("BMT", 3600)
	// Get local location
	localLoc := time.Now().Location()

	// Calculate start and end seconds offset from midnight BMT
	// 1 beat = 86.4 seconds
	startOffsetSec := float64(beat) * 86.4
	endOffsetSec := float64(beat+1) * 86.4

	// Get current time in BMT to determine the current BMT date
	nowBMT := time.Now().UTC().In(bmtLoc)
	// Get midnight BMT for the current BMT day
	midnightBMT := time.Date(nowBMT.Year(), nowBMT.Month(), nowBMT.Day(), 0, 0, 0, 0, bmtLoc)

	// Calculate start and end times in BMT by adding offset to midnight BMT
	// Use AddFloat for seconds to handle potential floats if needed, though Add works with Duration
	startTimeBMT := midnightBMT.Add(time.Duration(startOffsetSec * float64(time.Second)))
	// The end time is the start of the *next* beat, representing an exclusive upper bound.
	endTimeBMT := midnightBMT.Add(time.Duration(endOffsetSec * float64(time.Second)))

	// Convert BMT times to local time zone
	startTimeLocal := startTimeBMT.In(localLoc)
	endTimeLocal := endTimeBMT.In(localLoc)

	fmt.Printf("@%d corresponds to the time range %s - %s %s (on %s)\n",
		beat,
		startTimeLocal.Format("15:04:05.000"), // Include milliseconds for precision
		endTimeLocal.Format("15:04:05.000"),
		startTimeLocal.Format("MST"),        // Local timezone abbreviation
		startTimeLocal.Format("2006-01-02")) // Local date
}

func main() {
	// Define flags
	epochFlag := flag.Bool("epoch", false, "Print the current Unix epoch time.")
	beatFlag := flag.Bool("beat", false, "Print the current Swatch Internet Time (@beats).")
	helpFlag := flag.Bool("h", false, "Show help message.") // Explicit help flag

	// Customize usage message
	flag.Usage = usage

	// Parse command-line arguments
	flag.Parse()

	// Handle help flag explicitly
	if *helpFlag {
		usage()
		return
	}

	args := flag.Args() // Get non-flag arguments

	// --- Input Validation and Logic ---

	// Case 1: More than one flag set or flags with arguments
	if (*epochFlag && *beatFlag) || ((*epochFlag || *beatFlag) && len(args) > 0) {
		fmt.Fprintln(os.Stderr, "Error: Cannot combine -epoch or -beat flags with each other or with other arguments.")
		usage()
		return
	}
	// Case 2: More than one positional argument
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: Too many arguments.")
		usage()
		return
	}

	// --- Action Dispatch ---

	if *epochFlag {
		// Action: Print current epoch time
		printEpochTime()
	} else if *beatFlag {
		// Action: Print current beat time
		printBeatTime()
	} else if len(args) == 1 {
		arg := args[0]
		if strings.HasPrefix(arg, "@") {
			// Action: Convert beat time argument to local range
			convertBeatToLocalRange(arg)
		} else {
			// Action: Convert epoch time argument to local time
			convertEpochToLocal(arg)
		}
	} else {
		// Default Action: No flags, no args - print current times
		printCurrentTimes()
	}
}
