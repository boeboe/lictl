# LinkedIn CLI Tool

## Table of Contents
- [Introduction](#introduction)
- [Usage](#usage)
  - [Global Flags](#global-flags)
  - [Commands](#commands)
    - [job](#job)
      - [search](#search)
- [Release/Download Information](#releasedownload-information)
- [Development](#development)
  - [Repository structure](#repository-structure)
- [License](#license)
- [Contributing](#contributing)

## Introduction
`lictl` is a command-line tool designed to interact with LinkedIn functionalities. It provides a range of features that allow users to search for jobs, manage their LinkedIn profiles, and more, all from the comfort of their terminal.

## Usage

### Global Flags:
- `--help` or `-h`: Shows help for the command.

### Commands:

#### job
- **Usage**: `lictl job`
- **Description**: Interact with LinkedIn job functionalities.

##### search
- **Usage**: `lictl job search`
- **Description**: Search for LinkedIn jobs based on regions and keywords.
- **Flags**:
  - `--regions` or `-r`: Specify one or more regions. (Mandatory)
  - `--keywords` or `-k`: Specify one or more keywords. (Mandatory)
  - `--output` or `-o`: Specify the output directory. Default is the current working directory.
  - `--format` or `-f`: Specify the format (json/csv). Default is `json`.
  - `--debug` or `-d`: Enable or disable debug mode. Default is `false`.
  - `--interval` or `-i`: Specify the interval between web calls. Default is `100ms`.

**Example Usages**:

```bash
lictl job search --regions "New York" --keywords "Software Engineer"
lictl job search -r "San Francisco" -k "Data Scientist" -o "./results" -f "csv"
```

## Release/Download Information

For the latest releases and download options, please visit the [releases section](https://github.com/boeboe/lictl/releases) of the `lictl` [GitHub repository](https://github.com/boeboe/lictl/releases).


## Development

To build the lictl program, you can use the provided Makefile targets:

```bash
$ make
help                           This help
lint                           Run linter on all source code
test                           Run all tests recursively
build                          Build the project
release                        Create a GitHub release and upload the binary
```

### Repository structure



```bash
lictl/
├── bin                          # Created when building locally
│   ├── lictl
│   ├── lictl-arm64
│   ├── lictl-windows-amd64.exe
│   └── lictl-x86_64
├── cmd/                         # Command-line related code
├── pkg/                         # Reusable packages (your LinkedIn interaction logic)
└── testdata/                    # Test data used in testing (if any)
```

Details:

- `cmd/`: This directory contains application-specific code for your CLI. Each sub-command can have its own file.
- `pkg/`: This is where you'll place the reusable code that interacts with LinkedIn. By placing this code here, both your CLI and future REST API can use it without duplication. The linkedin package inside pkg will contain all the logic for interacting with LinkedIn.
- `api/`: In the future, when you implement the REST API, you can place the API-specific code here.
- `internal/`: Code you don't want to expose to other applications or libraries. It's a Go convention to prevent importing.
- `scripts/`: Any build or utility scripts.
- `testdata/`: If you have any data that you use for testing, it can be placed here.


## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

We welcome contributions from the community! If you'd like to contribute, please fork the repository, make your changes, and submit a pull request. Ensure that you've tested your changes before submitting the pull request.

