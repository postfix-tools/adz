package main

import (
	"fmt"
	"log"
	"os"
	"time"

	parser "github.com/postfix-tools/chisel"
	"github.com/urfave/cli"
)

var (
	connects      int
	lostconnects  int
	disconnects   int
	qidcount      int
	sends         int
	deferrals     int
	bounces       int
	local_pickups int
	app           *cli.App
	BUILD         string
	VERSION       string
)

func main() {
	app = cli.NewApp()
	app.Name = "adz"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "age,a",
			Usage:  "Max age to consider a record valid for the run, in minutes",
			EnvVar: "ME_MAXAGE",
			Value:  1,
		},
		cli.StringFlag{
			Name:   "file,f",
			Usage:  "Logfile to run against",
			EnvVar: "ME_LOGFILE",
			Value:  "/var/log/maillog",
		},
	}
	app.Action = handleLog
	app.Version = VERSION
	app.Usage = "Perform summary analyis of Postfix Mail Event logs"
	app.Run(os.Args)
}

func handleLog(c *cli.Context) error {
	var store parser.LogStore
	store.Filename = c.String("file")
	pstart := time.Now()
	store.ParseLogFile()
	elapsed := time.Now().Sub(pstart)
	fmt.Printf("Logparse took %q\n", elapsed)

	age_minutes := c.Int("age")
	max_age, err := time.ParseDuration(fmt.Sprintf("%dm", age_minutes))
	if err != nil {
		log.Printf("Invalid duration value.")
		log.Fatal(err)
	}
	events := store.GetRecords(max_age)
	for _, e := range events {
		//log.Printf("%+v\n", e)
		switch e.GetComponent() {
		case "smtpd":
			switch e.GetRecordType() {
			case 0:
				connects++
			case 1:
				disconnects++
			case 2:
				lostconnects++
			}
		case "pickup":
			local_pickups++
		case "cleanup":
			// don't do anything yet with these
		case "smtp":
			switch e.(*parser.SMTPRecord).Status {
			case "sent":
				sends++
			case "deferred":
				deferrals++
			case "bounced":
				bounces++
			case "":
				// represents a non-state smtp entry
			default:
				fmt.Printf("unhandled status entry: %q\n", e.(*parser.SMTPRecord).Status)
			}
		case "qmgr":
			// don't do anything yet with these
		case "local":
			switch e.(*parser.DeliveryRecord).Status {
			case "sent":
				sends++
			case "deferred":
				deferrals++
			case "bounced":
				bounces++
			default:
				fmt.Printf("unhandled status entry: %q\n", e.(*parser.SMTPRecord).Status)
			}

		}
	}
	fmt.Printf("Recorded %d Connects\n", connects)
	fmt.Printf("Recorded %d Lost Connections\n", lostconnects)
	fmt.Printf("Recorded %d Disconnects\n", disconnects)
	fmt.Printf("Recorded %d Bounces\n", bounces)
	fmt.Printf("Recorded %d Local Pickups\n", local_pickups)
	fmt.Printf("Recorded %d Sends\n", sends)
	fmt.Printf("Recorded %d Deferrals\n", deferrals)
	return nil
}
