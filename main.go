package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"log"
	"os"
	"path/filepath"
	"time"
)

type dirMap map[string]uint32

func picDate(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()
	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, fmt.Errorf("decode file: %w", err)
	}
	t, err := x.DateTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("parse EXIF date: %w", err)
	}
	return t, nil
}

func picSplit(path string) (dirMap, error) {
	var mappedDirs = make(dirMap)

	entries, err := os.ReadDir(path)
	if err != nil {
		return mappedDirs, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		source := filepath.Join(path, entry.Name())
		t, err := picDate(source)
		if err != nil {
			return mappedDirs, fmt.Errorf("retrieve picture date for %s: %w", entry.Name(), err)
		}
		dir := t.Format("2006-01-02")

		destinationDir := filepath.Join(path, dir)
		destination := filepath.Join(path, dir, entry.Name())
		if _, ok := mappedDirs[dir]; !ok {
			err = os.Mkdir(destinationDir, 0700)
			if err != nil {
				return mappedDirs, fmt.Errorf("create directory %s: %w", source, err)
			}
			mappedDirs[dir] = 0
		}

		fmt.Println(source, "=>", destination)
		err = os.Rename(source, destination)
		if err != nil {
			return mappedDirs, fmt.Errorf("move %s: %w", entry.Name(), err)
		}
		mappedDirs[dir]++
	}

	return mappedDirs, nil
}

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Println("Split pictures into sub-directories by camera shooting date")
		fmt.Println("Need path to directory")
		os.Exit(1)
	}

	path := os.Args[1]
	if path == "." {
		path, err = os.Getwd()
		if err != nil {
			fmt.Println("Split pictures into sub-directories by camera shooting date")
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	mappedDirs, err := picSplit(path)
	if err != nil {
		log.Fatalf("Unable to split pictures: %s", err)
	}
	fmt.Println("Done")
	for pName, count := range mappedDirs {
		fmt.Printf("- Moved %d pictures to %s\n", count, pName)
	}
	fmt.Println()
}
