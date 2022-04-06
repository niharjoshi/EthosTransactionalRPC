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

func init() {
	SetupMyRpcCreateAccount(createAccount)
	SetupMyRpcGetBalance(getBalance)
	SetupMyRpcTransferMoney(transferMoney)
}

func createAccount(accountHolderUserName string, startingBalance string) MyRpcProcedure {
	fd, status := altEthos.DirectoryOpen(path + "datastore/")
	if status != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", status)
		return &MyRpcCreateAccountReply{"Unable to open datastore", syscall.StatusFail}
	}
	var varName = accountHolderUserName
	var value = kernelTypes.String(startingBalance)
	status = altEthos.WriteVar(fd, varName, &value)
	if status != syscall.StatusOk {
		log.Printf("Error writing to %v: %v\n", path+"datastore/"+varName, status)
		return &MyRpcCreateAccountReply{"Unable to create account", syscall.StatusFail}
	}
	return &MyRpcCreateAccountReply{"Account created successfully", syscall.StatusOk}
}

func getBalance(accountHolderUserName string) MyRpcProcedure {
	_, status := altEthos.DirectoryOpen(path + "datastore/")
	if status != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", status)
		return &MyRpcCreateAccountReply{"Null", "Unable to open datastore", syscall.StatusFail}
	}
	var value kernelTypes.String
	status = altEthos.Read(path+"datastore/"+accountHolderUserName, &value)
	if status != syscall.StatusOk {
		log.Printf("Error reading %v: %v\n", path+"datastore/"+accountHolderUserName, status)
		return &MyRpcGetBalanceReply{"Null", "Unable to locate account", syscall.StatusFail}
	}
	return &MyRpcGetBalanceReply{string(value), "Balance fetched successfully", syscall.StatusOk}
}

func transferMoney(sourceAccountHolderUserName string, destinationAccountHolderUserName string, transferAmount string) MyRpcProcedure {
	fdSource, statusSource := altEthos.DirectoryOpen(path + "datastore/")
	if statusSource != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", statusSource)
		return &MyRpcCreateAccountReply{"Null", "Null", "Unable to open datastore", syscall.StatusFail}
	}
	var valueSource kernelTypes.String
	statusSource = altEthos.Read(path+"datastore/"+sourceAccountHolderUserName, &valueSource)
	if statusSource != syscall.StatusOk {
		log.Printf("Error reading source %v: %v\n", path+"datastore/"+sourceAccountHolderUserName, statusSource)
		return &MyRpcCreateAccountReply{"Null", "Null", "Unable to locate source account", syscall.StatusFail}
	}
	fdDestination, statusDestination := altEthos.DirectoryOpen(path + "datastore/")
	if statusDestination != syscall.StatusOk {
		log.Printf("Error fetching %v: %v\n", path+"datastore/", statusDestination)
		return &MyRpcCreateAccountReply{"Null", "Null", "Unable to open datastore", syscall.StatusFail}
	}
	var valueDestination kernelTypes.String
	statusDestination = altEthos.Read(path+"datastore/"+destinationAccountHolderUserName, &valueDestination)
	if statusSource != syscall.StatusOk {
		log.Printf("Error reading destination %v: %v\n", path+"datastore/"+destinationAccountHolderUserName, statusDestination)
		return &MyRpcCreateAccountReply{"Null", "Null", "Unable to locate destination account", syscall.StatusFail}
	}
	var valueSourceFloat int
	var valueDestinationFloat int
	valueSourceFloat, _ = strconv.Atoi(string(valueSource))
	valueDestinationFloat, _ = strconv.Atoi(string(valueDestination))
	var transferAmountFloat int
	transferAmountFloat, _ = strconv.Atoi(transferAmount)
	if valueSourceFloat < transferAmountFloat {
		log.Printf("Not enough balance for transfer\n")
		return &MyRpcCreateAccountReply{"Null", "Null", "Not enough balance for successful transfer", syscall.StatusFail}
	} else {
		valueSourceFloat -= transferAmountFloat
		valueDestinationFloat += transferAmountFloat

		var sourceVarName = kernelTypes.String(sourceAccountHolderUserName)
		var sourceValue = kernelTypes.String(strconv.Itoa(valueSourceFloat))
		statusSource = altEthos.WriteVar(fdSource, string(sourceVarName), &sourceValue)
		if statusSource != syscall.StatusOk {
			log.Printf("Error writing to %v: %v\n", path+"datastore/"+string(sourceVarName), statusSource)
			return &MyRpcCreateAccountReply{"Null", "Null", "Unable to transfer money from source account", syscall.StatusFail}
		}
		var destinationVarName = kernelTypes.String(destinationAccountHolderUserName)
		var destinationValue = kernelTypes.String(strconv.Itoa(valueDestinationFloat))
		statusDestination = altEthos.WriteVar(fdDestination, string(destinationVarName), &destinationValue)
		if statusDestination != syscall.StatusOk {
			log.Printf("Error writing to %v: %v\n", path+"datastore/"+string(destinationVarName), statusDestination)
			return &MyRpcCreateAccountReply{"Null", "Null", "Unable to transfer money from destination account", syscall.StatusFail}
		}

		log.Printf("New source balance: %d\n", valueSourceFloat)
		log.Printf("New destination balance: %d\n", valueDestinationFloat)
	}
	return &MyRpcTransferMoneyReply{strconv.Itoa(valueSourceFloat), strconv.Itoa(valueDestinationFloat), "Transfer completed successfully", syscall.StatusOk}
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
