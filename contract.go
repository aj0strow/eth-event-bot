package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"io/ioutil"
)

type Contract struct {
	ContractName string
	ABI          abi.ABI
	Address      string
}

func (c *Contract) GetEvent(topic common.Hash) (abi.Event, bool) {
	for _, event := range c.ABI.Events {
		if event.Id() == topic {
			return event, true
		}
	}
	return abi.Event{}, false
}

func ParseTruffleBuild(filePath string) (*Contract, error) {
	rawData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	contract := &Contract{}
	err = json.Unmarshal(rawData, contract)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

type Stringer interface {
	String() string
}

func PrintValue(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	if str, ok := value.(Stringer); ok {
		return str.String()
	}
	return fmt.Sprintf("%s", value)
}

func UnpackValues(event abi.Event, log types.Log) ([]interface{}, error) {
	values := []interface{}{}
	rawValues, err := event.Inputs.UnpackValues(log.Data)
	if err != nil {
		return nil, err
	}
	rawIndex := 0
	topicIndex := 1
	for _, input := range event.Inputs {
		if input.Indexed {
			rawValue, err := UnpackRaw(input, log.Topics[topicIndex].Bytes())
			if err != nil {
				return nil, err
			}
			values = append(values, rawValue)
			topicIndex++
		} else {
			values = append(values, rawValues[rawIndex])
			rawIndex++
		}
	}
	return values, nil
}

func UnpackRaw(arg abi.Argument, data []byte) (interface{}, error) {
	switch arg.Type.T {
	case abi.AddressTy:
		return common.BytesToAddress(data), nil
	default:
		return data, nil
	}
}
