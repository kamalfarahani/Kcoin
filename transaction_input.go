package kcoin

type TransactionInput struct {
	TxID        []byte
	OutputIndex int
	ScriptSig   string
}

func (in *TransactionInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
