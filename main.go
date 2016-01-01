package main

/*
    Do all these as administrator.

    - go build executable mclink.exe
	- Create system environment variable %minecraft_home% with value the path to minecraft.
	- Create system environment variable %mclink_log% with value the path to a log file.
	- create a link 'minecraft.lnk' to 'minecraft_launcher.exe' in the same folder

    - download nssm and use it to install mclink as service (nssm install mclink)
    - in nssm's gui provide path to mclink.exe

    - go to windows firewall -> inbound rules -> create New rule for 8123 port
    - go to services and start mclink service with some admin account
    - test app with http://localhost:8123 and check desktop

*/

import (
	"net/http"
	"log"
	"os"
	"fmt"
)

const (
	envMclinkLog = "mclink_log"
	port = 8123
)

var (
	startupErrors []string
)


func main() {
	logPath, ok := os.LookupEnv(envMclinkLog)
	if ok {
		f, err := os.OpenFile(logPath, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		if err != nil {
			startupErrors = append(startupErrors, err.Error())
		} else {
			defer f.Close()
			log.SetOutput(f)
		}
	}

	http.Handle("/", messageMiddleware(startupErrors, homeHandler))
	http.Handle("/a", messageMiddleware(startupErrors, addLinkHandler))
	http.Handle("/d", messageMiddleware(startupErrors, removeLinkHandler))
	http.Handle("/stop", messageMiddleware(startupErrors, stopProcessHandler))
	http.Handle("/s", messageMiddleware(startupErrors, scheduleHandler))

	log.Printf("Server listening on port %d.\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


