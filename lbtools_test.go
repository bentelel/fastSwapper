package main

import (
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

func Test_Typeof(t *testing.T) {
	anInt := 1
	aString := "test"
	aSlice := []int{1, 2, 3}
	aSliceAsText := "[]int{1, 2, 3}"
	wantIntType := "int"
	wantStringType := "string"
	wantSliceType := "[]int"
	if got := Typeof(anInt); got != wantIntType {
		t.Fatalf("Type to string conversion failed for %s. \nExpected: %s\n Got: %s\n", strconv.Itoa(anInt), wantIntType, got)
	}
	if got := Typeof(aString); got != wantStringType {
		t.Fatalf("Type to string conversion failed for %s. \nExpected: %s\n Got: %s\n", aString, wantStringType, got)
	}
	if got := Typeof(aSlice); got != wantSliceType {
		t.Fatalf("Type to string conversion failed for %s. \nExpected: %s\n Got: %s\n", aSliceAsText, wantSliceType, got)
	}
}

func Test_ContainsStringWord(t *testing.T) {
	slice := []string{"A", "B", "C"}
	word := "TestWordA"
	want := true
	got := ContainsStringWord(slice, word)
	if want != got {
		t.Fatalf("Positive test failed,\nExpected: %s\nGot: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	word = "TestWord"
	want = false
	got = ContainsStringWord(slice, word)
	if want != got {
		t.Fatalf("Negative test failed,\nExpected: %s\nGot: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
}

func Test_ContainsString(t *testing.T) {
	haystack := []string{"A", "B", "C"}
	needle := "C"
	want := true
	got := ContainsString(haystack, needle)
	if want != got {
		t.Fatalf("Positive test failed,\nExpected: %s\nGot: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	needle = "D"
	want = false
	got = ContainsString(haystack, needle)
	if want != got {
		t.Fatalf("Negative test failed,\nExpected: %s\nGot: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
}

func Test_Exists(t *testing.T) {
	// ideally we would create a test dir/file here and test against that and then delete it.
	// however i couldnt make that work, so for now we test against the Tagetik folder.
	want := true
	path := "C:\\Tagetik\\"
	if got := Exists(path); want != got {
		t.Fatalf("Could not find file %s even though it exists. Expected Exists() to return %s but got %s.\n", path, strconv.FormatBool(want), strconv.FormatBool(got))
	}
}

// func Test_GetDirsInDir() > how to test, for this we would need to create a dir and some dirs underneath and then test that and afterwards delete that again..
// same for Test_GetAllInDir()..

func Test_All(t *testing.T) {
	fun := func(n int) bool { return n >= 0 }
	nums := []int{-1, -1, 0, 1}
	want := false
	got := All(nums, fun)
	if want != got {
		t.Fatalf("Check failed, expected %s, got %s.\n", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	nums = []int{1, 1, 0, 1}
	want = true
	got = All(nums, fun)
	if want != got {
		t.Fatalf("Check failed, expected %s, got %s.\n", strconv.FormatBool(want), strconv.FormatBool(got))
	}
}

func Test_main(t *testing.T) {
	// run test functions as subtests so they run sequencially. We do this because both tests test against the Excel-process and might run into raceconditions if run without waiting each other out.
	t.Run("Restart Test", func(t *testing.T) { TRestartProgramByName(t) })
	t.Run("Kill Test", func(t *testing.T) { TKillProcessByName(t) })
}

// kills excel (if running) and then restarts it (and then kills it again so the test leaves no trace)
func TRestartProgramByName(t *testing.T) {
	programName := "excel"
	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		t.Logf("Restarting %s.\n", programName)
		err := RestartProgramByName(programName)
		errChan <- err
	}(&wg)
	wg.Wait()
	close(errChan)
	if err, ok := <-errChan; ok && err != nil {
		t.Fatalf("Could not stop and start %s due to: %s", programName, err)
	}
	processName := strings.ToUpper(programName) + ".EXE"
	t.Logf("Killing %s agian.\n", programName)
	err := KillProcessByName(processName)
	if err != nil {
		t.Fatalf("Could not kill %s, error: \n%s", processName, err)
	}
}

// starts kills an excel process if all works well
func TKillProcessByName(t *testing.T) {
	programName := "excel"
	processName := strings.ToUpper(programName) + ".EXE"
	var err error
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	// First we spin up the process and test the StartProgramByName() func
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		t.Logf("Starting %s\n", programName)
		err = StartProgramByName(programName)
		errChan <- err
	}(&wg)
	wg.Wait()
	close(errChan)
	if err, ok := <-errChan; ok && err != nil {
		t.Fatalf("Could not start program %s, error: \n%s\n", programName, err)
	}
	// If we made it here the program is started.
	// now we can go about terminating it.
	// first check if we can get and access the process at all
	processes, err := process.Processes()
	if err != nil {
		t.Fatalf("Error fetching processes to check if %s is running: \n%s\n", processName, err)
	}
	var processNames []string
	for _, p := range processes {
		n, err := p.Name()
		// Some processes do not let me access their names p.e. "Secure System". for those we need to skip ahead.
		// for now we dont handle err and just skip ahead
		if err != nil {
			// t.Fatalf("Error fetching process name while looking for %s: %s", processName, err)
			continue
		}
		processNames = append(processNames, n)
	}
	if !ContainsString(processNames, processName) {
		t.Fatalf("%s not found in running processes, cannot commence test.", processName)
	}
	errChan2 := make(chan error, 1)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		t.Logf("Killing %s\n", programName)
		err = KillProcessByName(processName)
		errChan2 <- err
	}(&wg)
	wg.Wait()
	close(errChan2)
	if err, ok := <-errChan2; ok && err != nil {
		t.Fatalf("Could not kill process %s, error: %s", processName, err)
	}
	processes, err = process.Processes()
	if err != nil {
		t.Fatalf("Error fetching processes to check if %s is still running: %s", processName, err)
	}
	for _, p := range processes {
		n, err := p.Name()
		// Some processes do not let me access their names p.e. "Secure System". for those we need to skip ahead.
		// for now we dont handle err and just skip ahead
		if err != nil {
			// t.Fatalf("Error fetching process name while looking for %s: %s", processName, err)
			continue
		}
		if n == processName {
			t.Fatalf("Process %s still running. KillProcess did not work.", processName)
		}
	}
}
