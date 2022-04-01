package main

import (
	"ethos/altEthos"
	"ethos/kernelTypes"
	"ethos/syscall"
	"fmt"
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
	var varName = kernelTypes.String(accountHolderUserName)
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
	var valueSourceFloat float64
	var valueDestinationFloat float64
	valueSourceFloat, statusSource = strconv.ParseFloat(valueSource, 32)
	if statusSource != nil {
		log.Printf("Error converting source %v: %v\n", path+"datastore/"+sourceAccountHolderUserName, statusSource)
	}
	valueDestinationFloat, statusDestination = strconv.ParseFloat(valueDestination, 32)
	if statusSource != nil {
		log.Printf("Error converting destination %v: %v\n", path+"datastore/"+destinationAccountHolderUserName, statusDestination)
	}
	var transferAmountFloat float64
	transferAmountFloat, _ = strconv.ParseFloat(transferAmount, 32)
	if valueSourceFloat < transferAmountFloat {
		log.Printf("Not enough balance for transfer\n")
	} else {
		valueSourceFloat -= transferAmountFloat
		valueDestinationFloat += transferAmountFloat

		var sourceVarName = kernelTypes.String(sourceAccountHolderUserName)
		var sourceValue = kernelTypes.String(fmt.Sprintf("%f", valueSourceFloat))
		statusSource = altEthos.WriteVar(fdSource, sourceVarName, &sourceValue)
		if statusSource != syscall.StatusOk {
			log.Printf("Error writing to %v: %v\n", path+"datastore/"+sourceVarName, statusSource)
		}

		var destinationVarName = kernelTypes.String(destinationAccountHolderUserName)
		var destinationValue = kernelTypes.String(fmt.Sprintf("%f", valueDestinationFloat))
		statusDestination = altEthos.WriteVar(fdDestination, destinationVarName, &destinationValue)
		if statusDestination != syscall.StatusOk {
			log.Printf("Error writing to %v: %v\n", path+"datastore/"+destinationVarName, statusDestination)
		}

		log.Printf("New source balance: %f\n", valueSourceFloat)
		log.Printf("New destination balance: %f\n", valueDestinationFloat)
	}
	return &MyRpcTransferMoneyReply{fmt.Sprintf("%f", valueSourceFloat), fmt.Sprintf("%f", valueDestinationFloat)}
}
