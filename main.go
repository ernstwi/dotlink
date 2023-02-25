package main

import (
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

func main() {
	sourceDir, err := filepath.Abs(os.Args[1])
	targetDir, err := filepath.Abs(os.Args[2])
	fatal(err)

	sourceDirParts := strings.Split(sourceDir, string(os.PathSeparator))
	targetDirParts := strings.Split(targetDir, string(os.PathSeparator))

	fatal(filepath.WalkDir(sourceDir,
		func(sourceEntry string, d os.DirEntry, err error) error {
			fatal(err)

			if d.Name()[0] == '.' {
				return fs.SkipDir
			}

			if _, ln := link(d.Name()); ln || !d.IsDir() {

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

				fatal(os.MkdirAll(filepath.Dir(targetEntry), 0777))
				if err := os.Symlink(sourceEntry, targetEntry); err != nil {
					fmt.Println("exists")
				}

				fmt.Println()

				// Do not recurse into linked dirs
				if d.IsDir() {
					return fs.SkipDir
				}
			}

			return nil
		}))
}
