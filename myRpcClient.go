package main

import (
	"ethos/altEthos"
	"ethos/defined"
	"ethos/kernelTypes"
	"ethos/syscall"
	"log"
)

func init() {
	SetupMyRpcCreateAccountReply(createAccountReply)
	SetupMyRpcGetBalanceReply(getBalanceReply)
	SetupMyRpcTransferMoneyReply(transferMoneyReply)
}

func createAccountReply(message string, status syscall.Status) MyRpcProcedure {
	if status == syscall.StatusOk {
		EthosSTDOUT(kernelTypes.String("Account created successfully\n"))
	} else {
		EthosSTDOUT(kernelTypes.String(message + "\n"))
	}
	return nil
}

func getBalanceReply(balance string, message string, status syscall.Status) MyRpcProcedure {
	if status == syscall.StatusOk {
		EthosSTDOUT(kernelTypes.String("Balance: " + balance + "\n"))
	} else {
		EthosSTDOUT(kernelTypes.String(message + "\n"))
	}
	return nil
}

func transferMoneyReply(sourceBalance string, destinationBalance string, message string, status syscall.Status) MyRpcProcedure {
	if status == syscall.StatusOk {
		EthosSTDOUT(kernelTypes.String("New source balance: " + sourceBalance + "\n"))
		EthosSTDOUT(kernelTypes.String("New destination balance: " + destinationBalance + "\n"))
		EthosSTDOUT(kernelTypes.String("Transfer successful\n"))
	} else {
		EthosSTDOUT(kernelTypes.String(message + "\n"))
	}
	return nil
}

func accountCreation(accountHolderUserName string, startingBalance string) {
	call := MyRpcCreateAccount{accountHolderUserName, startingBalance}
	ipcCall(&call)
}

func balanceCheck(accountHolderUserName string) {
	call := MyRpcGetBalance{accountHolderUserName}
	ipcCall(&call)
}

func moneyTransfer(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount string) {
	call := MyRpcTransferMoney{sourceAccountHolderUserName, destinationAccountHolderUserName, transferAmount}
	ipcCall(&call)
}

func ipcCall(call defined.Rpc) {
	fd, status := altEthos.IpcRepeat("myRpc", "", nil)
	if status != syscall.StatusOk {
		log.Printf("IP call failed: %v\n", status)
		altEthos.Exit(status)
	}
	status = altEthos.ClientCall(fd, call)
	if status != syscall.StatusOk {
		log.Printf("Client call failed: %v\n", status)
		altEthos.Exit(status)
	}
}

func EthosSTDIN() kernelTypes.String {
	var input kernelTypes.String
	status := altEthos.ReadStream(syscall.Stdin, &input)
	if status != syscall.StatusOk {
		log.Printf("Error while reading syscall.Stdin: %v\n", status)
	}
	return input
}

func EthosSTDOUT(prompt kernelTypes.String) {
	status := altEthos.WriteStream(syscall.Stdout, &prompt)
	if status != syscall.StatusOk {
		log.Printf("Error writing to syscall.Stdout: %v\n", status)
	}
}

func menu() {
	EthosSTDOUT("Options:\n")
	EthosSTDOUT("1. Create account\n")
	EthosSTDOUT("2. Get balance\n")
	EthosSTDOUT("3. Transfer money\n")
	var userInput kernelTypes.String
	status := altEthos.ReadStream(syscall.Stdin, &userInput)
	if status != syscall.StatusOk {
		log.Printf("Main menu - error while reading syscall.Stdin: %v\n", status)
	}
	handler(userInput)
}

func handler(userInput string) {
	switch userInput {
	case "1\n":
		EthosSTDOUT("Enter a username: ")
		var username = string(EthosSTDIN())
		EthosSTDOUT("Enter starting balance: ")
		var balance = string(EthosSTDIN())
		accountCreation(username, balance)
	case "2\n":
		EthosSTDOUT("Enter a username: ")
		var username = string(EthosSTDIN())
		balanceCheck(username)
	case "3\n":
		EthosSTDOUT("Enter the source account username: ")
		var source = string(EthosSTDIN())
		EthosSTDOUT("Enter the destination account username: ")
		var destination = string(EthosSTDIN())
		EthosSTDOUT("Enter transfer amount: ")
		var amount = string(EthosSTDIN())
		moneyTransfer(source, destination, amount)
	default:
		EthosSTDOUT("Wrong input, try again.\n")
	}
}

func main() {
	altEthos.LogToDirectory("/home/me/EthosTransactionalRPC/client/")
	log.Printf("Starting RPC client\n")
	menu()
	log.Printf("Shutting down RPC client\n")
}
