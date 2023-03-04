package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func fatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func concat[T any](slices [][]T) []T {
	var size int
	for _, s := range slices {
		size += len(s)
	}
	res := make([]T, size)
	var i int
	for _, s := range slices {
		i += copy(res[i:], s)
	}
	return res
}

func link(s string) (string, bool) {
	if len(s) > 5 && s[:5] == "link-" {
		return s[5:], true
	}
	return s, false
}

func dot(s string) string {
	if len(s) > 4 && s[:4] == "dot-" {
		return "." + s[4:]
	}
	return s
}

func parseArgs(args []string) (string, string, bool) {
	var source, target string
	rm := false
	switch len(args) {
	case 3:
		source = args[1]
		target = args[2]
	case 4:
		if args[1] == "--rm" {
			rm = true
		}
		source = args[2]
		target = args[3]
	default:
		fatal(errors.New("Bad input"))
	}
	sourceDir, err := filepath.Abs(source)
	targetDir, err := filepath.Abs(target)
	fatal(err)
	return sourceDir, targetDir, rm
}

func sliceHasPrefix[T comparable](s, prefix []T) bool {
	if len(prefix) > len(s) {
		return false
	}
	for i := range prefix {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

func main() {
	sourceDir, targetDir, rm := parseArgs(os.Args)

	sourceDirParts := strings.Split(sourceDir, string(os.PathSeparator))
	targetDirParts := strings.Split(targetDir, string(os.PathSeparator))

	fatal(filepath.WalkDir(sourceDir,
		func(sourceEntry string, d os.DirEntry, err error) error {
			fatal(err)

			// Dir or leaf starting with `.`
			if d.Name()[0] == '.' {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}

			// Normal dir
			if _, ln := link(d.Name()); d.IsDir() && !ln {
				return nil
			}

			// link-dir or leaf
			sourceEntryParts := strings.Split(sourceEntry, string(os.PathSeparator))
			sourceTailParts := sourceEntryParts[len(sourceDirParts):]
			targetTailParts := make([]string, len(sourceTailParts))
			for i, s := range sourceTailParts {
				res, _ := link(s)
				res = dot(res)
				targetTailParts[i] = res
			}
			targetEntry := filepath.Join(concat([][]string{
				{"/"},
				targetDirParts,
				targetTailParts})...)
			fmt.Println(sourceEntry)
			fmt.Println(targetEntry)

			if rm { // Remove symlink
				// Verify it is pointing to sourceDir
				unlinked, err := filepath.EvalSymlinks(targetEntry)
				fatal(err)
				abs, err := filepath.Abs(unlinked)
				if !sliceHasPrefix(strings.Split(abs, string(os.PathSeparator)), sourceDirParts) {
					fmt.Println("Not pointing to source")
				} else {
					fatal(os.Remove(targetEntry))
				}
			} else { // Create symlink
				fatal(os.MkdirAll(filepath.Dir(targetEntry), 0777))
				if err := os.Symlink(sourceEntry, targetEntry); err != nil {
					fmt.Println("exists")
				}
			}
			fmt.Println()

			// Do not recurse into subdirs
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}))
}
