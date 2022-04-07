package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"log"
	"strconv"
)

var path = "/home/me/EthosTransactionalRPC/server/"
var eventFd = make(map[syscall.EventId]syscall.Fd)

var datastore = make(map[string]string)

func init() {
	SetupMyRpcCreateAccount(createAccount)
	SetupMyRpcGetBalance(getBalance)
	SetupMyRpcTransferMoney(transferMoney)
}

func createAccount(accountHolderUserName string, startingBalance string) MyRpcProcedure {
	value, ok := datastore[accountHolderUserName]
	if ok {
		log.Println("Account already exists")
		return &MyRpcCreateAccountReply{"Account already exists", syscall.StatusFail}
	} else {
		log.Println("Account already exists with balance " + value)
		datastore[accountHolderUserName] = startingBalance
	}
	return &MyRpcCreateAccountReply{"Account created successfully", syscall.StatusOk}
}

func getBalance(accountHolderUserName string) MyRpcProcedure {
	value, ok := datastore[accountHolderUserName]
	if ok {
		log.Println("Account has a balance of " + value)
	} else {
		log.Println("Account does not exist")
		return &MyRpcGetBalanceReply{"Null", "Account does not exist", syscall.StatusFail}
	}
	return &MyRpcGetBalanceReply{value, "Balance fetched successfully", syscall.StatusOk}
}

func transferMoney(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount string) MyRpcProcedure {
	valueSource, okSource := datastore[sourceAccountHolderUserName]
	if !okSource {
		log.Println("Source account does not exist")
		return &MyRpcTransferMoneyReply{"Null", "Null", "Unable to locate source account", syscall.StatusFail}
	}
	valueDestination, okDestination := datastore[destinationAccountHolderUserName]
	if !okDestination {
		log.Println("Destination account does not exist")
		return &MyRpcTransferMoneyReply{"Null", "Null", "Unable to locate destination account", syscall.StatusFail}
	}
	var valueSourceFloat int
	var valueDestinationFloat int
	valueSourceFloat, _ = strconv.Atoi(valueSource)
	valueDestinationFloat, _ = strconv.Atoi(valueDestination)
	var transferAmountFloat int
	transferAmountFloat, _ = strconv.Atoi(transferAmount)
	if valueSourceFloat < transferAmountFloat {
		log.Println("Not enough balance for transfer")
		return &MyRpcTransferMoneyReply{"Null", "Null", "Not enough balance for transfer", syscall.StatusFail}
	} else {
		valueSourceFloat -= transferAmountFloat
		valueDestinationFloat += transferAmountFloat
		valueSource = strconv.Itoa(valueSourceFloat)
		valueDestination = strconv.Itoa(valueDestinationFloat)
		datastore[sourceAccountHolderUserName] = valueSource
		datastore[destinationAccountHolderUserName] = valueDestination
		log.Println("New source balance: " + valueSource)
		log.Println("New destination balance: " + valueDestination)
	}
	return &MyRpcTransferMoneyReply{valueSource, valueDestination, "Transfer completed successfully", syscall.StatusOk}
}

func main() {
	statusLogs := altEthos.LogToDirectory("/home/me/EthosTransactionalRPC/server/")
	if statusLogs != syscall.StatusOk {
		log.Printf("Service logs directory create failed: %v\n", statusLogs)
		altEthos.Exit(statusLogs)
	}
	listeningFd, status := altEthos.Advertise("myRpc")
	if status != syscall.StatusOk {
		log.Printf("Advertising service failed: %s \n", status)
		altEthos.Exit(status)
	}
	for {
		_, fd, statusEvent := altEthos.Import(listeningFd)
		if statusEvent != syscall.StatusOk {
			log.Printf("Error calling import: %v\n", statusEvent)
			altEthos.Exit(statusEvent)
		}
		log.Printf("Server: new client connection accepted\n")
		t := MyRpc{}
		altEthos.Handle(fd, &t)
	}
}
