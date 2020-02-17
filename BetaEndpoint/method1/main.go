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

// FulfillmentBody is the overall structure to be marshalled into the json.
// There is a very particular format the webhook is expecting.
type FulfillmentBody struct {
	FulfillmentText     FulfillmentText
	FulfillmentMessages FulfillmentMessages
}

// FulfillmentMessages holds information about the latency, download speed and upload speed from a remote server.
type FulfillmentMessages struct {
	Latency  string `json:"Latency"`
	Download string `json:"Download"`
	Upload   string `json:"Upload"`
}

// FulfillmentText holds the fulfillment text response expected by the webhook
type FulfillmentText struct {
	Text string `json:"Text"`
}

var fBody = FulfillmentBody{}
var fText = FulfillmentText{}
var fMessages = FulfillmentMessages{}

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

// // alphaTest runs the speed test on the ComblastAlpha EC2 and returns a json containign the results
// func alphaTest(w http.ResponseWriter, r *http.Request) {
// 	var err error
// 	var dfr *df.Request
// 	//var p params

// 	decoder := json.NewDecoder(req.Body)
// 	if err = decoder.Decode(&dfr); err != nil {
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	defer req.Body.Close()

// 	// Filter on action, using a switch for example

// 	// Retrieve the params of the request
// 	/*	if err = dfr.GetParams(&p); err != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		// Retrieve a specific context
// 		if err = dfr.GetContext("my-awesome-context", &p); err != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 	*/
// 	// Do things with the context you just retrieved
// 	fmt.Println("test")
// 	dff := &df.Fulfillment{
// 		FulfillmentMessages: df.Messages{
// 			df.ForGoogle(df.SingleSimpleResponse("hello", "hello")),
// 			{RichMessage: df.Text{Text: []string{"hello"}}},
// 		},
// 	}
// 	rw.Header().Set("Content-Type", "application/json")
// 	rw.WriteHeader(http.StatusOK)
// 	json.NewEncoder(rw).Encode(dff)

// 	// run ssh command
// 	cmd := exec.Command("ssh", "-i", awsK, awsP, awsC)

// 	// Read the output into a string and manipulate it to make it readable
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(fmt.Sprint(err) + "; " + string(output))
// 		return
// 	}
// 	res := strings.Split(string(output), "\n")
// 	fMessages.Latency = strings.Replace(res[0], "Latency: ", "", -1)
// 	fMessages.Download = strings.Replace(res[1], "Download: ", "", -1)
// 	fMessages.Upload = strings.Replace(res[2], "Upload: ", "", -1)
// 	fText.Text = "Here are the speed test results:"

// 	// Fill in the body of the fulfillment message
// 	fBody.FulfillmentMessages = fMessages
// 	fBody.FulfillmentText = fText

// 	// Opens the logs json file, this will create one if none are present
// 	alphaL, err := os.OpenFile("alphaLogs.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatal("ERROR OPENING JSON LOGS: ", err)
// 		return
// 	}
// 	defer alphaL.Close()

// 	// Marshall the raw data into the logs file
// 	raw, err := json.MarshalIndent(fMessages, "", "    ")
// 	if err != nil {
// 		log.Fatal("ERROR MARSHALLING JSON CURRENTS: ", err)
// 		return
// 	}
// 	alphaL.Write(raw)
// 	alphaL.Close()

// 	// Opens the currents json file, this will create one if none are present
// 	alphaC, err := os.OpenFile("alphaCurrent.json", os.O_RDWR|os.O_CREATE, 0666)
// 	if err != nil {
// 		log.Fatal("ERROR OPENING JSON CURRENTS: ", err)
// 		return
// 	}
// 	defer alphaC.Close()

// 	// Marshall the formated response into the response file
// 	bod, err := json.MarshalIndent(fBody, "", "    ")
// 	if err != nil {
// 		log.Fatal("ERROR MARSHALLING JSON CURRENTS: ", err)
// 		return
// 	}
// 	alphaC.Write(bod)
// 	alphaC.Close()

// 	// return with currents json file
// 	temp, err := template.ParseFiles("alphaStatic.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Execute
// 	temp.Execute(w, r)
// }
