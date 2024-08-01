package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
    "sort"
	"strconv"
	"strings"
	"syscall"
)

const (
	reset         = "\033[0m"
	black         = "\033[30m"
	red           = "\033[31m"
	green         = "\033[32m"
	yellow        = "\033[33m"
	blue          = "\033[34m"
	magenta       = "\033[35m"
	cyan          = "\033[36m"
	white         = "\033[37m"
	gray          = "\033[90m"
	orange        = "\033[38;5;208m"
	lightRed      = "\033[91m"
	lightGreen    = "\033[92m"
	lightYellow   = "\033[93m"
	lightBlue     = "\033[94m"
	lightMagenta  = "\033[95m"
	lightCyan     = "\033[96m"
	lightWhite    = "\033[97m"
	lightGray     = "\033[37m"
	lightOrange   = "\033[38;5;214m"
	lightPink     = "\033[38;5;218m"
	lightPurple   = "\033[38;5;183m"
	lightBrown    = "\033[38;5;180m"
	lightCyanBlue = "\033[38;5;117m"
	brightOrange  = "\033[38;5;214m"
	brightPink    = "\033[38;5;213m"
	brightCyan    = "\033[38;5;51m"
	brightPurple  = "\033[38;5;135m"
	brightYellow  = "\033[38;5;226m"
	brightGreen   = "\033[38;5;46m"
	brightBlue    = "\033[38;5;33m"
	brightRed     = "\033[38;5;196m"
	brightMagenta = "\033[38;5;198m"
	darkGray      = "\033[38;5;236m"
	darkOrange    = "\033[38;5;208m"
	darkGreen     = "\033[38;5;22m"
	darkCyan      = "\033[38;5;23m"
	darkMagenta   = "\033[38;5;90m"
	darkYellow    = "\033[38;5;172m"
	darkRed       = "\033[38;5;124m"
	darkBlue      = "\033[38;5;18m"

	version       = "gols: 1.2.1"
)

var (
	longListing      bool
	humanReadable    bool
	fileSize         bool
    orderBySize      bool
    orderByTime      bool
    showOnlySymlinks bool
    showHidden       bool
    recursiveListing bool
    dirOnLeft	  	 bool
	oneColumn	     bool
	showSummary		 bool
	showVersion		 bool
	maxDepth		 int = -1

	fileIcons = map[string]string{
		".go":   " ",
        ".mod":  " ",
		".sh":   " ",
		".cpp":  " ",
		".hpp":  " ",
		".cxx":  " ",
		".hxx":  " ",
		".css":  " ",
		".c":    " ",
		".h":    " ",
		".cs":   "󰌛 ",
		".png":  " ",
		".jpg":  "󰈥 ",
		".JPG":  "󰈥 ",
		".jpeg": " ",
		".webp": " ",
		".xcf":  " ",
		".xml":  "󰗀 ",
		".htm":  " ",
		".html": " ",
		".txt":  " ",
		".mp3":  " ",
		".m4a":  " ",
		".ogg":  " ",
		".flac": " ",
		".mp4":  " ",
		".mkv":  " ",
		".webm": "󰃽 ",
		".zip":  "󰿺 ",
		".tar":  "󰛫 ",
		".gz":   "󰛫 ",
		".bz2":  "󰿺 ",
		".xz":   "󰿺 ",
		".jar":  " ",
		".java": " ",
		".js":   " ",
		".json": " ",
		".py":   " ",
		".rs":   " ",
		".yml":  " ",
		".yaml": " ",
		".toml": " ",
		".deb":  " ",
		".md":   " ",
		".rb":   " ",
		".php":  " ",
		".pl":   " ",
		".svg":  "󰜡 ",
		".eps":  " ",
		".ps":   " ",
		".git":  " ",
		".zig":  " ",
		".xbps": " ",
		".el":   " ",
		".vim":  " ",
		".lua":  " ",
		".pdf":  " ",
		".epub": "󰂺 ",
		".conf": " ",
		".iso":  " ",
        ".exe":  " ",
        ".odt":  "󰷈 ",
        ".ods":  "󰱾 ",
        ".odp":  "󰈧 ",
        ".gif":  "󰵸 ",
        ".tiff": "󰋪 ",
        ".7z":   " ",
        ".bat":  " ",
        ".app":  " ",
        ".log":  " ",
        ".sql":  " ",
        ".db":   " ",
		".org":  " ",
	}
)

func main() {
	args := os.Args[1:]
	nonFlagArgs, hasFlags, hasSpecificFlags := parseFlags(args)

	if showVersion {
		fmt.Println(version)
		return
	}

	var directory string
	var fileExtension string

	if len(nonFlagArgs) > 0 {
		directory = nonFlagArgs[0]
	}
	if len(nonFlagArgs) > 1 {
		fileExtension = strings.TrimPrefix(filepath.Ext(nonFlagArgs[1]), ".")
	}

	if directory == "" {
		directory = "."
	}

	var files []os.DirEntry
	var err error

	if fileExtension != "" {
		files, err = listFilesWithExtension(directory, fileExtension)
		if err != nil {
			log.Fatalf("Error listing files with extension %s: %v", fileExtension, err)
		}
	} else {
		files, err = os.ReadDir(directory)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(files) == 0 {
		fmt.Println("No files found.")
		return
	}

	if !showHidden {
		files = filterHidden(files)
	}

	if showOnlySymlinks {
		files = filterSymlinks(files, directory)
	}

	if orderBySize {
		sort.Slice(files, func(i, j int) bool {
			info1, _ := files[i].Info()
			info2, _ := files[j].Info()
			return info1.Size() < info2.Size()
		})
	}

	if orderByTime {
		sort.Slice(files, func(i, j int) bool {
			info1, _ := files[i].Info()
			info2, _ := files[j].Info()
			return info1.ModTime().Before(info2.ModTime())
		})
	}

	if recursiveListing {
		printTree(directory, "", true, 0, maxDepth)
	} else if longListing {
		printLongListing(files, directory, humanReadable)
	} else if fileSize {
		getFileSize(files, directory, humanReadable, dirOnLeft)
	} else {
		printFilesInColumns(files, directory, dirOnLeft, showSummary)
	}

	if (hasSpecificFlags && !longListing) || !hasFlags {
		fmt.Println()
	}
}

func listFilesWithExtension(dir string, ext string) ([]os.DirEntry, error) {
	var result []os.DirEntry

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), "."+ext) {
			result = append(result, entry)
		}
	}
	return result, nil
}

func filterHidden(entries []os.DirEntry) []os.DirEntry {
	var result []os.DirEntry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			result = append(result, entry)
		}
	}
	return result
}

func filterSymlinks(entries []os.DirEntry, dir string) []os.DirEntry {
	var result []os.DirEntry
	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if info, err := os.Lstat(fullPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
			result = append(result, entry)
		}
	}
	return result
}

func parseFlags(args []string) ([]string, bool, bool) {
	var nonFlagArgs []string
	hasFlags := false
	hasSpecificFlags := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 1 && arg[0] == '-' {
			hasFlags = true
			for j := 1; j < len(arg); j++ {
				switch arg[j] {
				case 'l':
					longListing = true
				case 'h':
					humanReadable = true
				case 's':
					fileSize = true
				case 'o':
					orderBySize = true
					hasSpecificFlags = true
				case 't':
					orderByTime = true
					hasSpecificFlags = true
				case 'm':
					showOnlySymlinks = true
					hasSpecificFlags = true
				case 'a':
					showHidden = true
					hasSpecificFlags = true
				case 'r':
					recursiveListing = true
				case 'i':
					dirOnLeft = true
					hasSpecificFlags = true
				case 'c':
					oneColumn = true
				case 'f':
					showSummary = true
				case 'v':
					showVersion = true
				case 'd':
					// Check if there is a number immediately after 'd'
					if j+1 < len(arg) && arg[j+1] >= '0' && arg[j+1] <= '9' {
						depthValue := arg[j+1:]
						maxDepthValue, err := strconv.Atoi(depthValue)
						if err != nil {
							fmt.Println("Invalid value for -d")
							os.Exit(1)
						}
						maxDepth = maxDepthValue
						hasSpecificFlags = true
						// Skip the rest of the string since we processed the depth value
						break
					} else if i+1 < len(args) {
						// Check if the next argument is the depth value
						depthValue := args[i+1]
						maxDepthValue, err := strconv.Atoi(depthValue)
						if err != nil {
							fmt.Println("Invalid value for -d")
							os.Exit(1)
						}
						maxDepth = maxDepthValue
						hasSpecificFlags = true
						// Skip the next argument since it's the depth value
						i++
						break
					} else {
						fmt.Println("Missing value for -d")
						os.Exit(1)
					}
				default:
					showHelp()
					os.Exit(1)
				}
			}
		} else {
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}
	return nonFlagArgs, hasFlags, hasSpecificFlags
}

func showHelp() {
	fmt.Println()
	fmt.Println("Usage: gols [FLAG] [DIRECTORY] [FILES]")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println()
    fmt.Println("	-?        Help")
	fmt.Println()
    fmt.Println("	-a        Show Hidden files")
    fmt.Println("	-c        Don't use spacing, print all files in one column")
    fmt.Println("	-f        Show summary of directories and files")
	fmt.Println("	-h        Human-readable file sizes")
    fmt.Println("	-i        Show directory icon on left")
	fmt.Println("	-l        Long listing format")
    fmt.Println("	-m        Only symbolic links are showing")
    fmt.Println("	-o        Sort by size")
    fmt.Println("	-r d n    Tree like listing, set the depth of the directory tree (n is an integer)")
	fmt.Println("	-s        Print files size")
    fmt.Println("	-t        Order by time")
    fmt.Println("	-v        Version")
	fmt.Println()
}

func printFilesInColumns(files []os.DirEntry, directory string, dirOnLeft bool, showSummary bool) {
	maxFilesInLine := 4
	maxFileNameLength := 19

	filesInLine := 0
	dirCount := 0
	fileCount := 0

	for _, file := range files {
		if file.IsDir() {
			dirCount++
			if dirOnLeft {
				fmt.Print(blue + "  " + file.Name() + reset)
			} else {
				printFile(file, directory)
			}
		} else {
			fileCount++
			printFile(file, directory)
		}

		if !oneColumn {
			filesInLine++
			if filesInLine >= maxFilesInLine || len(file.Name()) > maxFileNameLength {
				fmt.Println()
				filesInLine = 0
			} else {
				printPadding(file.Name(), maxFileNameLength)
			}
		} else {
			fmt.Println()
		}
	}

	if showSummary {
		fmt.Println()
		fmt.Printf("Directories: %s%d%s\n", blue, dirCount, reset)
		fmt.Printf("Files: %s%d%s\n", red, fileCount, reset)
	}
}

func getFileSize(files []os.DirEntry, directory string, humanReadable, dirOnLeft bool) {
    const sizeFieldWidth = 10
    const spaceBetweenSizeAndIcon = 2

    for _, file := range files {
        info, err := file.Info()
        if err != nil {
            log.Fatal(err)
        }

        size := info.Size()
        sizeStr := formatSize(size, humanReadable)

        sizeStr = fmt.Sprintf("%*s", sizeFieldWidth, sizeStr)

        fmt.Print(sizeStr)
        for i := 0; i < spaceBetweenSizeAndIcon; i++ {
            fmt.Print(" ")
        }

        if file.IsDir() {
            if dirOnLeft {
                fmt.Println(blue + " " + file.Name() + reset)
            } else {
                fmt.Println(blue + file.Name() + " " + reset)
            }
        } else {
            fmt.Println(getFileIcon(file, info.Mode(), directory) + " " + file.Name())
        }
    }
	if showSummary {
		fileCount, dirCount := countFilesAndDirs(files)
		fmt.Printf("Directories: %s%d%s\n", blue, dirCount, reset)
		fmt.Printf("Files: %s%d%s\n", red, fileCount, reset)
	}
}

func padRight(str string, length int) string {
    for len(str) < length {
        str += " "
    }
    return str
}

func formatSize(size int64, humanReadable bool) string {
	const (
		_  = iota
		KB = 1 << (10 * iota) // 1024 bytes
		MB
		GB
		TB
	)

	if humanReadable {
		// Human-readable format
		switch {
		case size >= TB:
			return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
		case size >= GB:
			return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
		case size >= MB:
			return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
		case size >= KB:
			return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
		default:
			return fmt.Sprintf("%d B", size)
		}
	} else {
		// Plain format with units
		switch {
		case size >= TB:
			return fmt.Sprintf("%d TB", size)
		case size >= GB:
			return fmt.Sprintf("%d GB", size)
		case size >= MB:
			return fmt.Sprintf("%d MB", size)
		case size >= KB:
			return fmt.Sprintf("%d KB", size)
		default:
			return fmt.Sprintf("%d B", size)
		}
	}
}

func printLongListing(files []os.DirEntry, directory string, humanReadable bool) {
	maxLen := map[string]int{
		"permissions": 0,
		"size":        0,
		"owner":       0,
		"group":       0,
		"month":       0,
		"day":         0,
		"time":        0,
		"linkTarget":  0,
	}

	var filteredFiles []os.DirEntry
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		permissions := formatPermissions(file, info.Mode(), directory)
		size := info.Size()
		sizeStr := formatSize(size, humanReadable) // Use unified formatSize function
		owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
		if err != nil {
			log.Fatal(err)
		}
		group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Gid))
		if err != nil {
			log.Fatal(err)
		}
		modTime := info.ModTime()
		month := modTime.Format("Jan")
		day := fmt.Sprintf("%2d", modTime.Day())
		timeStr := modTime.Format("15:04:05 2006")

		maxLen["permissions"] = max(maxLen["permissions"], len(permissions))
		maxLen["size"] = max(maxLen["size"], len(sizeStr))
		maxLen["owner"] = max(maxLen["owner"], len(owner.Username))
		maxLen["group"] = max(maxLen["group"], len(group.Name))
		maxLen["month"] = max(maxLen["month"], len(month))
		maxLen["day"] = max(maxLen["day"], len(day))
		maxLen["time"] = max(maxLen["time"], len(timeStr))

		if file.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
			if err == nil {
				maxLen["linkTarget"] = max(maxLen["linkTarget"], len(linkTarget)+5)
			}
		}

		filteredFiles = append(filteredFiles, file)
	}

	for _, file := range filteredFiles {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		permissions := formatPermissions(file, info.Mode(), directory)
		size := info.Size()
		sizeStr := formatSize(size, humanReadable) // Use unified formatSize function
		owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
		if err != nil {
			log.Fatal(err)
		}
		group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Gid))
		if err != nil {
			log.Fatal(err)
		}
		modTime := info.ModTime()
		month := modTime.Format("Jan")
		day := fmt.Sprintf("%2d", modTime.Day())
		timeStr := modTime.Format("15:04:05 2006")

		permissions = green + permissions + reset
		sizeStr = fmt.Sprintf("%*s", maxLen["size"], sizeStr)
		ownerStr := cyan + owner.Username + reset
		groupStr := brightBlue + group.Name + reset
		monthStr := magenta + month + reset
		dayStr := magenta + day + reset
		timeStr = magenta + timeStr + reset

		line := fmt.Sprintf(
			"%-*s  %s  %-*s  %-*s %-*s %-*s %-*s %s %s",
			maxLen["permissions"], permissions,
			sizeStr,
			maxLen["owner"], ownerStr,
			maxLen["group"], groupStr,
			maxLen["month"], monthStr,
			maxLen["day"], dayStr,
			maxLen["time"], timeStr,
			getFileIcon(file, info.Mode(), directory), file.Name(),
		)

		if file.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
			if err == nil {
				line += fmt.Sprintf(" %s==> %s%s", cyan, linkTarget, reset)
			}
		}

		fmt.Println(line)
	}

	if showSummary {
		fileCount, dirCount := countFilesAndDirs(files)
		fmt.Printf("Directories: %s%d%s\n", blue, dirCount, reset)
		fmt.Printf("Files: %s%d%s\n", red, fileCount, reset)
	}
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func countFilesAndDirs(files []os.DirEntry) (int, int) {
    fileCount := 0
    dirCount := 0
    for _, file := range files {
        if file.IsDir() {
            dirCount++
        } else {
            fileCount++
        }
    }
    return fileCount, dirCount
}

func formatPermissions(file os.DirEntry, mode os.FileMode, directory string) string {
    perms := make([]byte, 10)
    for i := range perms {
        perms[i] = '-'
    }

    if file.Type()&os.ModeSymlink != 0 {
        linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
        if err == nil {
            symlinkTarget := filepath.Join(directory, linkTarget)
            targetInfo, err := os.Stat(symlinkTarget)
            if err == nil && targetInfo.IsDir() {
                perms[0] = 'l'
                perms[1] = 'd'
            } else {
                perms[0] = 'l'
            }
        }
    } else if mode.IsDir() {
        perms[0] = 'd'
    }

    for i, s := range []struct {
        bit os.FileMode
        char byte
    }{
        {0400, 'r'}, {0200, 'w'}, {0100, 'x'},
        {0040, 'r'}, {0020, 'w'}, {0010, 'x'},
        {0004, 'r'}, {0002, 'w'}, {0001, 'x'},
    } {
        if mode&s.bit != 0 {
            perms[i+1] = s.char
        }
    }

    return string(perms)
}

func rwx(perm os.FileMode) string {
	var b strings.Builder

	if perm&0400 != 0 {
		b.WriteString("r")
	} else {
		b.WriteString("-")
	}
	if perm&0200 != 0 {
		b.WriteString("w")
	} else {
		b.WriteString("-")
	}
	if perm&0100 != 0 {
		if perm&os.ModeSetuid != 0 {
			b.WriteString("s")
		} else {
			b.WriteString("x")
		}
	} else {
		if perm&os.ModeSetuid != 0 {
			b.WriteString("S")
		} else {
			b.WriteString("-")
		}
	}

	return b.String()
}

func printFile(file os.DirEntry, directory string) {
	info, err := file.Info()
	if err != nil {
		log.Fatal(err)
	}
	if file.IsDir() {
		fmt.Print(blue + file.Name() + " " + reset)
	} else {
		fmt.Print(getFileIcon(file, info.Mode(), directory) + file.Name())
	}
	fmt.Print(" ")
}

func getFileIcon(file os.DirEntry, mode os.FileMode, directory string) string {
	if file.Type()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
		if err == nil {
			symlinkTarget := filepath.Join(directory, linkTarget)
			targetInfo, err := os.Stat(symlinkTarget)
			if err == nil && targetInfo.IsDir() {
				return brightMagenta + " " + reset
			} else {
				return brightCyan + " " + reset
			}
		}
	}

	if mode.IsDir() {
		return blue + " " + reset
	}

	ext := filepath.Ext(file.Name())
	icon, exists := fileIcons[ext]
	if exists {
		switch ext {
		case ".go":
			return cyan + icon + reset
        case ".sh":
			if mode&os.ModePerm&0111 != 0 {
				return brightGreen + icon + reset
			} else {
				return white + icon + reset
			}
		case ".cpp", ".hpp", ".cxx", ".hxx":
			return blue + icon + reset
		case ".css":
			return lightBlue + icon + reset
		case ".c", ".h":
			return blue + icon + reset
		case ".cs":
			return darkMagenta + icon + reset
		case ".png", ".jpg", ".jpeg", ".JPG", ".webp":
			return darkBlue + icon + reset
		case ".gif":
			return magenta + icon + reset
		case ".xcf":
			return magenta + icon + reset
		case ".xml":
			return lightCyan + icon + reset
		case ".htm", ".html":
			return orange + icon + reset
		case ".txt", ".app":
			return white + icon + reset
		case ".mp3", ".m4a", ".ogg", ".flac":
			return brightBlue + icon + reset
		case ".mp4", ".mkv", ".webm":
			return darkMagenta + icon + reset
		case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z":
			return lightPurple + icon + reset
		case ".jar", ".java":
			return orange + icon + reset
		case ".js":
			return yellow + icon + reset
		case ".json", ".tiff":
			return brightYellow + icon + reset
		case ".py":
			return darkYellow + icon + reset
		case ".rs":
			return darkGray + icon + reset
		case ".yml", ".yaml":
			return brightRed + icon + reset
		case ".toml":
			return darkOrange + icon + reset
		case ".deb":
			return lightRed + icon + reset
		case ".md":
			return cyan + icon + reset
		case ".rb":
			return red + icon + reset
		case ".php":
			return brightBlue + icon + reset
		case ".pl":
			return red + icon + reset
		case ".svg":
			return lightPurple + icon + reset
		case ".eps", ".ps":
			return orange + icon + reset
		case ".git":
			return orange + icon + reset
		case ".zig":
			return darkOrange + icon + reset
		case ".xbps":
			return darkGreen + icon + reset
		case ".el":
			return magenta + icon + reset
		case ".vim":
			return darkGreen + icon + reset
		case ".lua", ".sql":
			return brightBlue + icon + reset
		case ".pdf", ".db":
			return brightRed + icon + reset
		case ".epub":
			return cyan + icon + reset
		case ".conf", ".bat":
			return darkGray + icon + reset
		case ".iso":
			return gray + icon + reset
		case ".exe":
			return brightCyan + icon + reset
		case ".org":
			return darkMagenta + icon + reset
		default:
			return icon
		}
	}

	if mode&os.ModePerm&0111 != 0 {
		return green + " " + reset
	}

	return " " + reset
}

func printPadding(name string, maxFileNameLength int) {
	padding := maxFileNameLength - len(name)
	fmt.Print(strings.Repeat(" ", padding))
}

func getFileNameAndExtension(file os.DirEntry) (string, string) {
	ext := filepath.Ext(file.Name())
	name := strings.TrimSuffix(file.Name(), ext)
	return name, ext
}

func printTree(path, prefix string, isLast bool, currentDepth, maxDepth int) {
	if maxDepth != -1 && currentDepth > maxDepth {
		return
	}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var filteredFiles []os.DirEntry
	for _, file := range files {
		if showHidden || !strings.HasPrefix(file.Name(), ".") {
			filteredFiles = append(filteredFiles, file)
		}
	}

	for i, file := range filteredFiles {
		isLastFile := i == len(filteredFiles)-1
		if isLastFile {
			fmt.Printf("%s└── ", prefix)
		} else {
			fmt.Printf("%s├── ", prefix)
		}

		printFile(file, path)
		fmt.Println()

		if file.IsDir() {
			newPrefix := prefix
			if isLastFile {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(filepath.Join(path, file.Name()), newPrefix, isLastFile, currentDepth+1, maxDepth)
		}
	}

	if showSummary && currentDepth == 0 {
		fileCount, dirCount := countFilesAndDirs(files)
		fmt.Printf("Directories: %s%d%s\n", blue, dirCount, reset)
		fmt.Printf("Files: %s%d%s\n", red, fileCount, reset)
	}
}
