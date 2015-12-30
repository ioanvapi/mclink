package main

import (
    "os/exec"
    "fmt"
    "strings"
    "bufio"
    "log"
    "syscall"
)

const (
    markerToken = "minecraft"
    appToken = "javaw.exe"
)


func minecraftPIDs() ([]string, error) {
    var pids []string

    out, err := exec.Command("TASKLIST", "-V", "-FI", "\"javaw.exe\"").Output()
//    out, err := exec.Command("cmd", "/C", "TASKLIST", "/V").Output()
    if err != nil {
        return pids, err
    }

    log.Println(string(out))

    scanner := bufio.NewScanner(strings.NewReader(string(out)))

    // scan output line by line
    for scanner.Scan() {
        // split a line in tokens
        line := scanner.Text()
        tokens := strings.Fields(line)
        // there must be at least 2 tokens (PID is the second) and first should be appToken
        if len(tokens) < 2 || !strings.HasPrefix(tokens[0], appToken) {
            continue
        }

        // search for marker and get the PID
        for _, token := range tokens {
            if strings.HasPrefix(strings.ToLower(token), markerToken) {
                pids = append(pids, tokens[1])
                log.Println(line)
            }
        }
    }

    if len(pids) == 0 {
        return pids, fmt.Errorf("Cannot find any Minecraft process.")
    }

    log.Println(pids)
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
