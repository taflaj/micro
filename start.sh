#! /bin/sh

# Start microservices
cd $HOME/projects/go/src/github.com/taflaj/micro
go run github.com/taflaj/micro/pubsub |& tee -a pubsub.log &
sleep 1
# go run github.com/taflaj/micro/access access.db |& tee -a access.log &
# go run github.com/taflaj/micro/digest digest.db |& tee -a digest.log &
# go run github.com/taflaj/micro/pubkey pubkey.db |& tee -a pubkey.log &
go run github.com/taflaj/micro/random |& tee -a random.log &
sleep 1
go run github.com/taflaj/micro/server ad18d93d20f79d |& tee -a server.log &