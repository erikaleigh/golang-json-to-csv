# JSON to CSV Converter

This module takes in multiple `.json` files, and converts them to a single `.csv` file.

Built with [Golang](https://golang.org/).

## Setup and Usage

### Build
`go build cmd/main.go`

### To use
- Remove example .json files from the `files/json` folder, and add the .json files to be converted
- Run the script with `go run cmd/main.go`
- The output `.csv` file will be in the `files/csv` folder
