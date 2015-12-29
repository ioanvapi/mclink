package main

/*
    Do all these as administrator.

    1. go build executable mclink.exe
	2. Create system environment variable %minecraft_home% with value the path to minecraft.
	3. create a link 'minecraft.lnk' to 'minecraft_launcher.exe' in the same folder

    4. download nssm and use it to install mclink as service (nssm install mclink)
    4.1 in nssm's gui provide path to mclink.exe

    5. go to windows firewall -> inbound rules -> create New rule for 8123 port
    6. go to services and start mclink service
    7. test app with http://localhost:8123 and check desktop

*/

import (
    "net/http"
    "log"
    "os"
    "fmt"
    "io/ioutil"
    "text/template"
    "bufio"
    "strings"
    "os/exec"
)

const (
    envMinecraft = "minecraft_home"
    envPublic = "public"
    linkName = "minecraft.lnk"
    homeTemplate = `
<html><body>
<button type="button" style="height:300px; width:500px" onclick="location.href='/d'">Delete Link</button><br/><br/>
<button type="button" style="height:300px; width:500px" onclick="location.href='/a'">Add Link</button><br/><br/>
<button type="button" style="height:300px; width:500px" onclick="location.href='/kill'">KILL</button><br/><br/>

{{range .}}
<p>{{.}}</p><br/>
{{end}}
</body></html>`
)

var (
    desktopLink string
    minecraftPath string
)


func init() {
    public, ok := os.LookupEnv(envPublic)
    if !ok {
        log.Panicf("Cannot find environment variable '%s'.", envPublic)
    }

    desktopLink = public + string(os.PathSeparator) + "Desktop" + string(os.PathSeparator) + linkName

    mcPath, ok := os.LookupEnv(envMinecraft)
    if !ok {
        log.Panicf("Cannot find environment variable '%s'.", envMinecraft)
    }

    minecraftPath = mcPath + string(os.PathSeparator) + linkName
}

func main() {
//    log.Printf("Desktop link path is '%s'", desktopLink)

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/a", addLinkHandler)
    http.HandleFunc("/d", removeLinkHandler)
    http.HandleFunc("/kill", killProcessHandler)

    log.Println("Server started ...")
    err := http.ListenAndServe(":8123", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
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


func renderTemplate(w http.ResponseWriter, message ...string) {
    t, err := template.New("home").Parse(homeTemplate)
    if err != nil {
        fmt.Fprintln(w, "Error")
        return
    }

    t.Execute(w, message)
}


func minecraftPIDs() ([]string, error) {

    out, err := exec.Command("TASKLIST", "/V").Output()
    if err != nil {
        return nil, err
    }

    pids := make([]string, 0)
    scanner := bufio.NewScanner(strings.NewReader(string(out)))

    for scanner.Scan() {
        line := scanner.Text()
        if !strings.HasPrefix(line, "javaw.exe") {
            continue
        }

        tokens := strings.Fields(line)
        if !strings.HasPrefix(tokens[9], "Minecraft") &&
           !strings.HasPrefix(tokens[10], "Minecraft"){
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