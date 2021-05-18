#! /bin/sh

# Start microservices
cd $HOME/projects/go/src/github.com/taflaj/services
go run github.com/taflaj/services/router |& tee -a router.log &
sleep 1
go run github.com/taflaj/services/access access.db |& tee -a access.log &
# go run github.com/taflaj/services/digest digest.db |& tee -a digest.log &
go run github.com/taflaj/services/pubkey pubkey.db |& tee -a pubkey.log &
go run github.com/taflaj/services/random |& tee -a random.log &
sleep 1
go run github.com/taflaj/services/server ad18d93d20f79d |& tee -a server.log &