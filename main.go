//
// main.go --- e3db command line tool.
//
// Copyright (C) 2017, Tozny, LLC.
// All Rights Reserved.
//

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jawher/mow.cli"
	"github.com/tozny/e3db-go"
)

type Options struct {
	Logging *bool
	Profile *string
}

func (o *Options) getClient() *e3db.Client {
	var client *e3db.Client
	var err error

	if *o.Profile == "" {
		client, err = e3db.GetDefaultClient()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		client, err = e3db.GetClient(*o.Profile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *o.Logging {
		client.Logging = true
	}

	return client
}

var options Options

func CmdList(cmd *cli.Cmd) {
	data := cmd.BoolOpt("d data", false, "include data in JSON format")
	outputJSON := cmd.BoolOpt("j json", false, "output in JSON format")
	contentTypes := cmd.StringsOpt("t type", nil, "record content type")
	recordIDs := cmd.StringsOpt("r record", nil, "record ID")
	writerIDs := cmd.StringsOpt("w writer", nil, "record writer ID")
	userIDs := cmd.StringsOpt("u user", nil, "record user ID")

	cmd.Action = func() {
		client := options.getClient()

		cursor := client.Query(context.Background(), e3db.Q{
			ContentTypes: *contentTypes,
			RecordIDs:    *recordIDs,
			WriterIDs:    *writerIDs,
			UserIDs:      *userIDs,
			IncludeData:  *data,
		})

		if *outputJSON {
			fmt.Println("[")
		}

		first := true
		for cursor.Next() {
			record := cursor.Get()
			if *outputJSON {
				if first {
					first = false
				} else {
					fmt.Printf(",\n")
				}

				bytes, _ := json.MarshalIndent(record, "  ", "  ")
				fmt.Printf("  %s", bytes)
			} else {
				fmt.Printf("%-40s %s\n", record.Meta.RecordID, record.Meta.Type)
			}
		}

		if *outputJSON {
			fmt.Println("\n]")
		}
	}
}

func CmdWrite(cmd *cli.Cmd) {
}

func CmdRead(cmd *cli.Cmd) {
	recordIDs := cmd.Strings(cli.StringsArg{
		Name:      "RECORD_ID",
		Desc:      "record ID to read",
		Value:     nil,
		HideValue: true,
	})

	cmd.Spec = "RECORD_ID..."
	cmd.Action = func() {
		client := options.getClient()

		for _, recordID := range *recordIDs {
			record, err := client.Get(context.Background(), recordID)
			if err != nil {
				log.Fatal(err)
			}

			bytes, err := json.MarshalIndent(record, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(bytes))
		}
	}
}

func main() {
	app := cli.App("e3db-cli", "E3DB Command Line Interface")

	app.Version("v version", "e3db-cli 0.0.1")

	options.Logging = app.BoolOpt("d debug", false, "enable debug logging")
	options.Profile = app.StringOpt("p profile", "", "e3db configuration profile")

	app.Command("ls", "list records", CmdList)
	app.Command("read", "read records", CmdRead)
	app.Command("write", "write a record", CmdWrite)
	app.Run(os.Args)
}
