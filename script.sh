#! /bin/sh
while true; do kill -9 $(ps aux | grep "[p]tpcamera" | awk '{print $2}'); done
