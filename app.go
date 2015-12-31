package main

import (
	"os/exec"
	"fmt"
	"strings"
	"bufio"
	"log"
	"strconv"
	"os"
	"io/ioutil"
)

const (
	envMinecraftHome = "minecraft_home"
	envPublic = "public"
	linkName = "minecraft.lnk"
)


var (
	desktopLink string
	minecraftPath string
)

func init() {
	public, ok := os.LookupEnv(envPublic)
	if !ok {
		startupErrors = append(startupErrors, fmt.Sprintf("Cannot find environment variable '%s'.", envPublic))
	}

	desktopLink = public + string(os.PathSeparator) + "Desktop" + string(os.PathSeparator) + linkName

	mcPath, ok := os.LookupEnv(envMinecraftHome)
	if !ok {
		startupErrors = append(startupErrors, fmt.Sprintf("Cannot find environment variable '%s'.", envMinecraftHome))
	}

	minecraftPath = mcPath + string(os.PathSeparator) + linkName
}

func removeLink() error {
	// check if link does not exist on desktop
	if _, err := os.Stat(desktopLink); os.IsNotExist(err) {
		return fmt.Errorf("Link '%s' does not exist on desktop.", linkName)
	}

	err := os.Remove(desktopLink)
	if err != nil {
		return fmt.Errorf("Error removing the link '%s' from desktop: \n %s", linkName, err.Error())
	}

	return nil
}


func addLink() error {
	// check link already exists on desktop
	if _, err := os.Stat(desktopLink); err == nil {
		return fmt.Errorf("Link '%s' already exist at '%s'.", linkName, desktopLink)
	}

	data, err := ioutil.ReadFile(minecraftPath)
	if err != nil {
		return fmt.Errorf("Cannot read file '%s'. \n '%s'.", minecraftPath, err.Error())
	}
	// Write data to dst
	err = ioutil.WriteFile(desktopLink, data, 0666)
	if err != nil {
		return fmt.Errorf("Cannot write to file '%s'. \n '%s'.", desktopLink, err.Error())
	}

	return nil
}

func minecraftPIDs() ([]string, error) {
	var pids []string

	// get pid and main window title for every javax process
	out, err := exec.Command("powershell", "get-process javaw | select-object id,mainwindowtitle").Output()
	if err != nil {
		log.Println(err.Error())
		return pids, fmt.Errorf("Cannot find any Minecraft process.")
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
