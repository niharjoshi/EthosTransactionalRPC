package EthosTransactionalRPC

import (
	"log"
	"syscall"
)

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
