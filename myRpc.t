MyRpc interface {
    CreateAccount(accountHolderUserName string, startingBalance string) (status string)
    GetBalance(accountHolderUserName string) (balance string)
    TransferMoney(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount string) (sourceBalance string, destinationBalance string)
}
