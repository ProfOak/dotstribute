package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

func getIgnoredFiles(dotignoreFile string) []string {
	ignoredFiles := []string{}
	buff, err := os.ReadFile(dotignoreFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Unable to find the .dotstribute file. Continuing without it.")
		} else if errors.Is(err, os.ErrPermission) {
			log.Println(".dotstribute file permission error. Continuing without it.")
		} else {
			log.Fatalln(err)
		}
	}

	for _, ignoredFile := range strings.Split(string(buff), "\n") {
		if ignoredFile != "" {
			ignoredFiles = append(ignoredFiles, ignoredFile)
		}
	}
	return ignoredFiles
}

func resolvePaths(path string) (string, string) {
	symlinkPath := homePath(path)
	dotfilePath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}
	return symlinkPath, dotfilePath
}

func shouldIgnore(path string, ignoredFiles []string) bool {
	for _, ignoredFile := range ignoredFiles {
		if strings.HasPrefix(path, ignoredFile) {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}

func isDir(path string) bool {
	if info, err := os.Stat(path); err != nil && info.IsDir() {
		return true
	}
	return false
}

func isCorrectSymlink(symlinkPath, realpath string) bool {
	// Stat checks original file, Lstat checks link itself.
	f, err := os.Lstat(symlinkPath)
	if err != nil {
		log.Fatal(err)
	}

	if f.Mode()&os.ModeSymlink == 0 {
		return false
	}

	if path, err := filepath.EvalSymlinks(symlinkPath); err == nil && path != realpath {
		return false
	}

	return true
}

func isHomeDir(directory string) bool {
	directory = filepath.Clean(directory)
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Unable to find the home directory.")
	}
	return directory == home
}

func getDotFiles(starting string, ignoredFiles []string) []string {
	var dotfiles []string
	err := filepath.WalkDir(starting, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalln(err)
		}
		if shouldIgnore(info.Name(), ignoredFiles) {
			if info.IsDir() {
				return filepath.SkipDir
			}
		} else if !info.IsDir() {
			// Only save leaves of the filepath tree.
			dotfiles = append(dotfiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	// Previously the first element was the current working directory because
	// of how the skipping logic worked but it seems to be working correctly
	// now. I've kept the check in case it ever happens again.
	if dotfiles[0] == "." {
		dotfiles = dotfiles[1:]
	}

	return dotfiles
}

func homePath(dotfile string) string {
	// Files might not have the necessary dot prefix in the dotfile repo.
	if !strings.HasPrefix(dotfile, ".") {
		dotfile = "." + dotfile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Unable to find the home directory.")
	}
	return filepath.Join(home, dotfile)
}

func symlink(symlinkPath, dotfile string, preview bool, ask bool) {
	if fileExists(symlinkPath) {
		if !isCorrectSymlink(symlinkPath, dotfile) {
			fmt.Print(" - Cannot create symlink")
		}
		return
	}

	if !ask || shouldContinue() {
		dir, _ := filepath.Split(symlinkPath)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("Unable to create the directory needed.")
			log.Fatalln(err)
		}
		err = os.Symlink(dotfile, symlinkPath)
		if err != nil {
			log.Println("Unable to symlink the file needed.")
			log.Fatalln(err)
		}
	}
}

// unsymlink will remove symlinked files from a user's home directory. It also
// removes the empty directory above any symlinked files if there are no more
// files stored there.
func unsymlink(symlinkPath, dotfile string, preview bool, ask bool) {
	if _, err := os.Stat(symlinkPath); os.IsNotExist(err) {
		log.Printf("Skipping: No symlink '%s' in home directory.", symlinkPath)
		return
	}
	if !isCorrectSymlink(symlinkPath, dotfile) {
		fmt.Println(symlinkPath, dotfile)
		fmt.Printf("%s is not a symlink, ignoring.\n", symlinkPath)
		return
	}
	fmt.Printf("Removing: %s\n", symlinkPath)

	if !ask || shouldContinue() {
		err := os.Remove(symlinkPath)
		if err != nil {
			log.Println("Unable to remove the symlink")
			log.Fatalln(err)
		}
	}

	// Delete the current directory if there are no more files in it.
	dir, _ := filepath.Split(symlinkPath)
	if isHomeDir(dir) {
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) > 0 {
		return
	}
	fmt.Printf("Removing: %s\n", dir)
	err = os.Remove(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func shouldContinue() bool {
	var choice string
	fmt.Print(" - Continue? (y/N): ")
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return false
	}
	return strings.ToLower(choice) == "y"
}

func main() {
	const (
		emojiFile   = "📄"
		emojiFolder = "📁"
		emojiLink   = "🔗"
		emojiShipit = "🚀"
		emojiError  = "❌"
	)

	var (
		emoji    string
		generate bool
		preview  bool
		unlink   bool
		ask      bool
	)
	flag.BoolVarP(&generate, "generate", "g", false, "")
	flag.BoolVarP(&preview, "preview", "p", false, "")
	flag.BoolVarP(&unlink, "unlink", "u", false, "")
	flag.BoolVarP(&ask, "ask", "a", false, "")
	flag.Parse()

	action := symlink
	if unlink {
		action = unsymlink
	}

	if generate {
		fmt.Println("TODO")
		return
	}

	ignoredFiles := getIgnoredFiles(".dotignore")
	dotfiles := getDotFiles(".", ignoredFiles)
	for _, dotfile := range dotfiles {
		symlinkPath, dotfilePath := resolvePaths(dotfile)
		if fileExists(symlinkPath) {
			if isDir(symlinkPath) {
				emoji = emojiFolder
			} else if isCorrectSymlink(symlinkPath, dotfilePath) {
				emoji = emojiLink
			} else {
				emoji = emojiError
			}
		} else {
			emoji = emojiShipit
		}

		fmt.Printf("%s %s", emoji, dotfile)
		if preview {
			continue
		}
		action(symlinkPath, dotfilePath, preview, ask)
		fmt.Println("")
	}
}
