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

const TwoWeeks = 24 * 14 // 2 weeks = 24x14 = 336 hours

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
		start, err := value.DateDecode()
		if err != nil {
			log.Printf("ERROR error-msg=\"unable to decode start date (DTSTART)\" uid=%v\n", uid)
			continue
		}

		value, ok = node["SUMMARY"]
		if !ok {
			log.Printf("ERROR error-msg=\"missing summary\" uid=%v\n", uid)
			continue
		}
		summary := value.Val

		log.Printf("OK summary=\"%v\" start=\"%v\" since-start=%v uid=%v\n",
			summary, start, time.Since(start), uid)

		// A valid event must have at least:
		// - an UID;
		// - a start date (DTSTART);
		// and a summary.

		eventCounter++

		var end time.Time
		value, haveEndDate := node["DTEND"]
		if haveEndDate {
			end, err = value.DateDecode()
			if err != nil {
				log.Printf("- error-msg=\"unable to decode end date (DTEND)\"\n")
				continue
			}
			log.Printf("- end=\"%v\" since-end=%v\n", end, time.Since(end))
		}

		// XXX Discard the event if it's out of our time window.
		//       s       e
		// -n----[-------]------
		//   <--> Must be less than 2 weeks in the future.
		//
		// ------[-n-----]------
		//
		// ------[-------]-----n
		//                <---> Must be less than 2 weeks in the past.

		fmt.Printf("\n")
		fmt.Printf("* %s\n", summary)
		if haveEndDate {
			fmt.Printf("  %s\n", orgDateRange(start, end))
		} else {
			fmt.Printf("  %s\n", orgDate(start))
		}

		fmt.Printf("  :PROPERTIES:\n")
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
	return fmt.Sprintf("<%04d-%02d-%02d %s %02d:%02d>",
		dt.Year(),
		dt.Month(),
		dt.Day(),
		dt.Weekday().String()[0:3],
		dt.Hour(),
		dt.Minute(),
	)
}

func orgDateRange(start time.Time, end time.Time) string {
	return fmt.Sprintf("%s--%s", orgDate(start), orgDate(end))
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
