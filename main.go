package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var processName string
	var dryRun bool

	var rootCmd = &cobra.Command{
		Use:   "kp",
		Short: "A CLI to kill processes by regex",
		Run: func(cmd *cobra.Command, args []string) {
			if processName == "" {
				fmt.Println("Please provide a process name regex using --process-name or -p")
				return
			}

			processes, err := listProcesses()
			if err != nil {
				fmt.Printf("Error listing processes: %v\n", err)
				return
			}

			matchedProcesses := filterProcesses(processes, processName)
			if dryRun {
				fmt.Println("Dry run mode. The following processes would be killed:")
				for _, p := range matchedProcesses {
					fmt.Printf("PID: %s, Name: %s\n", p.PID, p.Name)
				}
			} else {
				if len(matchedProcesses) > 0 {
					fmt.Println("The following processes match the regex:")
					for _, p := range matchedProcesses {
						fmt.Printf("PID: %s, Name: %s\n", p.PID, p.Name)
					}

					fmt.Print("Do you want to kill these processes? (y/n): ")
					reader := bufio.NewReader(os.Stdin)
					response, _ := reader.ReadString('\n')
					response = strings.TrimSpace(response)

					if strings.ToLower(response) == "y" {
						for _, p := range matchedProcesses {
							fmt.Printf("Killing process PID: %s, Name: %s\n", p.PID, p.Name)
							err := killProcess(p.PID)
							if err != nil {
								fmt.Printf("Failed to kill process PID: %s, error: %v\n", p.PID, err)
							}
						}
					} else {
						fmt.Println("No processes were killed.")
					}
				} else {
					fmt.Println("No matching processes found.")
				}
			}
		},
	}

	rootCmd.Flags().StringVarP(&processName, "process-name", "p", "", "Regex to match process names")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode, do not kill processes")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Process struct {
	PID  string
	Name string
}

func listProcesses() ([]Process, error) {
	out, err := exec.Command("ps", "-e", "-o", "pid,comm").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var processes []Process
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			processes = append(processes, Process{PID: fields[0], Name: fields[1]})
		}
	}
	return processes, nil
}

func filterProcesses(processes []Process, regex string) []Process {
	var matchedProcesses []Process
	r, err := regexp.Compile(regex)
	if err != nil {
		fmt.Printf("Invalid regex: %v\n", err)
		return matchedProcesses
	}

	for _, p := range processes {
		if r.MatchString(p.Name) {
			matchedProcesses = append(matchedProcesses, p)
		}
	}
	return matchedProcesses
}

func killProcess(pid string) error {
	cmd := exec.Command("kill", pid)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill process PID: %s, error: %w", pid, err)
	}
	return nil
}
