中文 | [English](README-cn.md)
# udptool
udptool是用于测试两点之间网络连通情况，以及网络质量(带宽，rtt，丢包率等)。

udptool采用golang开发，跨平台，高性能。

# 如何使用
## udpserver
``
./udpserver -s 192.168.1.100 -p 9797
``
* -s 192.168.1.100 是服务端udp listen ip
* -p 9797 是服务端udp listen port

## udpclient
./udpclient -c 10.0.0.1 -d 30300 -p 9797 -s 192.168.1.100
* -c 10.0.0.1 是客户端udp listen ip
* -d 300300 是客户端udp listen port
* -p 9797 是服务端udp listen port
* -s 192.168.1.100 是服务端udp listen ip
