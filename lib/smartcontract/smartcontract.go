package smartcontract

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ABIPackRaw returns []byte instead of string.
func ABIPackRaw(abiBytes []byte, method string, args ...any) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return []byte{}, fmt.Errorf("failed to parse the ABI: %s", err.Error())
	}

	callData, err := parsedABI.Pack(method, args...)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to pack the ABI: %s", err.Error())
	}

	return callData, nil
}

func ABIPack(abiBytes []byte, method string, args ...any) (string, error) {
	callData, err := ABIPackRaw(abiBytes, method, args...)
	if err != nil {
		return "", err
	}

	return "0x" + hex.EncodeToString(callData), nil
}

func StringsToABIArguments(args []string) ([]any, error) {
	callArgs := []any{}
	for _, v := range args {
		value := strings.Split(v, ":")
		switch value[0] {
		case "a":
			// Address
			callArgs = append(callArgs, common.HexToAddress(value[1]))
		case "n":
			// Numbers
			num := new(big.Int)
			_, valid := num.SetString(value[1], 10)
			if !valid {
				return []any{}, fmt.Errorf("error converting the number")
			}
			callArgs = append(callArgs, num)
		default:
			return callArgs, fmt.Errorf("invalid param type")
		}
	}

	return callArgs, nil
}
