package ui

import (
	"bufio"
	"clipboard_manager/search"
	"clipboard_manager/storage"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Terminal struct {
	db *storage.Database
}

func NewTerminal(db *storage.Database) *Terminal {
	return &Terminal{db: db}
}

func (t *Terminal) Run(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	time.Sleep(1 * time.Second)
	t.printHelp()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("\n📋 > ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "" {
				continue
			}

			t.handleCommand(input)
		}
	}
}

func (t *Terminal) handleCommand(input string) {
	parts := strings.SplitN(input, " ", 2)
	command := strings.ToLower(parts[0])

	switch command {
	case "list", "l":
		limit := 10
		if len(parts) > 1 {
			if n, err := strconv.Atoi(parts[1]); err == nil {
				limit = n
			}
		}
		t.listEntries(limit)

	case "search", "s":
		if len(parts) < 2 {
			fmt.Println("❌ Usage: search <query>")
			return
		}
		t.searchEntries(parts[1])

	case "fuzzy", "f":
		if len(parts) < 2 {
			fmt.Println("❌ Usage: fuzzy <query>")
			return
		}
		t.fuzzySearch(parts[1])

	case "view", "v":
		if len(parts) < 2 {
			fmt.Println("❌ Usage: view <id>")
			return
		}
		if id, err := strconv.Atoi(parts[1]); err == nil {
			t.viewEntry(id)
		}

	case "delete", "d":
		if len(parts) < 2 {
			fmt.Println("❌ Usage: delete <id>")
			return
		}
		if id, err := strconv.Atoi(parts[1]); err == nil {
			t.deleteEntry(id)
		}

	case "clear":
		t.clearHistory()

	case "help", "h":
		t.printHelp()

	case "quit", "q", "exit":
		fmt.Println("👋 Goodbye!")
		os.Exit(0)

	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
	}
	
}

func (t *Terminal) listEntries(limit int) {
	entries, err := t.db.GetRecent(limit)
	if err != nil {	
		fmt.Println(errText(fmt.Sprintf("Error: %v", err)))
		return
	}

	if len(entries) == 0 {
		fmt.Println(info("No clipboard history yet!"))
		return
	}

	fmt.Println("\n" + colorize(ColorCyan, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Printf("%s Recent %d entries:\n", colorize(ColorYellow, "📝"), len(entries))
	fmt.Println(colorize(ColorCyan, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))

	for _, entry := range entries {
		preview := t.formatPreview(entry.Text, 100, true)
		timeAgo := t.formatTimeAgo(entry.Timestamp)
		
		idStr := colorize(ColorBlue, fmt.Sprintf("[%d]", entry.ID))
		timeStr := colorize(ColorDim, fmt.Sprintf("⏰ %s", timeAgo))
		
		fmt.Printf("%s %s\n    %s\n", idStr, preview, timeStr)
	}
	
	fmt.Println("\n" + info("Tip: Use 'view <id>' to see full formatted content"))
}

func (t *Terminal) viewEntry(id int) {
	entries, err := t.db.GetRecent(1000)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	for _, entry := range entries {
		if entry.ID == id {
			fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("📄 Entry #%d (Copied %s)\n", entry.ID, t.formatTimeAgo(entry.Timestamp))
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			
			t.displayFormatted(entry.Text)
			
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			return
		}
	}
	
	fmt.Printf("❌ Entry #%d not found\n", id)
}

func (t *Terminal) searchEntries(query string) {
	entries, err := t.db.Search(query)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Printf("🔍 No results for: %s\n", query)
		return
	}

	fmt.Printf("\n🔍 Found %d results:\n", len(entries))
	for _, entry := range entries {
		preview := t.formatPreview(entry.Text, 100, true)
		fmt.Printf("[%d] %s\n", entry.ID, preview)
	}
}

func (t *Terminal) fuzzySearch(query string) {
	allEntries, _ := t.db.GetRecent(100)
	results := search.FuzzySearch(allEntries, query, 2)

	if len(results) == 0 {
		fmt.Printf("🔮 No fuzzy matches for: %s\n", query)
		return
	}

	fmt.Printf("\n🔮 Fuzzy search: %d results\n", len(results))
	for _, entry := range results {
		preview := t.formatPreview(entry.Text, 100, true)
		fmt.Printf("[%d] %s\n", entry.ID, preview)
	}
}

func (t *Terminal) deleteEntry(id int) {
	t.db.DeleteEntry(id)
	fmt.Printf("✅ Deleted #%d\n", id)
}

func (t *Terminal) clearHistory() {
	fmt.Print("⚠️  Clear all? (yes/no): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) == "yes" {
		t.db.Clear()
		fmt.Println("✅ Cleared!")
	}
}

func (t *Terminal) printHelp() {
	fmt.Println("\n" + colorize(ColorCyan, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(bold(colorize(ColorYellow, "📚 Available Commands:")))
	fmt.Println(colorize(ColorCyan, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Printf("  %s - Show last n entries (compact)\n", colorize(ColorGreen, "list [n]"))
	fmt.Printf("  %s - View full entry with formatting\n", colorize(ColorGreen, "view <id>"))
	fmt.Printf("  %s - Search clipboard\n", colorize(ColorGreen, "search <text>"))
	fmt.Printf("  %s - Fuzzy search\n", colorize(ColorGreen, "fuzzy <text>"))
	fmt.Printf("  %s - Add tags to entry\n", colorize(ColorGreen, "tag <id> <tags>"))
	fmt.Printf("  %s - Show statistics\n", colorize(ColorGreen, "stats"))
	fmt.Printf("  %s - Export to file\n", colorize(ColorGreen, "export <file>"))
	fmt.Printf("  %s - Delete entry\n", colorize(ColorRed, "delete <id>"))
	fmt.Printf("  %s - Clear all\n", colorize(ColorRed, "clear"))
	fmt.Printf("  %s - Show this help\n", colorize(ColorBlue, "help"))
	fmt.Printf("  %s - Exit program\n", colorize(ColorBlue, "quit"))
	fmt.Println(colorize(ColorCyan, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
}

func (t *Terminal) formatPreview(text string, maxLen int, compact bool) string {
	if compact {
		lines := strings.Split(text, "\n")
		firstLine := strings.TrimSpace(lines[0])
		
		if len(lines) > 1 {
			// Show it's multi-line
			if len(firstLine) > maxLen-10 {
				return firstLine[:maxLen-10] + "... [+" + fmt.Sprintf("%d", len(lines)-1) + " lines]"
			}
			return firstLine + " [+" + fmt.Sprintf("%d", len(lines)-1) + " lines]"
		}
		
		if len(firstLine) > maxLen {
			return firstLine[:maxLen] + "..."
		}
		return firstLine
	}
	
	// Full preview
	if len(text) > maxLen {
		return text[:maxLen] + "..."
	}
	return text
}

func (t *Terminal) displayFormatted(text string) {
	lines := strings.Split(text, "\n")
	
	for i, line := range lines {
		fmt.Printf("  %s\n", line)
		
		if i > 200 {
			fmt.Printf("  ... [%d more lines truncated]\n", len(lines)-i-1)
			break
		}
	}
}

func (t *Terminal) formatTimeAgo(timestamp time.Time) string {
	duration := time.Since(timestamp)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		return fmt.Sprintf("%d min ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}



