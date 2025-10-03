# ⏱️ Timer

A minimal, blazing-fast TUI countdown timer and stopwatch for the terminal. Built for efficiency with **<5MB memory usage**.

## ✨ Features

- ⏱️ **Countdown Timer** - Set durations with intuitive syntax (`5s`, `2m`, `1h`)
- ⏲️ **Stopwatch Mode** - Count up from 00:00 when no duration is specified
- 🖥️ **Fullscreen TUI** - Large ASCII art display with centered output
- 📟 **Inline Mode** - Compact display option for command-line use
- ⏸️ **Pause/Resume** - Pause and resume timers with spacebar
- 🎨 **Color Indicators** - Visual feedback (red warning <5min, blue when paused)
- ⚡ **Low Resource Usage** - Optimized adaptive ticker intervals
- ⌨️ **Simple Controls** - Intuitive keyboard shortcuts

## 🚀 Installation

### Using Go Install

```bash
# Install latest stable version
go install github.com/Zihad550/timer@v0.1.0

# Or install latest (may show development version)
go install github.com/Zihad550/timer@latest
```

### Build from Source

```bash
git clone https://github.com/Zihad550/timer.git
cd timer
go build -o timer .
```

## 📖 Usage

### Basic Examples

```bash
# Stopwatch mode (counts up from 00:00)
timer

# Countdown timer - 5 seconds
timer 5

# Countdown timer - 2 minutes
timer 2m

# Countdown timer - 1 hour
timer 1h

# Inline mode (no fullscreen)
timer -i 30s

# Display version
timer -version
```

### Duration Format

- **Numbers only**: Interpreted as seconds (e.g., `timer 60` = 60 seconds)
- **With units**: `s` (seconds), `m` (minutes), `h` (hours)
- **Examples**: `5s`, `90s`, `2m`, `1h30m`

### Command-Line Options

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--inline` | `-i` | Run in inline mode (disable fullscreen TUI) |
| `--version` | `-v` | Display version information |

## ⌨️ Keyboard Controls

| Key | Action |
|-----|--------|
| <kbd>Space</kbd> | Pause/Resume timer |
| <kbd>q</kbd> / <kbd>Q</kbd> / <kbd>ESC</kbd> | Quit |
| <kbd>Ctrl</kbd>+<kbd>C</kbd> | Force quit |

## 🎨 Visual Indicators

- **Default** - Normal white/terminal color
- **🔴 Red** - Countdown timer with <5 minutes remaining
- **🔵 Blue** - Timer is paused

## ⚙️ Technical Details

### Architecture

- **Language**: Go 1.24.0+
- **Dependencies**: `golang.org/x/term`
- **Memory**: <5MB footprint
- **Performance**: Adaptive ticker intervals based on duration
  - Fast (100ms) - Durations <1 minute
  - Medium (500ms) - Durations 1-10 minutes
  - Slow (1s) - Durations >10 minutes

### Project Structure

```
timer/
├── main.go         # CLI entry point and argument parsing
├── timer.go        # Core timer logic and event loop
├── display.go      # Text formatting and rendering
├── terminal.go     # Terminal control and raw mode
├── config.go       # Configuration constants
├── glyphs.go       # ASCII art character definitions
└── utils.go        # Helper functions
```

## 🛠️ Development

### Requirements

- Go 1.24.0 or higher
- Unix-like terminal (Linux, macOS, WSL)

### Building

```bash
go build -o timer .
```

### Running Tests

```bash
go test ./...
```

## 📝 Examples

### Pomodoro Timer (25 minutes)
```bash
timer 25m
```

### Quick Break (5 minutes)
```bash
timer 5m
```

### Meeting Timer (1 hour)
```bash
timer 1h
```

### Track Work Session (Stopwatch)
```bash
timer
```

### Inline Timer for Scripts
```bash
timer -i 10s && echo "Task complete!"
```

## 🤝 Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is open source. Check the LICENSE file for details.

## 💡 Tips

- Use **stopwatch mode** (no arguments) for tracking open-ended tasks
- Combine with shell commands: `timer 25m && notify-send "Break time!"`
- Perfect for Pomodoro technique, cooking, workouts, and more
- Runs entirely in the terminal—no GUI required

---

Made with ⚡ by developers who value simplicity and performance.
