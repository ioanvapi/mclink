package main

/*
	1. Create system environment variable %minecraft_path% with value the path to minecraft.
	2. create a link 'minecraft.lnk' to 'minecraft_launcher.exe' in the same folder

	//regedit.exe -> key below -> add exec as 'data' field
	HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Run
*/

import (
	"net/http"
	"log"
	"os"
	"fmt"
	"io/ioutil"
)

const (
	envMinecraft = "minecraft_path"
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
	log.Printf("Desktop link path is '%s'", desktopLink)

	// add link
	http.HandleFunc("/a", addLink)
	// create link
	http.HandleFunc("/c", addLink)

	// remove link
	http.HandleFunc("/r", removeLink)
	// delete link
	http.HandleFunc("/d", removeLink)

	// kill the process
	http.HandleFunc("/kill", killProcess)


	log.Println("Server started ...")
	err := http.ListenAndServe(":8123", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}



func addLink(w http.ResponseWriter, req *http.Request) {
	// check link already exists on desktop
	if _, err := os.Stat(desktopLink); err == nil {
		fmt.Fprintf(w, "Link '%s' already exist at '%s'.", linkName, desktopLink)
		return
	}

	data, err := ioutil.ReadFile(minecraftPath)
	if err != nil {
		fmt.Fprintf(w, "Cannot read file '%s'. \n '%s'.", minecraftPath, err.Error())
		return
	}
	// Write data to dst
	err = ioutil.WriteFile(desktopLink, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "Cannot write to file '%s'. \n '%s'.", desktopLink, err.Error())
		return
	}

	fmt.Fprintf(w, "Link successful added.")
	return
}

func removeLink(w http.ResponseWriter, req *http.Request) {
	// check link does not exist on desktop
	if _, err := os.Stat(desktopLink); os.IsNotExist(err) {
		fmt.Fprintf(w, "Link '%s' does not exist on desktop.", linkName)
		return
	}

	err := os.Remove(desktopLink)
	if err != nil {
		fmt.Fprintf(w, "Error removing the link '%s' from desktop: \n %s", linkName, err.Error())
		return
	}

	fmt.Fprintf(w, "Link successful removed from desktop.")
	return
}


func killProcess(w http.ResponseWriter, req *http.Request) {
	//todo
}
