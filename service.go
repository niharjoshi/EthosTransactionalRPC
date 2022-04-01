package main

import (
	"ethos/altEthos"
	"ethos/kernelTypes"
	"ethos/syscall"
	"log"
	"strconv"
)

var path = "/user/" + altEthos.GetUser() + "/server/"

func init() {
	SetupMyRpcCreateAccount(createAccount)
	SetupMyRpcGetBalance(getBalance)
	SetupMyRpcTransferMoney(transferMoney)
}

func createAccount(accountHolderUserName string, startingBalance string) MyRpcProcedure {
	fd, status := altEthos.DirectoryOpen(path + "datastore/")
	if status != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", status)
	}
	var varName = accountHolderUserName
	var value = kernelTypes.String(startingBalance)
	status = altEthos.WriteVar(fd, varName, &value)
	if status != syscall.StatusOk {
		log.Printf("Error writing to %v: %v\n", path+"datastore/"+varName, status)
	}
	return &MyRpcCreateAccountReply{"Account created successfully"}
}

func getBalance(accountHolderUserName string) MyRpcProcedure {
	_, status := altEthos.DirectoryOpen(path + "datastore/")
	if status != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", status)
	}
	var value kernelTypes.String
	status = altEthos.Read(path+"datastore/"+accountHolderUserName, &value)
	if status != syscall.StatusOk {
		log.Printf("Error reading %v: %v\n", path+"datastore/"+accountHolderUserName, status)
	}
	return &MyRpcGetBalanceReply{string(value)}
}

func transferMoney(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount string) MyRpcProcedure {
	fdSource, statusSource := altEthos.DirectoryOpen(path + "datastore/")
	if statusSource != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", statusSource)
	}
	var valueSource kernelTypes.String
	statusSource = altEthos.Read(path+"datastore/"+sourceAccountHolderUserName, &valueSource)
	if statusSource != syscall.StatusOk {
		log.Printf("Error reading source %v: %v\n", path+"datastore/"+sourceAccountHolderUserName, statusSource)
	}
	fdDestination, statusDestination := altEthos.DirectoryOpen(path + "datastore/")
	if statusDestination != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", statusDestination)
	}
	var valueDestination kernelTypes.String
	statusDestination = altEthos.Read(path+"datastore/"+destinationAccountHolderUserName, &valueDestination)
	if statusSource != syscall.StatusOk {
		log.Printf("Error reading destination %v: %v\n", path+"datastore/"+destinationAccountHolderUserName, statusDestination)
	}
	var valueSourceFloat int
	var valueDestinationFloat int
	valueSourceFloat, statusSource = strconv.Atoi(valueSource)
	if statusSource != nil {
		log.Printf("Error converting source %v: %v\n", path+"datastore/"+sourceAccountHolderUserName, statusSource)
	}
	valueDestinationFloat, statusDestination = strconv.Atoi(valueDestination)
	if statusSource != nil {
		log.Printf("Error converting destination %v: %v\n", path+"datastore/"+destinationAccountHolderUserName, statusDestination)
	}
	var transferAmountFloat int
	transferAmountFloat, _ = strconv.Atoi(transferAmount)
	if valueSourceFloat < transferAmountFloat {
		log.Printf("Not enough balance for transfer\n")
	} else {
		valueSourceFloat -= transferAmountFloat
		valueDestinationFloat += transferAmountFloat

		var sourceVarName = kernelTypes.String(sourceAccountHolderUserName)
		var sourceValue = kernelTypes.String(strconv.Itoa(valueSourceFloat))
		statusSource = altEthos.WriteVar(fdSource, sourceVarName, &sourceValue)
		if statusSource != syscall.StatusOk {
			log.Printf("Error writing to %v: %v\n", path+"datastore/"+sourceVarName, statusSource)
		}

		var destinationVarName = kernelTypes.String(destinationAccountHolderUserName)
		var destinationValue = kernelTypes.String(strconv.Itoa(valueDestinationFloat))
		statusDestination = altEthos.WriteVar(fdDestination, destinationVarName, &destinationValue)
		if statusDestination != syscall.StatusOk {
			log.Printf("Error writing to %v: %v\n", path+"datastore/"+destinationVarName, statusDestination)
		}

		log.Printf("New source balance: %d\n", valueSourceFloat)
		log.Printf("New destination balance: %d\n", valueDestinationFloat)
	}
	return &MyRpcTransferMoneyReply{strconv.Itoa(valueSourceFloat), strconv.Itoa(valueDestinationFloat)}
}
