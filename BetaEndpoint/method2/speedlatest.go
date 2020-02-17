package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	df "github.com/leboncoin/dialogflow-go-webhook"
)

var awsP = "ubuntu@ec2-54-147-41-245.compute-1.amazonaws.com"
var awsK = "./comblastalpha.pem"
var commandDown = "./speedtest-go | grep -i \"Download:\" | awk '{print $2}'"

/*type params struct {
	City   string `json:"city"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}
*/

func webhook(rw http.ResponseWriter, req *http.Request) {
	// run the speed test with an ssh command
	cmd := exec.Command("ssh", "-i", awsK, awsP, commandDown)

	// Read the output into a string
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + "; " + string(output))
		return
	}
	var dfr *df.Request
	//var p params

	decoder := json.NewDecoder(req.Body)
	if err = decoder.Decode(&dfr); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Filter on action, using a switch for example

	// Retrieve the params of the request
	/*	if err = dfr.GetParams(&p); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}


		// Retrieve a specific context
		if err = dfr.GetContext("my-awesome-context", &p); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	*/

	// Do things with the context you just retrieved
	fmt.Println("Download Speed Test")
	dff := &df.Fulfillment{
		FulfillmentMessages: df.Messages{
			df.ForGoogle(df.SingleSimpleResponse(string(output), string(output))),
			{RichMessage: df.Text{Text: []string{"download speed"}}},
		},
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(dff)
}

func main() {
	http.HandleFunc("/webhook", webhook)
	log.Fatal(http.ListenAndServe(":80", nil))
	//log.Fatal(http.ListenAndServeTLS(":443", "https-server.crt", "https-server.key", nil))
}
