# GrassMudHorse
GrassMudHorse

a tcp tunnel that continuously send icmp pings to target server(s) to check latency and packet drop rate, then pick the best server and tunnel data to the preset port.

use for choose and balance multiple servers.

eg. if you have mutltiple shadowsocks server that outside GFW, and one fast server inside GFW, you can use the server inside gfw as a forwarding server, run this program on it, it will choose the best server and make a tunnel to the shadowsocks server.

shadowsocks client ==> inside-gfw-server:port ==> the "choosed" remote shadowsocks server:port

this program only forwards data, it will not decrypt anything, so you should set encrypt method and encrypt key all the same between different shadowsocks servers, and also the shadowsocks client.
