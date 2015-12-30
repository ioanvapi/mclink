package main

import (
    "os/exec"
    "fmt"
    "strings"
    "bufio"
    "log"
    "strconv"
)

const (
    markerToken = "minecraft"
    appToken = "javaw.exe"
)


func minecraftPIDs() ([]string, error) {
    var pids []string

    out, err := exec.Command("powershell", "get-process javaw | select-object id,mainwindowtitle").Output()
    if err != nil {
        return pids, err
    }

    log.Printf("Tasklist output: '%s'\n", string(out))

    scanner := bufio.NewScanner(strings.NewReader(string(out)))

    // scan output line by line
    for scanner.Scan() {
        // split a line in tokens
        line := scanner.Text()
        tokens := strings.Fields(line)
        // there must be at least 2 tokens (PID is the 1st)
        if len(tokens) < 1 {
            continue
        }

        // running as service we don't get the 'window title' as running from cmd
        // and we cannot select minecraft app only. We have to kill all javaw
        token := strings.TrimSpace(tokens[0])
        _, err = strconv.Atoi(token)
        if err == nil {
            pids = append(pids, tokens[0])
        }
/*
        // search for marker and get the PID
        for _, token := range tokens {
            if strings.HasPrefix(strings.ToLower(token), markerToken) {
                pids = append(pids, tokens[0])
                log.Println(line)
            }
        }
*/
    }

    if len(pids) == 0 {
        return pids, fmt.Errorf("Cannot find any Minecraft process.")
    }

    log.Println(pids)
    return pids, nil
}

func minecraftPIDs3() ([]string, error) {
    var pids []string

    cmd := "TASKLIST"
    args := []string{
        "/V",
        "/FI",
        fmt.Sprintf("\"IMAGENAME eq %s\"", appToken),
    }

    out, err := exec.Command(cmd, args...).Output()
//    out, err := exec.Command("cmd", "/C", "TASKLIST", "/V").Output()
    if err != nil {
        return pids, err
    }

    log.Printf("Tasklist output: '%s'\n", string(out))

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
