package kcoin

type TransactionOutput struct {
	Value        int
	ScriptPubKey string
}

func (out *TransactionOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
