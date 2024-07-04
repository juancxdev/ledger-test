package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Transaction struct {
	SourceAccountID      string  `db:"source_account_id" json:"source_account_id"`
	DestinationAccountID string  `db:"destination_account_id" json:"destination_account_id"`
	Amount               float64 `db:"amount" json:"amount"`
	TransactionType      string  `db:"transaction_type" json:"transaction_type"`
	TransactionDate      string  `db:"transaction_date" json:"transaction_date"`
	Status               string  `db:"status" json:"status"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, transactionID, sourceAccountID, destinationAccountID string, amount float64, transactionType, transactionDate, status string) error {
	resTransaction, err := s.Query(ctx, transactionID)
	if resTransaction != nil {
		fmt.Printf("Transaction already exist error: %s", err.Error())
		return err
	}

	transaction := Transaction{
		SourceAccountID:      sourceAccountID,
		DestinationAccountID: destinationAccountID,
		Amount:               amount,
		TransactionType:      transactionType,
		TransactionDate:      transactionDate,
		Status:               status,
	}

	transactionAsBytes, err := json.Marshal(transaction)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	return ctx.GetStub().PutState(transactionID, transactionAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, transactionID string) (*Transaction, error) {
	transactionAsBytes, err := ctx.GetStub().GetState(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if transactionAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", transactionID)
	}

	transaction := new(Transaction)
	err = json.Unmarshal(transactionAsBytes, transaction)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error. %s", err.Error())
	}

	return transaction, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create Bjaguar Ledger chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Bjaguar Ledger chaincode: %s", err.Error())
	}
}
