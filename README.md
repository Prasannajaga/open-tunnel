# open-tunnel
Simple ngrok style open tunnel for share your files and expose your data through raw TCP connection


working module

Step 1:

there's simple server running in e2-micro VM (1GB RAM) in compute engine which will open three ports for reverse tunnel proxy

port action
9000 open connection with ur local running client  
9001 expose the connection globally 
9002 data transfer portal

step 2:

when u run the "opentunnel 8080" (assume 8080 should be the local server u want to expose globally)
you'll get the http:sdsjdsb:9001 some sample URL like this 

client will execute connection to VM:9000 keep it open
when u hit the http:sdsjdsb:9001 server will notify the client there's new request came so 
open a new TCP connection to 9002 

client -> 9002 stay connected 

this 9002 will be mapped to 9001
which routes all the traffic to 9002 when it's being hit 
this connection will kept open until there's intteruption happend if 








