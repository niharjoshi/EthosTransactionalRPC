# CS 587 Assignment 1

[x] Name: Nihar Shailesh Joshi
[x] NetID: njoshi27@uic.edu
[x] UIN: 677063712


## Instructions

- Boot up an Ethos instance in Oracle VirtualBox and log in
- Open a new terminal window and clone this repository onto the VM
```console
git clone https://github.com/niharjoshi/EthosTransactionalRPC.git
```
- Run the commands below to start the RPC server
```console
make install
cd server
sudo -E ethosRun
```
- Open a new terminal window and run the commands below to start the RPC client
```console
cd server
etAl server.ethos
cd /programs
myRpcClient
```
