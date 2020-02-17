// Copyright 2017, Google, Inc.
// Licensed under the Apache License, Version 2.0 (the 'License');
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an 'AS IS' BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

'use strict';

const http = require('http');
const functions = require('firebase-functions');

const host = 'http://ec2-54-92-157-95.compute-1.amazonaws.com/webhook';

exports.dialogflowFirebaseFulfillment = functions.https.onRequest((req, res) => {
 

  // Which parameter do we want to test for?
  let endpoint = '/download'; // For now we will only run the download test
  // if (req.body.queryRestult.parameters['endpoint']) {
  //   specifics = req.body.queryRestult.parameters['endpoint'];
  //   console.log('Specifics: ' + endpoint);
  // }

  // Call the speedtest API
  callSpeedTestAPI(host, endpoint).then((output) => {
    res.json({ 'fulfillmentText': output }); // Return the results of the Speed Test API to Dialogflow
  }).catch(() => {
    res.json({ 'fulfillmentText': 'Error accessing the Speed Test API, check connections.'}); // For connection breaks
  });

});

function callSpeedTestAPI(host, endpoint) {
  return new Promise((resolve, reject) => {
    // Create the path for the HTTP request to get the results of the speed test
    let path = host + endpoint;
    console.log('API Request: ' + host + path); 

    // Make the HTTP request to get the results of the speed test
    http.get({host: host, path: path}, (res) => {
      let body = ''; 
      res.on('data', (d) => { body += d; });
      res.on('end', () =>{
        // After all the data has been received parse the JSON for desired data
        let response = JSON.parse(body);
        let download = response['data']['download'][0];
        
        // let latency = response['data']['latency'][0];
        // let upload = response['data']['upload'][0];

        // Create response
        let output = `Here are the speed test results for ${ip['type']} : \n
        Download Speed: ${download['download']} \n`;

        // let output = `Here are the speed test results for ${ip['type']} : \n
        // Latency: ${latency['latency']} \n
        // Download Speed: ${download['download']} \n
        // Upload Speed: ${upload['upload']}`;

        // Resolve the promise with the output text
        console.log(output);
        resolve(output);
      });
      res.on('error', (error) => {
        console.log(`Error calling the speed test API: ${error}`)
        reject();
      });
    });
  });
}