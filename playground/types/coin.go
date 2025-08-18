package types

import (
	"fmt"
	"regexp"
	"strconv"
)

// Coin is akin to the Cosmos SDK's implementation by
// combining a given amount of a denomination.
type Coin struct {
	amount uint64
	denom  string
}

// coinPattern defines the expected regular expression pattern, that matches the coin
// representation as it is defined in the Cosmos SDK v0.50.13.
var coinPattern = regexp.MustCompile(`^([0-9]+)([a-zA-Z][a-zA-Z0-9/:._-]{2,127})$`)

func ParseCoin(s string) (Coin, error) {
	matches := coinPattern.FindStringSubmatch(s)
	if len(matches) != 3 {
		return Coin{}, fmt.Errorf("invalid coin: %s", s)
	}

	amount, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return Coin{}, fmt.Errorf("invalid amount: %s", matches[1])
	}

	denom := matches[2]

	return Coin{amount, denom}, nil
}

func (c Coin) Amount() uint64 { return c.amount }
func (c Coin) Denom() string  { return c.denom }
func (c Coin) String() string { return fmt.Sprintf("%d%s", c.amount, c.denom) }
