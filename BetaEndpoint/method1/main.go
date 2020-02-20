package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"text/template"
)

func main() {
	http.HandleFunc("/", homePage)
	// http.HandleFunc("/alpha", alphaTest)
	http.HandleFunc("/webhook", webHook)
	http.HandleFunc("/webhook/latency", webLatency)
	http.HandleFunc("/webhook/download", webDownload)
	http.HandleFunc("/webhook/upload", webUpload)
	http.ListenAndServe(":80", nil)
	// err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}

var awsP = "ubuntu@ec2-54-147-41-245.compute-1.amazonaws.com"
var awsK = "./comblastalpha.pem"

var commandFull = "./speedtest-go | grep -E \"Latency|Download:|Upload:\""
var commandLate = "./speedtest-go | grep -i \"Latency:\" | awk '{print $2}'"
var commandDown = "./speedtest-go | grep -i \"Download:\" | awk '{print $2}'"
var commandUplo = "./speedtest-go | grep -i \"Upload:\" | awk '{print $2}'"

// homePage is mostly for redirection
func homePage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
	temp.Execute(w, nil)
}

// webHook runs the speed test on the ComblastAlpha EC2 and return a string containing the full results
func webHook(w http.ResponseWriter, r *http.Request) {
	// run the speed test with an ssh command
	cmd := exec.Command("ssh", "-i", awsK, awsP, commandFull)

	// Read the output into a string
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + "; " + string(output))
		return
	}

	// Return the output
	fmt.Fprintf(w, string(output))
}

// webLatency runs the speed test on the ComblastAlpha EC2 and return a string containing the results of latency
func webLatency(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ssh", "-i", awsK, awsP, commandLate)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + "; " + string(output))
		return
	}
	fmt.Fprintf(w, string(output))
}

// webHook runs the speed test on the ComblastAlpha EC2 and return a string containing the download speed
func webDownload(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ssh", "-i", awsK, awsP, commandDown)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + "; " + string(output))
		return
	}
	fmt.Fprintf(w, string(output))
}

// webHook runs the speed test on the ComblastAlpha EC2 and return a string containing the upload speed
func webUpload(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ssh", "-i", awsK, awsP, commandUplo)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + "; " + string(output))
		return
	}
	fmt.Fprintf(w, string(output))
}
