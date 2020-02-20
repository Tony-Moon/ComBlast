# Beta Endpoint

The software in this folder will all the source code for the beta EC2 endpoint.
The Beta EC2 makes an ssh connection to Alpha EC2 and runs an internet speed test.
Then, The Beta EC2 will gather, log and report that information.
The Beta EC2 will have API endpoints accessable via HTTP (current) and HTTPS (in progress). 
These endpoints will be able to give specific or generic data based on the endpoints. 

There are currently two methods of setting up these endpoints, so they will have their respective directories.