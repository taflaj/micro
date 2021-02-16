#! /bin/sh

# Stop microservices
killall server
sleep 1
killall random
killall pubkey
killall access
sleep 1
killall router