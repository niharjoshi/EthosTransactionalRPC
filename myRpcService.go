package main

import (
	"ethos/altEthos"
	"ethos/kernelTypes"
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

func CustomHandleImport(eventInfo altEthos.ImportEventInfo) {
	event, status := altEthos.ReadRpcStreamAsync(eventInfo.ReturnedFd, eventInfo.I, altEthos.HandleRpc)
	if status != syscall.StatusOk {
		log.Println("RPC stream read failed")
		return
	}
	eventFd[event] = eventInfo.ReturnedFd
	altEthos.PostEvent(event)
	event, status = altEthos.ImportAsync(eventInfo.Fd, eventInfo.I, CustomHandleImport)
	if status != syscall.StatusOk {
		log.Println("Async import failed")
		return
	}
	altEthos.PostEvent(event)
}

func main() {
	var pathService = path
	var pathClient = "/home/me/EthosTransactionalRPC/client/"
	var pathDatastore = "/home/me/EthosTransactionalRPC/server/datastore/"
	var pathType kernelTypes.String
	var checkPathService = altEthos.IsDirectory(pathService)
	if checkPathService == false {
		log.Printf("Creating service logs directory\n")
		var status1 = altEthos.DirectoryCreate(pathService, &pathType, "all")
		if status1 != syscall.StatusOk {
			log.Println("Service logs directory create failed: ", pathService, status1)
			altEthos.Exit(status1)
		}
	}
	var checkPathClient = altEthos.IsDirectory(pathClient)
	if checkPathClient == false {
		log.Printf("Creating client logs directory\n")
		var status2 = altEthos.DirectoryCreate(pathClient, &pathType, "all")
		if status2 != syscall.StatusOk {
			log.Println("Client logs directory create failed: ", pathClient, status2)
			altEthos.Exit(status2)
		}
	}
	var checkPathDatastore = altEthos.IsDirectory(pathDatastore)
	if checkPathDatastore == false {
		log.Printf("Creating datastore directory\n")
		var status3 = altEthos.DirectoryCreate(pathDatastore, &pathType, "all")
		if status3 != syscall.StatusOk {
			log.Println("Datastore directory create failed: ", pathDatastore, status3)
			altEthos.Exit(status3)
		}
	}
	log.Printf("Starting RPC service\n")
	listeningFd, status := altEthos.Advertise("myRpc")
	if status != syscall.StatusOk {
		log.Printf("Advertising service failed: %s\n", status)
		altEthos.Exit(status)
	}
	var tree altEthos.EventTreeSlice
	var next []syscall.EventId
	t := MyRpc{}
	event, status := altEthos.ImportAsync(listeningFd, &t, CustomHandleImport)
	if status != syscall.StatusOk {
		log.Println("Import failed")
		return
	}
	next = append(next, event)
	tree = altEthos.WaitTreeCreateOr(next)
	for {
		tree, _ = altEthos.Block(tree)
		completed, pending := altEthos.GetTreeEvents(tree)
		for _, eventId := range completed {
			eventInfo, status := altEthos.OnComplete(eventId)
			if status != syscall.StatusOk {
				log.Println("OnComplete failed", eventInfo, status)
				return
			}
			eventInfo.Do()
		}
		next = nil
		next = append(next, pending...)
		next = append(next, altEthos.RetrievePostedEvents()...)
		tree = altEthos.WaitTreeCreateOr(next)
	}
	log.Printf("Shutting down RPC server\n")
}
