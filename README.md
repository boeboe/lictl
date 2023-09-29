# LinkedIn CLI Tool


```bash
lictl/
│
├── cmd/                  # Command-line related code
│   ├── root.go
│   └── search.go
│
├── pkg/                  # Reusable packages (your LinkedIn interaction logic)
│   ├── linkedin/         # LinkedIn specific logic
│   │   ├── auth.go       # Authentication related functions
│   │   ├── job.go        # Job search related functions
│   │   └── ...           # Other LinkedIn functionalities
│   │
│   └── utils/            # Any utility functions
│
├── api/                  # For future REST API related code
│
├── internal/             # Private application and library code
│
├── scripts/              # Scripts to perform various build, install, analysis, etc operations
│
├── testdata/             # Test data used in testing (if any)
│
├── .gitignore
├── go.mod                # Go module file
├── go.sum                # Go module checksum
└── README.md
```

Details:

- `cmd/`: This directory contains application-specific code for your CLI. Each sub-command can have its own file.
- `pkg/`: This is where you'll place the reusable code that interacts with LinkedIn. By placing this code here, both your CLI and future REST API can use it without duplication. The linkedin package inside pkg will contain all the logic for interacting with LinkedIn.
- `api/`: In the future, when you implement the REST API, you can place the API-specific code here.
- `internal/`: Code you don't want to expose to other applications or libraries. It's a Go convention to prevent importing.
- `scripts/`: Any build or utility scripts.
- `testdata/`: If you have any data that you use for testing, it can be placed here.
