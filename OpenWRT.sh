#!/bin/sh /etc/rc.common
#start script for OpenWRT
START=95
start() {
    /root/uestc/connect /root/uestc/settings.ini /root/uestc/connect.pid 1>/root/uestc/uestc.log 2>&1 &
}

stop(){
    kill `cat /root/uestc/connect.pid`
    rm /root/uestc/connect.pid
}