package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

// ANSI escape codes for colors
const (
	reset       = "\033[0m"
	red         = "\033[31m"
	green       = "\033[32m"
	yellow      = "\033[33m"
	blue        = "\033[34m"
	magenta     = "\033[35m"
	white       = "\033[97m"
	cyan        = "\033[36m"
	orange      = "\033[38;5;208m"
	purple      = "\033[35m"
	lightRed    = "\033[91m"
	lightPurple = "\033[95m"
	darkGreen   = "\033[38;5;22m"
	darkOrange  = "\033[38;5;208m"
	darkYellow  = "\033[38;5;172m"
	darkMagenta = "\033[38;5;125m"
)

var (
	longListing   bool
	humanReadable bool
	fileSize      bool

	// File icons based on extensions
	fileIcons = map[string]string{
		".go":   cyan + " " + reset,
		".sh":   white + " " + reset,
		".cpp":  blue + " " + reset,
		".hpp":  blue + " " + reset,
		".cxx":  blue + " " + reset,
		".hxx":  blue + " " + reset,
		".css":  blue + " " + reset,
		".c":    blue + " " + reset,
		".png":  magenta + " " + reset,
		".jpg":  magenta + " " + reset,
		".jpeg": magenta + " " + reset,
		".webp": magenta + " " + reset,
		".xcf":  white + " " + reset,
		".xml":  red + " " + reset,
		".htm":  red + " " + reset,
		".html": red + " " + reset,
		".txt":  white + " " + reset,
		".mp3":  cyan + " " + reset,
		".ogg":  cyan + " " + reset,
		".mp4":  cyan + " " + reset,
		".zip":  yellow + "󰿺 " + reset,
		".tar":  yellow + "󰿺 " + reset,
		".gz":   yellow + "󰿺 " + reset,
		".bz2":  yellow + "󰿺 " + reset,
		".xz":   yellow + "󰿺 " + reset,
		".jar":  white + " " + reset,
		".java": white + " " + reset,
		".js":   yellow + " " + reset,
		".py":   yellow + " " + reset,
		".rs":   orange + " " + reset,
		".deb":  red + " " + reset,
		".md":   blue + " " + reset,
		".rb":   red + " " + reset,
		".php":  purple + " " + reset,
		".pl":   orange + " " + reset,
		".svg":  magenta + " " + reset,
		".eps":  magenta + " " + reset,
		".ps":   magenta + " " + reset,
		".git":  orange + " " + reset,
		".zig":  darkOrange + " " + reset,
		".xbps": darkGreen + " " + reset,
	}
)

func main() {
	parseFlags()

	directory := "."
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[len(os.Args)-1], "-") {
		directory = os.Args[len(os.Args)-1]
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		fmt.Println("No files found.")
		return
	}

	if longListing {
		printLongListing(files, directory)
	} else if fileSize {
		getFileSize(files, directory)
	} else {
		printFilesInColumns(files, directory)
	}
}

func parseFlags() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-l":
			longListing = true
		case "-h":
			showHelp()
			os.Exit(0)
		case "-lh", "-hl":
			longListing = true
			humanReadable = true
		case "-s":
			fileSize = true
		case "-hs", "-sh":
			fileSize = true
		default:
			if !strings.HasPrefix(arg, "-") {
				continue
			}
			showHelp()
			os.Exit(1)
		}
	}
}

func showHelp() {
	fmt.Println("Usage: gols [options] [directory]")
	fmt.Println("Options:")
	fmt.Println("  -l    Long listing format")
	fmt.Println("  -lh   Human-readable file sizes")
	fmt.Println("  -s    print files size")
	fmt.Println("  -h    Show options")
}

func printFilesInColumns(files []os.DirEntry, directory string) {
	maxFilesInLine := 3
	maxFileNameLength := 19

	filesInLine := 0
	for _, file := range files {
		printFile(file, directory)
		filesInLine++
		if filesInLine >= maxFilesInLine || len(file.Name()) > maxFileNameLength {
			fmt.Println()
			filesInLine = 0
		} else {
			printPadding(file.Name(), maxFileNameLength)
		}
	}
	fmt.Println() // Ensure a newline at the end
}

func getFileSize(files []os.DirEntry, directory string) {
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		size := info.Size()
		sizeStr := fmt.Sprintf("%d", size)
		if humanReadable {
			sizeStr = humanizeSize(size)
		}
		var spaces = 10 - len(sizeStr)
		fmt.Print(sizeStr)
		for i := 0; i < spaces; i++ {
			fmt.Print(" ")
		}
		fmt.Println(getFileIcon(file.Name()) + file.Name())
	}
	fmt.Println()
}

func printLongListing(files []os.DirEntry, directory string) {
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		permissions := formatPermissions(info.Mode())
		size := info.Size()
		sizeStr := fmt.Sprintf("%d", size)
		if humanReadable {
			sizeStr = humanizeSize(size)
		}

		// Get owner and group names
		owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
		if err != nil {
			log.Fatal(err)
		}
		group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Gid))
		if err != nil {
			log.Fatal(err)
		}

		// Print long listing format with icons
		fmt.Printf("%s %10s %s %s", permissions, sizeStr, owner.Username, group.Name)
		fmt.Printf(" %s", info.ModTime().Format("Jan 02 15:04"))

		fmt.Printf(" %s %s\n", getFileIcon(file.Name()), file.Name())
	}
}

func formatPermissions(mode os.FileMode) string {
	var b strings.Builder

	if mode.IsDir() {
		b.WriteString("d")
	} else {
		b.WriteString("-")
	}

	b.WriteString(rwx(mode.Perm() >> 6)) // Owner permissions
	b.WriteString(rwx(mode.Perm() >> 3)) // Group permissions
	b.WriteString(rwx(mode.Perm()))      // Other permissions

	return b.String()
}

func rwx(perm os.FileMode) string {
	var b strings.Builder

	if perm&0400 != 0 {
		b.WriteString(green + "r")
	} else {
		b.WriteString("-")
	}
	if perm&0200 != 0 {
		b.WriteString(yellow + "w")
	} else {
		b.WriteString("-")
	}
	if perm&0100 != 0 {
		b.WriteString(red + "x")
	} else {
		b.WriteString("-")
	}

	b.WriteString(reset) // Reset colors

	return b.String()
}

func printFile(file os.DirEntry, directory string) {
	name := file.Name()
	icon := getFileIcon(name)

	if file.IsDir() {
		fmt.Print(blue + icon + name + "/" + reset)
	} else {
		fmt.Print(icon + name)
	}
}

func getFileIcon(name string) string {
	ext := filepath.Ext(name)
	icon, exists := fileIcons[ext]
	if exists {
		return icon
	}
	return white + " " + reset // Default icon
}

func printPadding(fileName string, maxFileNameLength int) {
	padding := maxFileNameLength - len(fileName)
	for i := 0; i < padding; i++ {
		fmt.Print(" ")
	}
}

func humanizeSize(size int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	switch {
	case size >= PB:
		return fmt.Sprintf("%.2fPB", float64(size)/PB)
	case size >= TB:
		return fmt.Sprintf("%.2fTB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2fKB", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}
