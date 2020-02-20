/*
Copyright 2017, Google, Inc.
Licensed under the Apache License, Version 2.0 (the 'License');
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.
	
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an 'AS IS' BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This code has been adapted by Anthony Moon

'use strict';

const http = require('http');
const functions = require('firebase-functions');
const host = 'http://ec2-54-92-157-95.compute-1.amazonaws.com/webhook';

exports.dialogflowFirebaseFulfillment = functions.https.onRequest((req, res) => {
  // Call the speed test API
  callSpeedTestAPI(host).then((output) => {
    //Return the results of the Speed Test API to Dialogflow
	res.json({ 'fulfillmentText': output });
  }).catch(() => {
    res.json({ 'fulfillmentText': 'Error accessing the Speed Test API, check connections.'}); // For connection breaks
  });

});

// callSpeedTestAPI simply curls an API to get the results of speed test on a remote AWS server
function callSpeedTestAPI(host) {
	return new Promise((resolve, reject) => {
    // Create the path for the HTTP request to get the results of the speed test
    console.log('API Request: ' + host); 

    // Make the HTTP request to get the results of the speed test
    http.get({host: host}, (res) => {
		res.setEncoding('utf8');
		let body = ''; // body will store the response chunks
		
		// Store each chuck to body one by one
		res.on('data', (d) => { body += d; });
		
		// After all the data has been received parse the JSON for desired data
		res.on('end', () =>{
			let response = JSON.parse(body);
			console.log(body);

			// Resolve the promise with the output text
			console.log(response);
			resolve(body);
		});
	  
		// Error handeling 
		res.on('error', (error) => {
			console.log(`Error calling the speed test API: ${error}`);
			reject();
		});
	});
  });
}