# Pitwall 🏎️

Pitwall is a Go-based command-line interface (CLI) tool designed to provide Formula 1 session data at your fingertips. By integrating with the **OpenF1 API**, Pitwall tracks race weekends and provides intelligent status updates for every session.

## 🌟 Features

- **Subcommand Architecture**: Clean CLI interface using `get_session` and `latest` commands.
- **Smart Session States**: Automatically classifies sessions into **Future**, **Live**, or **Finished** based on real-time data.
- **Human-Readable Timing**: 
  - Displays countdowns for upcoming races (e.g., `Starts in: 2d 4h 30m`).
  - Tracks live progress (e.g., `Ends in: 18m`).
  - Shows historical context (e.g., `Ended: 1h 42m ago`).
- **Timezone Awareness**: Handles UTC to UK (Europe/London) conversion automatically.
- **Robust Testing**: Table-driven tests for date parsing, domain logic, and even terminal output.

## 📂 Project Structure

```text
.
├── .pitwall_cache        # Cache file for storing session data
├── cmd/pitwall/          # Main entry point & CLI routing
├── domain/               # Domain entities (Session) and Business Logic
├── internal/
│   ├── openf1/           # API Client for OpenF1
│   └── services/
│       ├── getsession/   # GetSession service & UI formatting
│       └── latest/       # Logic for fetching the current/next session
│       └── weekend/      # Logic for fetching weekend sessions
│       └── cache/        # Logic for cache management
        └── remind/       # Reminder service for upcoming sessions
└── utils/                # Shared utilities (Date parsing, formatting)
```

## 🚀 Getting Started

### Prerequisites
- [Go](https://go.dev/doc/install) 1.21+
  
### Installation

```bash
git clone https://github.com/bhopalg/pitwall.git
cd pitwall
go build -o pitwall ./cmd/pitwall
```

### Usage

#### Check the most recent or active session:

```bash
./pitwall latest
```

#### Find a specific session with flags:

```bash
./pitwall get_session --country Belgium --type Race --year 2023
```

#### Find a specific weekend:
```bash
./pitwall weekend --country Belgium --year 2023
```

#### Clear the cache:
```bash
./pitwall cache clear
```

#### Information about the cache:
```bash
./pitwall cache info
```

|### Reminder for upcoming sessions:
```bash
./pitwall remind
```

## 🛠️ Development

### Running Tests

The project prioritizes testability by injecting time.Now into services. To run the full suite:

```bash
go test ./... -v
```

### VS Code Debugging
A `launch.json` is included to support debugging with arguments. You can test different flags by editing the `args` array in your debug configuration.

## 📊 Logic: Session State Rules
| State      | Condition                                         | Display               |
| ---------- | ------------------------------------------------- | --------------------- |
| **Future** | `now < start`                                     | Starts at [Time] (UK) |
| **Live**   | `now >= start` AND (`now < end` OR `end` is null) | Ends in [Duration]    |
| **Future** | `now >= end`                                      | Ended [Duration] ago  |