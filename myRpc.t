MyRpc interface {
    CreateAccount(accountHolderUserName string, startingBalance string) (message string, message string, status Status)
    GetBalance(accountHolderUserName string) (balance string, message string, status Status)
    TransferMoney(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount string) (sourceBalance string, destinationBalance string, message string, status Status)
}
