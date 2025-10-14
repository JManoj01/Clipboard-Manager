# 🖋️ Terminal Clipboard Manager

A **cross-platform clipboard manager** built in Go with a beautiful **Bubble Tea TUI**, stylish **Lipgloss design**, and **syntax highlighting** via Chroma. Supports **Linux, macOS, and Windows**.

---

## ✨ Features

### 📋 Browse & Search
Browse and search **clipboard history** directly from your terminal.

### 🔍 Fuzzy Search
Fast lookup with **fuzzy search**.

### 🧠 Auto-categorization
Automatically categorize by content type (**code, text, images, etc.**).

### 💾 Persistent Storage
Uses **JSON** to store clipboard history locally.

### 🎨 Syntax Highlighting
**Chroma** powers syntax highlighting for code snippets.

### ⚡ High Capacity
Handles **1,000+ entries** with **duplicate detection**.

### 📤 Export Functionality
Export **text and image history** easily.

---

## 🧰 Requirements

- **Go 1.21+**
- Clipboard access enabled on your system:
  - Linux: `xclip` / `xsel`  
  - macOS: built-in `pbcopy` / `pbpaste`

---

## 🚀 Installation & Usage

### 1️⃣ Build the Executable
```bash
go build -o clipboard_manager


# 2️⃣ Run the Program
./clipboard_manager

⚙️ Configuration

Clipboard history and configuration are stored in a backup JSON file:

~/.config/clipboard_manager/history.json


You can edit or back up this file as needed.
