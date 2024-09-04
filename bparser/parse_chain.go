/*
package bparser is for parsing bitcoin blockchain data directly from .dat files
*/
package bparser

import (
	"slices"
	"strings"
)

/*
byteSwapArray function will take a 64 character hexidecimal string and swap the bytes and return as a slice.

This is for converting from LittleEndian to BigEndian and vice-versa.

# Example

	genesisBlock := "6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000"
	byteSwap(genesisBlock)

returns [00 00 00 00 00 19 D6 68 9C 08 5A E1 65 83 1E 93 4F F7 63 AE 46 A2 A6 C1 72 B3 F1 B6 0A 8C E2 6F]
*/
func ByteSwapArray(hash string) []string {
	var hashHex []string
	for i := 0; i < len(hash); i += 2 {
		hashHex = append(hashHex, hash[i:i+2])
	}
	slices.Reverse(hashHex)

	return hashHex
}

/*
byteSwap function will take a 64 character hexidecimal string and swap the bytes and return a string.

This is for converting from LittleEndian to BigEndian and vice-versa.

# Example

	genesisBlock := "6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000"
	byteSwap(genesisBlock)

returns "000000000019D6689C085AE165831E934FF763AE46A2A6C172B3F1B60A8CE26F"
*/
func ByteSwap(hash string) string {
	var hashHex []string
	for i := 0; i < len(hash); i += 2 {
		hashHex = append(hashHex, hash[i:i+2])
	}
	slices.Reverse(hashHex)

	return strings.ToUpper(strings.Join(hashHex, ""))
}
