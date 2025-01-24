English | [中文](README-cn.md)
# udptool

udptool is used to test the network connectivity between two points, as well as the network quality (bandwidth, rtt, packet loss rate, etc.).

udptool is developed in golang, cross-platform, and high-performance.

# How to use
## udpserver
``
./udpserver -s 192.168.1.100 -p 9797
``
* -s 192.168.1.100: the server side's udp listen ip
* -p 9797: the server side's udp listen port

## udpclient
./udpclient -c 10.0.0.1 -d 30300 -p 9797 -s 192.168.1.100
* -c 10.0.0.1: client side's udp listen ip
* -d 300300: client side's udp listen port
* -p 9797: client side's udp listen port
* -s 192.168.1.100: client side's udp listen ip
