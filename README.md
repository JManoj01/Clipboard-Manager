🖋️ Terminal Clipboard Manager

A cross-platform clipboard manager built in Go with a beautiful Bubble Tea TUI, stylish Lipgloss design, and syntax highlighting via Chroma.
Supports Linux, macOS, and Windows.

✨ Features

📋 Browse and search clipboard history directly from your terminal

🔍 Fuzzy search for fast lookup

🧠 Auto-categorization by content type (code, text, images, etc.)

💾 Persistent storage using JSON

🎨 Syntax highlighting powered by Chroma

⚡ Handles 1,000+ entries with duplicate detection

📤 Export functionality for text and image history



🧰 Requirements

Go 1.21+

Clipboard access enabled on your system (e.g., xclip/xsel on Linux, built-in pbcopy/pbpaste on macOS)

🚀 Installation & Usage
1. Build the executable
go build -o clipboard_manager

2. Run the program
./clipboard_manager



⚙️ Configuration

Clipboard history and configuration are stored in a backup JSON file:

~/.config/clipboard_manager/history.json


You can edit or back up this file as needed.

🧩 Tech Stack
Component	Description
Go	Core programming language
Bubble Tea	Interactive terminal UI framework
Lipgloss	Stylish terminal UI styling
Chroma	Syntax highlighting for code snippets
JSON	Local data storage format


📜 License

MIT License © 2025 Justin Manoj