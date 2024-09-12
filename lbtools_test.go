package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

//func Test_RestartProgramByName(t *testing.T) {
//	programName := "excel"
//	err := RestartProgramByName(programName)
//	if err != nil {
//		t.Fatalf("Could not stop and start %s due to: %s", programName, err)
//	}
//}

// starts kills an excel process if all works well
func Test_KillProcessByName(t *testing.T) {
	programName := "excel"
	processName := strings.ToUpper(programName) + ".EXE"
	var err error
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	// First we spin up the process and test the StartProgramByName() func
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Printf("Starting %s\n", programName)
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
		fmt.Printf("Killing %s\n", programName)
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
