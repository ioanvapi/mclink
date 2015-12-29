package main

import (
	"os/exec"
	"fmt"
	"strings"
	"bufio"
)


func minecraftPIDs() ([]string, error) {
	out, err := exec.Command("TASKLIST", "/V").Output()
	if err != nil {
		return nil, err
	}

	nrTokens := 11
	pids := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "javaw.exe") {
			continue
		}

		tokens := strings.Fields(line)
		if len(tokens) < nrTokens {
			continue
		}

		if !strings.HasPrefix(tokens[nrTokens - 2], "Minecraft") &&
		!strings.HasPrefix(tokens[nrTokens - 1], "Minecraft") {
			continue
		}

		pids = append(pids, tokens[1])
	}

	if len(pids) == 0 {
		return nil, fmt.Errorf("Cannot find any Minecraft process.")
	}

	return pids, nil
}


func killProcesses(pids []string) ([]string) {
	messages := make([]string, 0)

	for _, pid := range pids {
		out, err := exec.Command("TASKKILL", "/PID", pid, "/F").Output()
		if err != nil {
			messages = append(messages, err.Error())
		} else {
			messages = append(messages, string(out))
		}
	}
	return messages
}
