package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joennespreuwers/freecam/internal/ui"
	"github.com/joennespreuwers/freecam/internal/watcher"
)

const version = "v1.0.0"

func main() {
	once := flag.Bool("once", false, "kill the target process once and exit")
	process := flag.String("process", "ptpcamera", "name of the process to kill")
	ver := flag.Bool("version", false, "print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "freecam %s — free your PTP camera (Canon, Nikon, Sony, Fujifilm, Olympus…) from macOS ptpcamera\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  freecam               launch TUI daemon mode\n")
		fmt.Fprintf(os.Stderr, "  freecam --once        kill target process once and exit\n")
		fmt.Fprintf(os.Stderr, "  freecam --process foo override target process (default: ptpcamera)\n")
		fmt.Fprintf(os.Stderr, "  freecam --version     print version\n")
		fmt.Fprintf(os.Stderr, "  freecam --help        show this help\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *ver {
		fmt.Println("freecam", version)
		return
	}

	if *once {
		results, err := watcher.FindAndKill(*process)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(results) == 0 {
			fmt.Printf("%s not found — nothing to kill\n", *process)
		} else {
			for _, r := range results {
				fmt.Printf("killed %s (PID %d)\n", r.ProcessName, r.PID)
			}
		}
		return
	}

	// TUI mode
	m := ui.New(version, *process)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
