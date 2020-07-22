package types

import "fmt"

type GenesisState struct {
	PDVRecords []PDV `json:"pdv_records"`
}

func NewGenesisState(pdvRecords []PDV) GenesisState {
	return GenesisState{PDVRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.PDVRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
		}
		if record.Value == "" {
			return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		PDVRecords: []PDV{},
	}
}
