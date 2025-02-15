package wallet

import (
	"dgc/util"
	"fmt"
	"time"

	"dgc/types" // Import the new types package
)

// NewTransaction creates a new transaction with outputs
func NewTransaction(senderWallet *types.Wallet, recipient string, amount int) *types.Transaction {
	if amount > senderWallet.Balance {
		fmt.Printf("Amount: %d exceeds balance.\n", amount)
		return nil
	}

	outputs := []types.Output{
		{Amount: senderWallet.Balance - amount, Address: senderWallet.PublicKey},
		{Amount: amount, Address: recipient},
	}
	return transactionWithOutputs(senderWallet, outputs)
}

// UpdateTransaction updates the transaction with a new amount and recipient.
func UpdateTransaction(t *types.Transaction, senderWallet *types.Wallet, recipient string, amount int) *types.Transaction {
	senderOutput := FindOutput(t, senderWallet.PublicKey)

	if senderOutput == nil || amount > senderOutput.Amount {
		fmt.Printf("Amount: %d exceeds balance.\n", amount)
		return nil
	}

	senderOutput.Amount -= amount
	t.Outputs = append(t.Outputs, types.Output{Amount: amount, Address: recipient})
	signTransaction(t, senderWallet)

	return t
}

// FindOutput finds the output for a given address in the transaction's outputs.
func FindOutput(t *types.Transaction, address string) *types.Output {
	for i := range t.Outputs {
		if t.Outputs[i].Address == address {
			return &t.Outputs[i]
		}
	}
	return nil
}

// transactionWithOutputs creates a transaction with specified outputs and signs it.
func transactionWithOutputs(senderWallet *types.Wallet, outputs []types.Output) *types.Transaction {
	transaction := &types.Transaction{
		ID:      util.ID(),
		Outputs: outputs,
	}
	signTransaction(transaction, senderWallet)
	return transaction
}

// RewardTransaction creates a reward transaction for miners.
func RewardTransaction(minerWallet *types.Wallet, blockchainWallet *types.Wallet) *types.Transaction {
	outputs := []types.Output{
		{Amount: util.MINING_REWARD, Address: minerWallet.PublicKey},
	}
	return transactionWithOutputs(blockchainWallet, outputs)
}

// signTransaction signs the transaction with the sender's wallet information.
func signTransaction(transaction *types.Transaction, senderWallet *types.Wallet) {
	transaction.Input = &types.Input{
		Timestamp: time.Now().Unix(),
		Amount:    senderWallet.Balance,
		Address:   senderWallet.PublicKey,
		Signature: Sign(senderWallet, util.ChainUtilHash(transaction.Outputs)),
	}
}

// VerifyTransaction verifies the signature of a given transaction.
func VerifyTransaction(transaction *types.Transaction) bool {
	publicKey, err := util.HexToPublicKey(transaction.Input.Address)
	if err != nil {
		fmt.Printf("Error converting public key: %v\n", err)
		return false
	}

	return util.VerifySignature(
		publicKey,
		transaction.Input.Signature,
		util.ChainUtilHash(transaction.Outputs),
	)
}
