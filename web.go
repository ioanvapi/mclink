package main
import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"html/template"
)


const (
	homeTemplate = `
<html>
<head>
	<meta charset="UTF-8">
	<title>Minecraft shortcut</title>
	<style>
	button {
		height:250px;
		width:400px;
		font-size:40px
	}
	</style>
</head>
<body>
<div><button type="button" onclick="location.href='/d'">Delete Link</button></div>
<div><button type="button" onclick="location.href='/a'">Add Link</button></div>
<div><button type="button" onclick="location.href='/kill'">KILL</button></div>

{{range .}}<div><p>{{.}}</p></div><br/>{{end}}
</body></html>`
)

func renderTemplate(w http.ResponseWriter, message ...string) {
	t, err := template.New("home").Parse(homeTemplate)
	if err != nil {
		fmt.Fprintln(w, "Error")
		return
	}

	t.Execute(w, message)
}


func messageMiddleware(message string, next func(w http.ResponseWriter, req *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if len(message) > 0 {
			renderTemplate(w, message)
			return
		}
		next(w, req)
	})
}

func addLinkHandler(w http.ResponseWriter, req *http.Request) {
	// check link already exists on desktop
	if _, err := os.Stat(desktopLink); err == nil {
		msg := fmt.Sprintf("Link '%s' already exist at '%s'.", linkName, desktopLink)
		renderTemplate(w, msg)
		return
	}

	data, err := ioutil.ReadFile(minecraftPath)
	if err != nil {
		msg := fmt.Sprintf("Cannot read file '%s'. \n '%s'.", minecraftPath, err.Error())
		renderTemplate(w, msg)
		return
	}
	// Write data to dst
	err = ioutil.WriteFile(desktopLink, data, 0666)
	if err != nil {
		msg := fmt.Sprintf("Cannot write to file '%s'. \n '%s'.", desktopLink, err.Error())
		renderTemplate(w, msg)
		return
	}

	renderTemplate(w, fmt.Sprintf("Link successful added."))
}


func removeLinkHandler(w http.ResponseWriter, req *http.Request) {
	// check link does not exist on desktop
	if _, err := os.Stat(desktopLink); os.IsNotExist(err) {
		msg := fmt.Sprintf("Link '%s' does not exist on desktop.", linkName)
		renderTemplate(w, msg)
		return
	}

	err := os.Remove(desktopLink)
	if err != nil {
		msg := fmt.Sprintf("Error removing the link '%s' from desktop: \n %s", linkName, err.Error())
		renderTemplate(w, msg)
		return
	}

	renderTemplate(w, fmt.Sprintf("Link successful removed from desktop."))
}


func killProcessHandler(w http.ResponseWriter, req *http.Request) {
	pids, err := minecraftPIDs()
	if err != nil {
		renderTemplate(w, err.Error())
		return
	}

	messages := killProcesses(pids)
	renderTemplate(w, messages...)

}


func homeHandler(w http.ResponseWriter, req *http.Request) {
	renderTemplate(w, "")
}
