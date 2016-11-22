package main

import "github.com/jordic/goics"
import "io/ioutil"
import "log"
import "os"
import "strings"
import "time"
import "fmt"

type Event struct {
	Start, End  time.Time
	ID, Summary string
}

type Events []Event

func (e *Events) ConsumeICal(c *goics.Calendar, err error) error {
	fmt.Println("-*- eval: (auto-revert-mode 1); -*-")

	eventCounter := 0
	noUidCounter := 0
	for _, el := range c.Events {
		node := el.Data

		value, ok := node["UID"]
		if !ok {
			noUidCounter++
			continue
		}
		uid := value.Val

		value, ok = node["DTSTART"]
		if !ok {
			log.Printf("ERROR error-msg=\"missing start date (dtstart)\" uid=%v\n", uid)
			continue
		}
		dtStart, err := value.DateDecode()
		if err != nil {
			log.Printf("ERROR error-msg=\"unable to decode start date (dtstart)\" uid=%v\n", uid)
			continue
		}

		value, ok = node["SUMMARY"]
		if !ok {
			log.Printf("ERROR error-msg=\"missing summary\" uid=%v\n", uid)
			continue
		}
		summary := value.Val
		log.Printf("OK uid=%v summary=\"%v\" dtStart=%v\n", uid,
			summary, dtStart)

		// A valid event must have at least both an UID and a summary.

		eventCounter++

		fmt.Printf("\n")
		fmt.Printf("* %s\n", summary)
		fmt.Printf("  %s\n", orgDate(dtStart))
		fmt.Printf("  :PROPERTIES\n")
		fmt.Printf("  :ID: %s\n", uid)

		value, ok = node["LOCATION"]
		if ok {
			location := value.Val
			log.Printf("- location=\"%v\"\n", location)
			fmt.Printf("  :LOCATION: %s\n", location)
		}

		value, ok = node["STATUS"]
		if ok {
			status := value.Val
			log.Printf("- status=\"%v\"\n", status)
			fmt.Printf("  :STATUS: %s\n", status)
		}

		fmt.Printf("  :END:\n")

		value, ok = node["DESCRIPTION"]
		if ok {
			description := value.Val
			log.Printf("- description=\"%v\"\n", description)
			fmt.Printf("\n  %s\n", description)
		}
	}
	log.Printf("STATS events=%v no-uid=%v\n", eventCounter, noUidCounter)
	return nil
}

func orgDate(dt time.Time) string {
	return fmt.Sprintf("<%04d-%02d-%02d %02d:%02d>",
		dt.Year(),
		dt.Month(),
		dt.Day(),
		dt.Hour(),
		dt.Minute(),
	)
}

func main() {
	fileContent, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Unable to read from stdin")
		os.Exit(65) // EX_DATAERR, see FreeBSD's sysexits(3),
	}

	decoder := goics.NewDecoder(strings.NewReader(string(fileContent)))
	events := Events{}
	err = decoder.Decode(&events)
	if err != nil {
		log.Fatal("Can't decode stdin")
		os.Exit(65) // EX_DATAERR, see FreeBSD's sysexits(3),
	}
}
