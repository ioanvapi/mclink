package main
import (
    "net/http"
    "fmt"
    "html/template"
    "strconv"
    "strings"
    "log"
    "time"
)


const (
    homeTemplate = `
<html>
<head>
	<meta charset="UTF-8">
	<title>Minecraft shortcut</title>
	<style>
	button {
		height:150px;
		width:350px;
		font-size:40px
	}
	input {
		height:50px;
		width:350px;
		font-size:30px
	}
	</style>
</head>
<body>
<div><button type="button" onclick="location.href='/d'">Delete Link</button></div>
<div><button type="button" onclick="location.href='/a'">Add Link</button></div>
<div><button type="button" onclick="location.href='/kill'">KILL</button></div>
<div><button type="button" onclick="schedule()">Schedule</button></div>
<div><input type="text" id="duration"></div>
<div>{{range .}}<p>{{.}}</p>{{end}}</div>

<script>
  function schedule() {
    var d = document.getElementById("duration").value;
    var xhttp;
	if (window.XMLHttpRequest) {
		xhttp = new XMLHttpRequest();
	} else {
		// code for IE6, IE5
		xhttp = new ActiveXObject("Microsoft.XMLHTTP");
	}
  	xhttp.open("GET", "/s?d=" + d, true);
  	xhttp.send();
  }
</script>
</body>
</html>`
)



func renderTemplate(w http.ResponseWriter, message ...string) {
    t, err := template.New("home").Parse(homeTemplate)
    if err != nil {
        fmt.Fprintln(w, "Error")
        return
    }

    t.Execute(w, message)
}


func messageMiddleware(messages []string, next func(w http.ResponseWriter, req *http.Request)) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        if len(messages) > 0 {
            renderTemplate(w, messages...)
            return
        }
        next(w, req)
    })
}

func addLinkHandler(w http.ResponseWriter, req *http.Request) {
    err := addLink()
    if err != nil {
        renderTemplate(w, err.Error())
        return
    }

    renderTemplate(w, fmt.Sprintf("Link successful added to desktop."))
}


func removeLinkHandler(w http.ResponseWriter, req *http.Request) {
    err := removeLink()
    if err != nil {
        renderTemplate(w, err.Error())
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

func scheduleHandler(w http.ResponseWriter, req *http.Request) {
    // check and parse request parameter
    dStr := req.URL.Query().Get("d")
    d, err := strconv.Atoi(strings.TrimSpace(dStr))
    if err != nil {
        log.Printf("Cannot convert scheduled duration value '%s'", dStr)
        return
    }

    // set Scheduler and the action taken after that duration
    dm := time.Duration(d) * time.Minute
    Scheduler.Reset(d, func() {
        pids, err := minecraftPIDs()
        if err != nil {
            log.Println(err.Error())
            return
        }
        killProcesses(pids)
        removeLink()

        log.Printf("Scheduled Kill executed for pids: '%v'", pids)
    })

    log.Printf("Kill scheduled after '%v' for time: '%s'", dm, time.Now().Add(dm).Format("15:04"))
}


func homeHandler(w http.ResponseWriter, req *http.Request) {
    msg := ""
    if Scheduler.When() != nil {
        msg = fmt.Sprintf("Kill schedulet at: '%s'", Scheduler.When().Format("15:04"))
    }
    renderTemplate(w, msg)
}
