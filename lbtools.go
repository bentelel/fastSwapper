package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
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

func Remove[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func KillProcess(name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, err := p.Name()
		// Some processes do not let me access their names p.e. "Secure System". for those we need to skip ahead.
		// for now we dont handle err and just skip ahead
		if err != nil {
			// t.Fatalf("Error fetching process name while looking for %s: %s", processName, err)
			continue
		}
		if n == name {
			return p.Kill()
		}
	}
	return fmt.Errorf("process not found")
}
