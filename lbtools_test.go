package main

import (
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

func Test_KillProcess_KillingProcess(t *testing.T) {
	processName := "EXCEL.EXE"
	// first check if we can get and access the process at all
	processes, err := process.Processes()
	if err != nil {
		t.Fatalf("Error fetching processes to check if %s is running: %s", processName, err)
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

	err = KillProcess(processName)
	if err != nil {
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
