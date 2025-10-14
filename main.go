package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clipboard_manager/clipboard"
	"clipboard_manager/storage"
	"clipboard_manager/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := clipboard.Init(); err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	db, err := storage.NewDatabase("clipboard_history.json")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	os.MkdirAll("clipboard_images", 0755)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	p := ui.NewProgram(db)

	go startEnhancedWatcher(ctx, db, p)

	if _, err := p.Run(); err != nil {
		log.Printf("UI error: %v", err)
	}

	cancel()
}

func startEnhancedWatcher(ctx context.Context, db *storage.Database, p *tea.Program) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastText string
	var lastImageHash string

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if clipboard.HasImage() {
				imageHash, _ := clipboard.GetImageAsBase64()
				if imageHash != lastImageHash && imageHash != "" {
					lastImageHash = imageHash
					filename := fmt.Sprintf("clipboard_images/img_%d.png", time.Now().Unix())
					if err := clipboard.SaveImageToFile(filename); err == nil {
						db.AddImageEntry(filename)
						p.Send(ui.StatusMsg("🖼️ Saved image: " + filename))
					}
				}
			}

			text, err := clipboard.ReadText()
			if err != nil {
				continue
			}

			if text != lastText && text != "" {
				lastText = text
				db.AddEntry(text)
				preview := text
				if len(preview) > 60 {
					preview = preview[:60] + "..."
				}
				p.Send(ui.StatusMsg("✓ Saved: " + preview))
			}
		}
	}
}
