package main

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func Typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsStringWord(sliceToCheckAgainst []string, wordToCheck string) bool {
	// this is probably highly inefficient, as we are looping over the complete list for each rune in wordToCheck, but whatever, well refactor later
	for _, r := range wordToCheck {
		if ContainsString(sliceToCheckAgainst, string(r)) {
			return true
		}
	}
	return false
}

func IsDir(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetDirsInDir(dir string) []string {
	// Returns slice of strings containing all directories within given directory
	// param dir: string -- directory to check
	entries, err := os.ReadDir(filepath.FromSlash(dir))
	if err != nil {
		log.Fatal(err)
	}
	result := make([]string, 0)
	// translate all vfiles to strings
	for _, e := range entries {
		if e.IsDir() {
			result = append(result, e.Name())
		}
	}
	return result
}

func GetAllInDir(dir string) []string {
	// Returns slice of strings containing all directories within given directory
	// param dir: string -- directory to check
	entries, err := os.ReadDir(filepath.FromSlash(dir))
	if err != nil {
		log.Fatal(err)
	}
	result := make([]string, 0)
	// translate all files to strings
	for _, e := range entries {
		result = append(result, e.Name())
	}
	return result
}

func All[T any](ts []T, pred func(T) bool) bool {
	// Takes in a slice and function (which returns a bool) and does the func check on all elements of slice
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}
	return true
}

func CombineString(s []string) (string, error) {
	var err error
	var ret string
	for _, w := range s {
		ret += w + " "
	}
	ret = TrimSuffix(ret, " ")
	return ret, err
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
