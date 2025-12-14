# Open-Tunnel: Simple Reverse TCP Tunnel

**Goal:** Create a lightweight, `ngrok`-style open tunnel to securely share local files or expose services via a raw TCP connection.

## Server Configuration (VM)

The tunnel server utilizes three ports to establish the proxy connection: 

Note: deployed on e2-micro instance on Google Compute Engine


| Port | Role | Description |
| :--- | :--- | :--- |
| **9000** | **Control** | Persistent link for the server to send commands   to the client. |
| **9001** | **Public** | The exposed port that receives all incoming internet traffic. |
| **9002** | **Data** | The channel used for proxying the actual application data. |

## Execution Flow

1.  **Startup:** Client runs `opentunnel 8080` (default: 8080) and establishes a permanent connection to **Server:9000**.
2.  **Public Hit:** An external user connects to the tunnel's public address (e.g., `http://35.224.59.81:9001`).
3.  **Command Sent:** Server:9001 receives the request and sends a `"new"` command over the Control Channel (**9000**) to the client.
4.  **Client Connects:** The client immediately opens two new connections:
    * Outbound Data Connection to **Server:9002**.
    * Local Connection to the service on **localhost:8080**.
5.  **Proxying:** The server stitches the Public Connection (9001) to the new Data Connection (9002), creating a real-time, bidirectional proxy pipe between the internet user and the local service.


This pipe remains active until either end closes the connection.


## Why Build?

Out of curiosity, I started digging into how ngrok functions internally. Realizing that my needs were limited to file serving, I figured it made more sense to implement a simplified version myself.