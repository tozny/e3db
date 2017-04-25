//
// main.go --- e3db command line tool.
//
// Copyright (C) 2017, Tozny, LLC.
// All Rights Reserved.
//

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jawher/mow.cli"
	"github.com/tozny/e3db-go"
)

type cliOptions struct {
	Logging *bool
	Profile *string
}

func dieErr(err error) {
	fmt.Fprintf(os.Stderr, "e3db-cli: %s\n", err)
	cli.Exit(1)
}

func dieFmt(format string, args ...interface{}) {
	fmt.Fprint(os.Stderr, "e3db-cli: ")
	fmt.Fprintf(os.Stderr, format, args)
	fmt.Fprint(os.Stderr, "\n")
	cli.Exit(1)
}

func (o *cliOptions) getClient() *e3db.Client {
	var client *e3db.Client
	var err error

	opts, err := e3db.GetConfig(*o.Profile)
	if err != nil {
		dieErr(err)
	}

	if *o.Logging {
		opts.Logging = true
	}

	client, err = e3db.GetClient(*opts)
	if err != nil {
		dieErr(err)
	}

	return client
}

var options cliOptions

func cmdList(cmd *cli.Cmd) {
	data := cmd.BoolOpt("d data", false, "include data in JSON format")
	outputJSON := cmd.BoolOpt("j json", false, "output in JSON format")
	contentTypes := cmd.StringsOpt("t type", nil, "record content types")
	recordIDs := cmd.StringsOpt("r record", nil, "record IDs")
	writerIDs := cmd.StringsOpt("w writer", nil, "record writer IDs or email addresses")
	userIDs := cmd.StringsOpt("u user", nil, "record user IDs")

	cmd.Action = func() {
		client := options.getClient()
		ctx := context.Background()

		// Convert e-mail addresses in write list to writer IDs.
		for ix, writerID := range *writerIDs {
			if strings.Contains(writerID, "@") {
				info, err := client.GetClientInfo(ctx, writerID)
				if err != nil {
					dieErr(err)
				}

				(*writerIDs)[ix] = info.ClientID
			}
		}

		cursor := client.Query(context.Background(), e3db.Q{
			ContentTypes: *contentTypes,
			RecordIDs:    *recordIDs,
			WriterIDs:    *writerIDs,
			UserIDs:      *userIDs,
			IncludeData:  *data,
		})

		first := true
		for cursor.Next() {
			record, err := cursor.Get()
			if err != nil {
				dieErr(err)
			}

			if *outputJSON {
				if first {
					first = false
					fmt.Println("[")
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

func cmdWrite(cmd *cli.Cmd) {
	recordType := cmd.String(cli.StringArg{
		Name:      "TYPE",
		Desc:      "type of record to write",
		Value:     "",
		HideValue: true,
	})

	data := cmd.String(cli.StringArg{
		Name:      "DATA",
		Desc:      "json formatted record data",
		Value:     "",
		HideValue: true,
	})

	cmd.Action = func() {
		client := options.getClient()
		record := client.NewRecord(*recordType)

		err := json.NewDecoder(strings.NewReader(*data)).Decode(&record.Data)
		if err != nil {
			dieErr(err)
		}

		id, err := client.Write(context.Background(), record)
		if err != nil {
			dieErr(err)
		}

		fmt.Println(id)
	}
}

func cmdRead(cmd *cli.Cmd) {
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
			record, err := client.Read(context.Background(), recordID)
			if err != nil {
				dieErr(err)
			}

			bytes, err := json.MarshalIndent(record, "", "  ")
			if err != nil {
				dieErr(err)
			}

			fmt.Println(string(bytes))
		}
	}
}

func cmdDelete(cmd *cli.Cmd) {
	recordIDs := cmd.Strings(cli.StringsArg{
		Name:      "RECORD_ID",
		Desc:      "record IDs to delete",
		Value:     nil,
		HideValue: true,
	})

	cmd.Spec = "RECORD_ID..."
	cmd.Action = func() {
		client := options.getClient()

		for _, recordID := range *recordIDs {
			err := client.Delete(context.Background(), recordID)
			if err != nil {
				dieErr(err)
			}
		}
	}
}

func cmdInfo(cmd *cli.Cmd) {
	clientID := cmd.String(cli.StringArg{
		Name:      "CLIENT_ID",
		Desc:      "client unique id or email",
		Value:     "",
		HideValue: true,
	})

	cmd.Spec = "[CLIENT_ID]"

	cmd.Action = func() {
		client := options.getClient()
		if *clientID == "" {
			fmt.Printf("Client ID:   %s\n", client.ClientID)
			fmt.Printf("Public Key:  %s\n", base64.RawURLEncoding.EncodeToString(client.PublicKey[:]))
			fmt.Printf("API Key ID:  %s\n", client.APIKeyID)
			fmt.Printf("API Secret:  %s\n", client.APISecret)
		} else {
			info, err := client.GetClientInfo(context.Background(), *clientID)
			if err != nil {
				dieErr(err)
			}

			fmt.Printf("Client ID:   %s\n", info.ClientID)
			fmt.Printf("Public Key:  %s\n", info.PublicKey.Curve25519)
		}
	}
}

func cmdShare(cmd *cli.Cmd) {
	recordType := cmd.String(cli.StringArg{
		Name:      "TYPE",
		Desc:      "type of records to share",
		Value:     "",
		HideValue: true,
	})

	clientID := cmd.String(cli.StringArg{
		Name:      "CLIENT_ID",
		Desc:      "client unique id or email",
		Value:     "",
		HideValue: true,
	})

	cmd.Action = func() {
		client := options.getClient()

		err := client.Share(context.Background(), *recordType, *clientID)
		if err != nil {
			dieErr(err)
		}

		fmt.Printf("Records of type '%s' are now shared with client '%s'\n", *recordType, *clientID)
	}
}

func cmdUnshare(cmd *cli.Cmd) {
	recordType := cmd.String(cli.StringArg{
		Name:      "TYPE",
		Desc:      "type of records to share",
		Value:     "",
		HideValue: true,
	})

	clientID := cmd.String(cli.StringArg{
		Name:      "CLIENT_ID",
		Desc:      "client unique id or email",
		Value:     "",
		HideValue: true,
	})

	cmd.Action = func() {
		client := options.getClient()

		err := client.Unshare(context.Background(), *recordType, *clientID)
		if err != nil {
			dieErr(err)
		}

		fmt.Printf("Records of type '%s' are no longer shared with client '%s'\n", *recordType, *clientID)
	}
}

func cmdRegister(cmd *cli.Cmd) {
	apiBaseURL := cmd.String(cli.StringOpt{
		Name:      "api",
		Desc:      "e3db api base url",
		Value:     "",
		HideValue: true,
	})

	authBaseURL := cmd.String(cli.StringOpt{
		Name:      "auth",
		Desc:      "e3db auth service base url",
		Value:     "",
		HideValue: true,
	})

	isPublic := cmd.Bool(cli.BoolOpt{
		Name:      "public",
		Desc:      "allow other clients to find you by email",
		Value:     false,
		HideValue: false,
	})

	email := cmd.String(cli.StringArg{
		Name:      "EMAIL",
		Desc:      "client e-mail address",
		Value:     "",
		HideValue: true,
	})

	// TODO: minimally validate that email looks like an email address

	cmd.Action = func() {
		// Preflight check for existing configuration file to prevent a later
		// failure writing the file (since we use O_EXCL) after registration.
		if e3db.ProfileExists(*options.Profile) {
			var name string
			if *options.Profile != "" {
				name = *options.Profile
			} else {
				name = "(default)"
			}

			dieFmt("register: profile %s already registered", name)
		}

		info, err := e3db.RegisterClient(*email, e3db.RegistrationOpts{
			APIBaseURL:  *apiBaseURL,
			AuthBaseURL: *authBaseURL,
			Logging:     *options.Logging,
			FindByEmail: *isPublic,
		})

		if err != nil {
			dieErr(err)
		}

		err = e3db.SaveConfig(*options.Profile, info)
		if err != nil {
			dieErr(err)
		}
	}
}

func main() {
	app := cli.App("e3db-cli", "E3DB Command Line Interface")

	app.Version("v version", "e3db-cli 0.0.1")

	options.Logging = app.BoolOpt("d debug", false, "enable debug logging")
	options.Profile = app.StringOpt("p profile", "", "e3db configuration profile")

	app.Command("register", "register a client", cmdRegister)
	app.Command("info", "get client information", cmdInfo)
	app.Command("ls", "list records", cmdList)
	app.Command("read", "read records", cmdRead)
	app.Command("write", "write a record", cmdWrite)
	app.Command("delete", "delete a record", cmdDelete)
	app.Command("share", "share records with another client", cmdShare)
	app.Command("unshare", "stop sharing records with another client", cmdUnshare)
	app.Run(os.Args)
}
