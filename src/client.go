package src

import (
	"ethos/altEthos"
	"log"
)

func init() {
	CreateAccountRPCReply(createAccountRPCReply)
	GetBalanceRPCReply(getBalanceRPCReply)
	TransferMoneyRPCReply(transferMoneyRPCReply)
}

func createAccountRPCReply(accountHolderName string, startingBalance float64) {

}

func getBalanceRPCReply(accountHolderName string) {

}

func transferMoneyRPCReply(sourceAccountHolderName string, destinationAccountHolderName string) {

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

		} else if input == "2\n" {

		} else if input == "3\n" {

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
