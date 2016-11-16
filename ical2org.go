package main

import "github.com/laurent22/ical-go/ical"
import "io/ioutil"
import "log"
import "os"

func main() {
	// icalFilePath := "sample-ics-files/bastille.ics"
	// icalFilePath := "sample-ics-files/small.ics"
	icalFilePath := "sample-ics-files/broken.ics"
	icalData, err := ioutil.ReadFile(icalFilePath)
	if err != nil {
		log.Fatalf("Unable to read file: %s", icalFilePath)
		os.Exit(65) // EX_DATAERR, see FreeBSD's sysexits(3),
	}
	// log.Println("Data:")
	// log.Println(icalData)
	icalNode, err := ical.ParseCalendar(string(icalData))
	if err != nil {
		log.Fatalf("Unable to read file: %s", icalFilePath)
		os.Exit(65) // EX_DATAERR, see FreeBSD's sysexits(3),
	}
	log.Println("Node:")
	log.Println(icalNode)
}
