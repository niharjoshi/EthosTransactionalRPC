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

func createAccountReply(status string) MyRpcProcedure {
	if status == "Account created successfully" {
		EthosSTDOUT(kernelTypes.String("Account created successfully\n"))
	} else {
		EthosSTDOUT(kernelTypes.String("Unable to create account\n"))
	}
	return nil
}

func getBalanceReply(balance string) MyRpcProcedure {
	EthosSTDOUT(kernelTypes.String("Balance: " + balance + "\n"))
	return nil
}

func transferMoneyReply(sourceBalance string, destinationBalance string) MyRpcProcedure {
	EthosSTDOUT(kernelTypes.String("New source balance: " + sourceBalance + "\n"))
	EthosSTDOUT(kernelTypes.String("New destination balance: " + destinationBalance + "\n"))
	EthosSTDOUT(kernelTypes.String("Transfer successful\n"))
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
	for {
		EthosSTDOUT("Options:\n")
		EthosSTDOUT("1. Create account\n")
		EthosSTDOUT("2. Get balance\n")
		EthosSTDOUT("3. Transfer money\n")
		EthosSTDOUT("4. Exit\n")
		var input = string(EthosSTDIN())
		if input == "1\n" {
			EthosSTDOUT("Enter a username: ")
			var username = string(EthosSTDIN())
			EthosSTDOUT("Enter starting balance: ")
			var balance = string(EthosSTDIN())
			accountCreation(username, balance)
		} else if input == "2\n" {
			EthosSTDOUT("Enter a username: ")
			var username = string(EthosSTDIN())
			balanceCheck(username)
		} else if input == "3\n" {
			EthosSTDOUT("Enter the source account username: ")
			var source = string(EthosSTDIN())
			EthosSTDOUT("Enter the destination account username: ")
			var destination = string(EthosSTDIN())
			EthosSTDOUT("Enter transfer amount: ")
			var amount = string(EthosSTDIN())
			moneyTransfer(source, destination, amount)
		} else if input == "4\n" {
			break
		} else {
			EthosSTDOUT("Wrong input, try again.\n")
		}
	}
}

func main() {
	altEthos.LogToDirectory("/home/me/EthosTransactionalRPC/client/")
	log.Printf("Starting RPC client\n")
	menu()
	log.Printf("Shutting down RPC client\n")
}
