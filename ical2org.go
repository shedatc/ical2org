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
	eventCounter := 0
	noUidCounter := 0
	for _, el := range c.Events {
		node := el.Data
		// value, ok := node["DTSTART"]
		// log.Output(0, fmt.Sprintf("value=%v ok=%v", value, ok))
		// if ok {
		// 	dtstart, err := value.DateDecode()
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		// value, ok = node["DTEND"]
		// log.Output(0, fmt.Sprintf("value=%v ok=%v", value, ok))
		// if ok {
		// 	dtend, err := value.DateDecode()
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		value, ok := node["UID"]
		if !ok {
			noUidCounter++
		}
		uid := value.Val

		value, ok = node["SUMMARY"]
		if !ok {
			log.Printf("uid=%v status=error what=\"missing/invalid summary\"\n", uid)
			continue
		}
		summary := value.Val
		log.Printf("uid=%v status=ok summary=\"%v\"\n", uid, summary)

		// A valid event must have at least both an UID and a summary.

		eventCounter++

		fmt.Printf("* %s\n", summary)
		fmt.Println("  :PROPERTIES")
		fmt.Printf("  :UID: %s\n", uid)
		fmt.Println("  :END")
	}
	log.Printf("stats: events=%v no-uid=%v\n", eventCounter, noUidCounter)
	return nil
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
