MyRpc interface {
    CreateAccount(accountHolderUserName string, startingBalance float64) (status string)
    GetBalance(accountHolderUserName string) (balance string)
    transferMoney(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount float64) (sourceBalance string, destinationBalance string)
}