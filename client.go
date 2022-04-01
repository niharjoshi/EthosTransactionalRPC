package EthosTransactionalRPC

import (
	"ethos/altEthos"
	"ethos/defined"
	"ethos/syscall"
	"log"
)

func init() {
	SetupCreateAccountReply(createAccountReply)
	SetupGetBalanceReply(getBalanceReply)
	SetupTransferMoneyReply(transferMoneyReply)
}

func createAccountReply(status string) MyRpcProcedure {
	if status == "Account created successfully" {
		EthosSTDOUT("Account created successfully\n")
	} else {
		EthosSTDOUT("Unable to create account\n")
	}
	return nil
}

func getBalanceReply(balance string) MyRpcProcedure {
	EthosSTDOUT("Balance: " + balance + "\n")
	return nil
}

func transferMoneyReply(sourceBalance string, destinationBalance string) MyRpcProcedure {
	EthosSTDOUT("New source balance: " + sourceBalance + "\n")
	EthosSTDOUT("New destination balance: " + destinationBalance + "\n")
	EthosSTDOUT("Transfer successful\n")
	return nil
}

func accountCreation(accountHolderUserName string, startingBalance float64) {
	call := MyRpcCreateAccount{accountHolderUserName, startingBalance}
	ipcCall(&call)
}

func balanceCheck(accountHolderUserName string) {
	call := MyRpcGetBalance{accountHolderUserName}
	ipcCall(&call)
}

func moneyTransfer(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount float64) {
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

func menu() {
	for {
		EthosSTDOUT("Options:\n")
		EthosSTDOUT("1. Create account\n")
		EthosSTDOUT("2. Get balance\n")
		EthosSTDOUT("3. Transfer money\n")
		EthosSTDOUT("4. Exit\n")
		var input = EthosSTDIN()
		if input == "1\n" {
			EthosSTDOUT("Enter a username: ")
			var username = EthosSTDIN()
			EthosSTDOUT("Enter starting balance: ")
			var balance = EthosSTDIN()
			accountCreation(username, balance)
		} else if input == "2\n" {
			EthosSTDOUT("Enter a username: ")
			var username = EthosSTDIN()
			balanceCheck(username)
		} else if input == "3\n" {
			EthosSTDOUT("Enter the source account username: ")
			var source = EthosSTDIN()
			EthosSTDOUT("Enter the destination account username: ")
			var destination = EthosSTDIN()
			EthosSTDOUT("Enter transfer amount: ")
			var amount = EthosSTDIN()
			transferMoney(source, destination, amount)
		} else if input == "4\n" {
			break
		} else {
			EthosSTDOUT("Wrong input, try again.\n")
		}
	}
}

func main() {
	altEthos.LogToDirectory("log/client")
	log.Printf("Starting RPC client\n")
	menu()
	log.Printf("Shutting down RPC client\n")
}
