package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCoin(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected Coin
		expErr   string
	}{
		{
			name:   "error - empty string",
			input:  "",
			expErr: "invalid coin",
		},
		{
			name:   "error - invalid string",
			input:  "abcde1234",
			expErr: "invalid coin",
		},
		{
			name:     "success - valid coin",
			input:    "1234abcde",
			expected: Coin{amount: 1234, denom: "abcde"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			coin, err := ParseCoin(tc.input)
			if tc.expErr != "" {
				require.ErrorContains(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, coin)
			}
		})
	}
}
