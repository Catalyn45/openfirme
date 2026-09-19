# OpenFirme

Openfirme is an application that provides details about Romanian companies (such as: Financial situations, status, administrators etc ) in the form of a website.

All the source data is retrieved from official government sources available to the public.

The application is currently hosted on [openfirme.ro](https://openfirme.ro).

## Structure

Although the main application is basically a website that queries an [sqlite](https://sqlite.org/) database, the entire project is composed of 3 parts:

- `Web Server` (lets users search for companies)
- `Downloader` (downloads the datasets from government websites)
- `Parser` (parses the downloaded datasets and iserts them into the db)

## Dependencies

Although [golang](https://go.dev/) is the only tool needed to build this project, `CGO` is used in the sqlite library, so you need to have a C compiler in your `PATH`.

If you are on `windows` and have `VisualStudio Build Tools` installed, the easiest workaround is to run the build/run commands in a `Developer Powershell for VS`.

## Building

To build the project packages do the follwing:

1. Enable `CGO` if not already
	```sh
	go env -w CGO_ENABLED=1
	```

1. Initialize the golang packages
	```sh
	go mod download
	```

1. Buid the binaries
	```sh
	go build -tags "fts5" ./cmd/website
	go build -tags "fts5" ./cmd/parser
	go build ./cmd/downloader
	```
	alternatively you can run them directly
	```sh
	go run -tags "fts5" ./cmd/website
	```

	*note: `-tags "fts5"` is needed by the sqlite library in order to enable the `fts5` extension used for fast searching*

## Constructing the database

There are 2 main steps in order to construct the database from scratch:

1. Download the datasets from the government website.

	```sh
	go run ./cmd/downloader
	```

	This should create a `data` folder in your working directory, where all the companies datasets will be downloaded.

	*note: the downloader uses `metadata.json` to keep track of it's downloaded files, so if downloader failed at some point because of some timeout, just run it again and it will take from where it stopped*


1. Parse the datasets and put them into the database
	```sh
	go run -tags "fts5" ./cmd/parser
	```

	This should create a database in your current directory named `foo.db`.

	*note: the database creation process can take a bit of time but keeps track of it's progress, so if you close it at some point just run the parser again and it should take from where it left.*

## Configuration

check `common/config.go` to check all the configurations used in the project.

If can also create a file named `config.json` in your current directory, in which you can specify any runtime config override.
See `config.example.json` for an example.