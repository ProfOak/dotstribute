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

func shouldIgnore(path string, ignoredFiles []string) bool {
	for _, ignoredFile := range ignoredFiles {
		if strings.HasPrefix(path, ignoredFile) {
			return true
		}
	}
	return false
}

func isSymlink(symlinkPath string) bool {
	// Stat checks original file, Lstat checks link itself.
	f, err := os.Lstat(symlinkPath)
	if err != nil {
		log.Fatal(err)
	}
	return f.Mode()&os.ModeSymlink != 0
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
	// First one is `.` and we don't want to skip it because filepath.SkipDir
	// will exclude the whole path under the current working directory.
	return dotfiles[1:]
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

func symlink(dotfile string, preview bool, ask bool) {
	symlinkPath := homePath(dotfile)
	dotfile, err := filepath.Abs(dotfile)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s -> %s\n", dotfile, symlinkPath)
	if preview {
		return
	}
	if _, err := os.Stat(symlinkPath); !errors.Is(err, os.ErrNotExist) {
		fmt.Println("File already exists")
		return
	}

	if !ask || shouldContinue() {
		dir, _ := filepath.Split(symlinkPath)
		err = os.MkdirAll(dir, os.ModePerm)
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
func unsymlink(dotfile string, preview bool, ask bool) {
	var err error
	symlinkPath := homePath(dotfile)

	if _, err := os.Stat(symlinkPath); os.IsNotExist(err) {
		log.Printf("Skipping: No symlink '%s' in home directory.", symlinkPath)
		return
	}
	if !isSymlink(symlinkPath) {
		fmt.Printf("%s is not a symlink, ignoring.\n", symlinkPath)
		return
	}
	fmt.Printf("Removing: %s\n", symlinkPath)
	if preview {
		return
	}

	if !ask || shouldContinue() {
		err = os.Remove(symlinkPath)
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
	fmt.Print("Continue? (y/N): ")
	fmt.Scanln(&choice)
	return strings.ToLower(choice) == "y"
}

func main() {
	var (
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
		return
	}

	ignoredFiles := getIgnoredFiles(".dotignore")
	dotfiles := getDotFiles(".", ignoredFiles)
	for _, f := range dotfiles {
		action(f, preview, ask)
	}
}
