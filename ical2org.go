package main

import "github.com/jordic/goics"
import "io/ioutil"
import "log"
import "os"
import "strings"
import "time"
import "fmt"

type Events struct{}

type TimeFrame struct {
	Start, End time.Time
}

// Tell if f totally or partially overlap g.
//
// Legend: f [] g <>
//
//  ---[---]---<------->----- false
//  ---<--->-----[-------]--- false
//  ---<-----[------>------]- true
//  -----[--<--------]--->--- true
func (f TimeFrame) overlap(g TimeFrame) bool {
	return f.End.After(g.Start) && f.Start.Before(g.End)
}

func (f TimeFrame) contains(t time.Time) bool {
	return f.Start.Before(t) && f.End.After(t)
}

const TwoWeeks = 336 * time.Hour // 2 weeks = 2 x 7 x 24 hours

func (e *Events) ConsumeICal(c *goics.Calendar, err error) error {
	fmt.Println("-*- eval: (auto-revert-mode 1); -*-")

	totalCounter := 0
	processedCounter := 0
	noUidCounter := 0
	noStartDateCounter := 0
	noSummaryCounter := 0
	invalidStartDateCounter := 0
	invalidEndDateCounter := 0
	outOfTimeFrameCounter := 0

	currentTimeFrame := TimeFrame{Start: time.Now().Add(-TwoWeeks), End: time.Now().Add(TwoWeeks)}
	log.Printf("CONF currentTimeFrame={Start: %v, End: %v}\n",
		currentTimeFrame.Start, currentTimeFrame.End)
	for _, el := range c.Events {
		node := el.Data

		totalCounter++

		value, ok := node["UID"]
		if !ok {
			log.Printf("ERROR error-msg=\"missing UID\"\n")
			noUidCounter++
			continue
		}
		uid := value.Val
		log.Printf("EVENT uid=\"%v\"\n", uid)

		value, ok = node["DTSTART"]
		if !ok {
			log.Printf("` ERROR error-msg=\"missing start date (DTSTART)\" uid=%v\n", uid)
			noStartDateCounter++
			continue
		}
		start, err := value.DateDecode()
		if err != nil {
			log.Printf("` ERROR error-msg=\"unable to decode start date (DTSTART)\" uid=%v\n", uid)
			invalidStartDateCounter++
			continue
		}
		log.Printf("| start=\"%v\"\n", start)

		value, ok = node["SUMMARY"]
		if !ok {
			log.Printf("` ERROR error-msg=\"missing summary\" uid=%v\n", uid)
			noSummaryCounter++
			continue
		}
		summary := value.Val

		log.Printf("| summary=\"%v\"\n", summary)

		var end time.Time
		value, haveEndDate := node["DTEND"]
		if haveEndDate {
			end, err = value.DateDecode()
			if err != nil {
				log.Printf("` ERROR error-msg=\"unable to decode end date (DTEND)\"\n")
				invalidEndDateCounter++
				continue
			}
			log.Printf("| end=\"%v\"\n", end)

			eventTimeFrame := TimeFrame{Start: start, End: end}
			if currentTimeFrame.overlap(eventTimeFrame) {
				log.Printf("| KEEP\n")
			} else {
				log.Printf("` DISCARD\n")
				outOfTimeFrameCounter++
				continue
			}
		} else {
			if currentTimeFrame.contains(start) {
				log.Printf("| KEEP\n")
			} else {
				log.Printf("` DISCARD\n")
				outOfTimeFrameCounter++
				continue
			}
		}

		// The event is valid and will be part of the output.
		processedCounter++

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
			log.Printf("| location=\"%v\"\n", location)
			fmt.Printf("  :LOCATION: %s\n", location)
		}

		value, ok = node["STATUS"]
		if ok {
			status := value.Val
			log.Printf("| status=\"%v\"\n", status)
			fmt.Printf("  :STATUS: %s\n", status)
		}

		fmt.Printf("  :END:\n")

		value, ok = node["DESCRIPTION"]
		if ok {
			description := value.Val
			log.Printf("| description=\"%v\"\n", description)
			fmt.Printf("\n  %s\n", description)
		}
		log.Printf("` PROCESSED\n")
	}
	log.Printf("STATS total=%v processed=%v no-uid=%v no-start-date=%v no-summary=%v out-of-time-frame=%v invalid-start-date=%v invalid-end-date=%v\n",
		totalCounter, processedCounter, noUidCounter, noStartDateCounter,
		noSummaryCounter, outOfTimeFrameCounter,
		invalidStartDateCounter, invalidEndDateCounter)
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
