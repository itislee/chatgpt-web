#!/bin/sh
pkill -9 auth
sleep 1
nohup ./auth >> auth.log &
ps aux|grep auth
