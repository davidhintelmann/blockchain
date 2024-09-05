/*
package bparser is for parsing bitcoin blockchain data directly from .dat files
*/
package bparser

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

/*
ByteSwapArray function will take a 64 character hexidecimal string and swap the bytes and return as a slice.

This is for converting from LittleEndian to BigEndian and vice-versa.

# Example

	genesisBlock := "6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000"
	ByteSwapArray(genesisBlock)

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
ByteSwap function will take a 64 character hexidecimal string and swap the bytes and return a string (all upper case).

This is for converting from LittleEndian to BigEndian and vice-versa.

# Example

	genesisBlock := "6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000"
	ByteSwap(genesisBlock)

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

/*
GlobDat function takes a single parameter for blocks data directory, and returns path of all dat files.
Unless there is an error.

# Example

	GlobDat("C:\\Users\\david\\OneDrive\\Documents\\code\\python\\Blockchain\\Bitcoin\\data\\bitcoin_data\\")

returns [C:\Users\david\OneDrive\Documents\code\python\Blockchain\Bitcoin\data\bitcoin_data\blk00000.dat C:\Users\david\OneDrive\Documents\code\python\Blockchain\Bitcoin\data\bitcoin_data\blk00001.dat C:\Users\david\OneDrive\Documents\code\python\Blockchain\Bitcoin\data\bitcoin_data\blk00002.dat, ...]
*/
func GlobDat(blocksFilePath string) ([]string, error) {
	matches, err := filepath.Glob(blocksFilePath + "*.dat")
	if err != nil {
		return []string{}, errors.New(err.Error())
	}

	return matches, nil
}

type BlockStructure struct {
	Version    int
	PrevBlock  []byte
	MerkleRoot []byte
	Timestamp  time.Time
	Bits       []byte
	Nonce      []byte
	TxIdsRaw   []byte
}

func ParseBlocks(blks []byte, block_height_start int, block_height_end int, input_remainder []byte) {
	// blocks := make(map[int]BlockStructure)
	// var blockRemainder byte
	if input_remainder[0] != 0 {
		blks = append(input_remainder, blks...)
	}

	fmt.Println(blks[:295])
}

/*
ParseMagicNumber function slices the first 4 bytes and returns the little-endian magic number.
*/
func ParseMagicNumber(blks []byte) []byte {
	return blks[:4]
}

func ParseBlockSizeRaw(blks []byte) []byte {
	return blks[4:8]
}

/*
ParseBlockSize function slices the first 4 to 8 bytes and returns an int.
First the little-endian magic number is converted into big-endian, by swapping the bytes, and finally converted into an int.
*/
func ParseBlockSize(blks []byte) (int64, error) {
	blockSize := blks[4:8]
	bs := fmt.Sprintf("%X", blockSize)
	// swap bytes
	blockSize2 := ByteSwap(bs)
	// convert string to int64
	bs2, err := strconv.ParseInt(blockSize2, 16, 16)
	if err != nil {
		errMsg := fmt.Sprintf("could not parse int, error: %v\n", err)
		return -1, errors.New(errMsg)
	}
	return bs2, nil
}
