package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var oval struct {
	operators string
	temp      string
}

func open(rawurl string) error {
	_, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return err
	}

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, rawurl)
	err = exec.Command(cmd, args...).Start()
	if err != nil {
		return errors.New("failed to open url" + err.Error())
	}
	return nil

}

func operations_for_server() string {
	fmt.Println(exec.Command(oval.operators))

	parts := strings.Split(oval.operators, " ")
	if len(parts) == 0 {
		fmt.Println("Empty command")
		return "empty"
	}

	// Create a new command with the parts
	cmd := exec.Command(parts[0], parts[1:]...)

	// Execute the command
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "Error 2"
	}

	// Print the output
	fmt.Println(string(output))
	return string(output)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Welcome" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello World!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	reader := bufio.NewReader(os.Stdin)
	var err error

	if err = r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	// The submit button was pressed
	oval.operators = r.FormValue("OperatorfCalc")
	fmt.Fprintf(w, "POST request successful!\n")
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if ip == "127.0.0.1" || ip == "::1" {
		fmt.Println("Hallo Localhost; keine Anmeldung erforderlich.")
		operations_for_server()
	} else {
		fmt.Println("Allow?")
		temp, error := reader.ReadString('\n')
		if error != nil {
			fmt.Println("Something went wrong with taking the input")
		}
		switch temp {

		case "y\n":
			fmt.Fprintf(w, operations_for_server())

		default:
			fmt.Fprintf(w, "your request was not accepted")
		}
	}
	//submitValue := r.FormValue("submit")
	//if submitValue != "" {

	//for oval.temp == "" {

	//}

}

func main() {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/Welcome", WelcomeHandler)

	fmt.Printf("Starting server at port 8080\n")
	err := open("http://127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
