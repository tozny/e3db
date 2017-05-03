# Tozny End-to-End Encrypted Database

The Tozny End-to-End Encrypted Database (E3DB) is a storage platform
with powerful sharing and consent management features. [Read more on
our
blog](https://tozny.com/blog/announcing-project-e3db-the-end-to-end-encrypted-database/).

Tozny's E3DB provides a familiar JSON-based NoSQL-style API for
reading, writing, and listing JSON data stored securely in the cloud.

## Quick Start

Please try out E3DB and give us feedback! Here are the basic steps.
E3DB has been tested on MacOS, Windows, and Linux:

 1. Download the appropriate binary from our [releases page](https://github.com/tozny/e3db-cli/releases) and save it somewhere in your PATH. For convenience, rename the binary to `e3db`.
 1. `$ e3db register` - then check your email!
 1. `$ e3db ls` - You should see nothing.
 1. Write a record: `$ recordID=$(e3db write address_book '{"name": "John Doe", "phone": "503-555-1212"}')`.
 1. `$ e3db ls` - You should see your new record.
 1. Read a record: `$ e3db read $recordID`.

## Terms of Service

Your use of E3DB must abide by our [Terms of Service](terms.pdf), as detailed in
the linked document.

# Installation & Use

To obtain the E3DB CLI binary, download the
[1.0.0 binary](https://github.com/tozny/e3db-cli/releases/tag/1.0.0)
for your platform from our [releases
page](https://github.com/tozny/e3db-cli/releases).

After downloading the file, rename it to `e3db` for convenience and
set it executable.

You should now be able to run the E3DB CLI via the `e3db` command:

```
$ e3db --help
Usage: e3db-cli [OPTIONS] COMMAND [arg...]

E3DB Command Line Interface
[...]
```

## Registration

Before you can use E3DB, you must register an account and receive API
credentials:

```
$ e3db register <your email>
```

After registration, API credentials and other configuration will be
written to the file `$HOME/.tozny/e3db.json` and can be displayed by
running `e3db info`.

To use your account, it must be verified. Tozny will send a
verification e-mail to the email address you provided during
registration. Simply click the link in the e-mail to complete
verification.

## Using the CLI

These examples demonstrate how to use the E3DB Command Line
Interface. Note that all E3DB commands have help, so anytime you can see
the documentation for a given command using the `--help` argument. For
example, you can see help on all commands:

```
$ e3db --help
...
```

Or help on a particular command, such as `register`:

```
$ e3db register --help
...
```

### Writing Records

To write a record containing free-form JSON data, use the
`e3db write` subcommand. Each record is tagged with a "content
type", which is a string that you choose used to identify the
structure of your data.

In this example, we write an address book entry into E3DB:

```
$ e3db write address_book '{"name": "John Doe", "phone": "503-555-1212"}'
874b41ff-ac84-4961-a91d-9e0c114d0e92
```

Once E3DB has written the record, it outputs the UUID of the newly
created data. This can be used later to retrieve the specific record.

### Searching & Listing Records

To list all records that we have access to in E3DB, use the
`e3db ls` command:

```
$ e3db ls
768d2ef7-36b4-4061-923c-d38bf72d03d3     message
50176d7a-c026-49bd-be1e-4df1e7c49b1f     message
```

#### Formats & Data

To see data associated with each record when using `ls`, use the `-jd` flags:

```
$ e3db ls -jd
[
  {
    "meta": {
      "record_id": "768d2ef7-36b4-4061-923c-d38bf72d03d3",
      "writer_id": "dac7899f-c474-4386-9ab8-f638dcc50dec",
      "user_id": "dac7899f-c474-4386-9ab8-f638dcc50dec",
      "type": "message",
      "plain": null,
      "created": "2017-05-03T16:38:00.692946Z",
      "last_modified": "2017-05-03T16:38:00.692946Z"
    },
    "data": {
      "content": "Twas brillig, and the slithy toves"
    }
  },
  {
    "meta": {
      "record_id": "50176d7a-c026-49bd-be1e-4df1e7c49b1f",
      "writer_id": "dac7899f-c474-4386-9ab8-f638dcc50dec",
      "user_id": "dac7899f-c474-4386-9ab8-f638dcc50dec",
      "type": "message",
      "plain": null,
      "created": "2017-05-03T16:38:00.692946Z",
      "last_modified": "2017-05-03T16:38:00.692946Z"
    },
    "data": {
      "content": "Did gyre and gimble in the wabe:"
    }
  }
]
```

#### Filters

Several filters are available for matching records. Each argument can take a comma-separated list of values:

- `-t`/`--type` - Retrieve records with the given content type.
- `-r`/`--record` - Retrieve records with the given ID.
- `-w`/`--writer` - Retrieve records written by the given writer. Each writer is identified by their unique ID or email address.
- `-u`/`--user` - Retrieve records written about the given user. Each user is identified by their unique ID.

#### Search & List Examples

Search for records (written by you) about a specific user ID:

```
$ e3db ls -jd --user e04af806-3bbf-41d6-9ed1-a12bdab671ee
....
```

Search for records written by a set of writer IDs. Note you can use email addresses or IDs here:

```
$ e3db ls -jd --writer mimsy@borogoves.com --writer 874b41ff-ac84-4961-a91d-9e0c114d0e92
...
```

### Sharing Records

E3DB allows you to share your data with another E3DB client. In order
to set up sharing, you must know the unique ID or email adddress of
the client you wish to share with. Similarly, if others wish to share
with you, they must know your unique ID or email address. To find the
unique ID of your client, run `e3db info`.

The E3DB client allows you to share records based on their content
type. For example, to share all address book entries with another
client (who registered as `beamish@gimble.com` and made their email
discoverable):

```
$ e3db share address_book beamish@gimble.com
```

This command will set up an access control policy to allow the
client associated with the email address `beamish@gimble.com`
to read your records with type `address_book`. It will also
securely share the encryption key for those records with the
client so they can decrypt the contents of each field.

## SDKs

Tozny provides SDKs for interacting with E3DB. We currently offer SDKs for the following languages:

- [Ruby](http://github.com/tozny/e3db-ruby)
- [Go](http://github.com/tozny/e3db-go)
- [Java](http://github.com/tozny/e3db-client)

Each repository contains information about how to use the SDK,
where to find hosted documentation, and more.