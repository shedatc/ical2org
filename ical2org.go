package main

import "github.com/laurent22/ical-go/ical"
import "io/ioutil"
import "log"
import "os"

func main() {
	icalData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Unable to read from stdin")
		os.Exit(65) // EX_DATAERR, see FreeBSD's sysexits(3),
	}
	// log.Println("Data:")
	// log.Println(icalData)
	icalNode, err := ical.ParseCalendar(string(icalData))
	if err != nil {
		log.Fatalf("Parse error")
		os.Exit(65) // EX_DATAERR, see FreeBSD's sysexits(3),
	}
	log.Println("Node:")
	log.Println(icalNode)
}
