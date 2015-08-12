# GrassMudHorse
`GrassMudHorse`  is a tcp tunnel that continuously send `ICMP` `PING`s to target servers to check latency and packet drop rate, then choose the best server and tunnel data to it.

It can be used to choose and balance multiple servers.



eg. If you have mutltiple shadowsocks server that outside GFW, and one fast server inside GFW, you can use the server inside GFW as a forwarding server, run `GrassMudHorse` on it, it will choose the best server and make a tunnel to the shadowsocks server.

` shadowsocks client ==> inside-gfw-server:port ==> the "chosen" remote shadowsocks server:port `

`GrassMudHorse` only forwards data, and will not decrypt anything. So you should set encrypt method and encrypt key all the same between different shadowsocks servers, and also the shadowsocks client.


## Usage:
first you should rename sample.json to config.json.

here's the values in sample.json means:

+ `Listen` : Local listening `IP:Port`

+ `Interval` : Interval between each `PING`s to one server in millisecond. If you have 5 remote servers and set `Interval` to 1000, `GrassMudHorse` will trigger 5 `PING`s/s.

+ `Timeout` : `PING` time out in millisecond.

+ `IPv6` : Use ipv6 ping or not.

+ `PayloadSize` : The payload size of `ICMP` `PING` packets. Don't set it larger than your network's `mtu-`28. FYI, the `mtu` of `pppoe` network is 1492.

+ `HistorySize` : The number of recent ping datas ( latency and timeouted-or-not infos ) GrassMudHorse keeps to caculate the average latency and droprate. A smaller number may cause the change of chosen server too often.

+ `Lua` : Filename of Lua script. `GrassMudHorse` will use native score function if this value is not set. See below for the details of Lua Script.

+ `Servers` : list of `IP:PORT` of remote servers.


## Lua Script

You can write your own score function in Lua to override the default one.

The Lua script file must have a function named `score()`, and returns a number. A bigger number means the server is better.

The script will be called once per server in a interval of one second.  You can use `droprate() averagelatency() address()` to get current server's state , caculate and return the score. Then the server which has max score will be chosen as remote server.


## Debug info
try run `GrassMudHorse` with flag --?
