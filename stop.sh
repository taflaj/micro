#! /bin/sh

# Stop microservices
killall -USR1 server
sleep 1
killall -USR1 random
killall -USR1 pubkey
# killall -USR1 digest
killall -USR1 access
sleep 1
killall -USR1 router