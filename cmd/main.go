package main

import (
	"diary-renamer/dict"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const DEFAULT_OUTPUT_LAYOUT = "02-01-2006"

var dir string
var outputLayout string

func init() {
	flag.StringVar(&dir, "dir", "", "Directory with files to rename (required)")
	flag.StringVar(&outputLayout, "o", DEFAULT_OUTPUT_LAYOUT, "Output layout (default to 02-01-2006)")
	flag.Parse()
}

func main() {
	if dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		time, err := parseToTime(file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		formatted := time.Format(outputLayout)

		oldPath := fmt.Sprintf("%s/%s", dir, file.Name())
		newPath := fmt.Sprintf("%s/%s.md", dir, formatted)
		os.Rename(oldPath, newPath)
	}
}

func parseToTime(fileName string) (time.Time, error) {
	split := strings.Fields(strings.Replace(fileName, ".md", "", -1))

	if len(split) != 6 {
		return time.Time{}, fmt.Errorf("could not rename file %s", fileName)
	}

	day := split[1]
	month := split[3]
	year := split[5]

	month, err := translateMonth(month)
	if err != nil {
		return time.Now(), err
	}

	date, err := time.Parse("2006 January 02", fmt.Sprintf("%s %s %s", year, month, day))
	if err != nil {
		return time.Now(), err
	}

	return date, nil
}

func translateMonth(ptMonth string) (string, error) {
	en, ok := dict.Months[strings.ToLower(ptMonth)]

	if !ok {
		return "", fmt.Errorf("could not translate %s", ptMonth)
	}

	return en, nil
}
