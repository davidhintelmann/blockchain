/*
package bparser is for parsing bitcoin blockchain data directly from .dat files
*/
package bparser

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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

/*
ParseBlocks function will parse an entire .dat bitcoin-core file and output a text file.
*/
func ParseBlocks(blks []byte, block_height_start int, block_height_end int, input_remainder []byte) error {
	// blocks := make(map[int]BlockStructure)
	// var blockRemainder byte
	if input_remainder[0] != 0 {
		blks = append(input_remainder, blks...)
	}

	// split bytes on magic number 'f9beb4d9'
	blocks := bytes.Split(blks, []byte{249, 190, 180, 217})[1:]

	for _, b := range blocks {
		// fmt.Printf("parsing block number: %d\n", i)
		blk := append([]byte{249, 190, 180, 217}, b...)
		_, err := ParseBlock(blk)
		if err != nil {
			errMsg := fmt.Sprintf("could not parse block, error: %v\n", err)
			return errors.New(errMsg)
		}

		// fmt.Printf("Magic Number: %v\nBlock Size: %v\n", block.MagicNumber, block.Size)
		// blockTimeStamp := time.Unix(block.BlockHeader.TimestampUnix, 0)
		// fmt.Printf("Version: %v\nPrev Block: %v\nMerkle Root: %v\nTimestamp: %v\nBits: %v\nNonce: %v\n", block.BlockHeader.Version, block.BlockHeader.PrevBlock, block.BlockHeader.MerkleRoot, blockTimeStamp, block.BlockHeader.Bits, block.BlockHeader.Nonce)
	}
	fmt.Printf("parsed %d blocks\n", len(blocks))
	return nil
}

/*
ParseMagicNumber function slices the first 4 bytes and returns the little-endian magic number.
*/
func ParseMagicNumber(blks []byte) []byte {
	return blks[:4]
}

// structs for raw block bytes in little-endian format
type ParseBlockBytes struct {
	MagicNumber []byte
	Size        []byte
	BlockHeader BlockHeaderBytes
	Tx          BlockTransactionsBytes
}

type BlockHeaderBytes struct {
	Version    []byte // int
	PrevBlock  []byte
	MerkleRoot []byte
	Timestamp  []byte // time.Time
	Bits       []byte
	Nonce      []byte
}

type BlockTransactionsBytes struct {
	TxCount []byte
	TxId    []byte
}

// structs for block strings data in big-endian format
type ParseBlockString struct {
	MagicNumber string
	Size        string
	BlockHeader BlockHeaderString
	Tx          BlockTransactionsBytes
}

type BlockHeaderString struct {
	Version    string
	PrevBlock  string
	MerkleRoot string
	Timestamp  string
	Bits       string
	Nonce      string
}

// type BlockTransactionsString struct {
// 	TxCount string
// 	TxId    string
// }

// structs for block strings/int/time data in big-endian format
type BlockData struct {
	MagicNumber string
	Size        int64
	BlockHeader BlockHeaderData
	Tx          BlockTransactionsBytes
}

type BlockHeaderData struct {
	Version       int64
	PrevBlock     string
	MerkleRoot    string
	TimestampUnix int64
	Bits          string
	Nonce         int64
}

// type BlockTransactionData struct {
// 	TxCount string
// 	TxId    string
// }

type ParseBlockSizeBytes struct {
	Size []byte
}

func (p ParseBlockSizeBytes) Raw() []byte {
	return p.Size
}

func (p ParseBlockSizeBytes) ParseInt() (int64, error) {
	bs := fmt.Sprintf("%X", p.Size)
	// swap bytes
	blockSize2 := ByteSwap(bs)
	// convert string to int64
	bs2, err := strconv.ParseInt(blockSize2, 16, 64)
	if err != nil {
		errMsg := fmt.Sprintf("could not parse int, error: %v\n", err)
		return -1, errors.New(errMsg)
	}
	return bs2, nil
}

func ParseBlockSizeRaw(blks []byte) ([]byte, error) {
	if len(blks) >= 8 {
		parseBlockSize := ParseBlockSizeBytes{Size: blks[4:8]}
		return parseBlockSize.Raw(), nil
	} else {
		errMsg := errors.New("can not slice bytes in order to read the size of the next block, index out of bounds")
		return nil, errMsg
	}
}

/*
ParseBlockSize function slices the first 4 to 8 bytes and returns an int.
First the little-endian magic number is converted into big-endian, by swapping the bytes, and finally converted into an int.
*/
func ParseBlockSize(blks []byte) (int64, error) {
	if len(blks) >= 8 {
		parseBlockSize := ParseBlockSizeBytes{Size: blks[4:8]}
		return parseBlockSize.ParseInt()
	} else {
		errMsg := errors.New("can not slice bytes in order to read the size of the next block, index out of bounds")
		return -1, errMsg
	}
}

/*
ParseBlockSizeFunc function slices the first 4 to 8 bytes and returns an int.
First the little-endian magic number is converted into big-endian, by swapping the bytes, and finally converted into an int.
*/
func ParseBlockSizeFunc(blks []byte) (int64, error) {
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

/*
ParseBlockRaw function will parse a single block at a time and return byte slices of little-endian numbers.
*/
func ParseBlockRaw(blks []byte) (ParseBlockBytes, error) {
	if len(blks) >= 8 {
		blockSize, err := ParseBlockSize(blks)
		if err != nil {
			return ParseBlockBytes{}, err
		}

		parseBlockHeader := BlockHeaderBytes{
			Version:    blks[8:12],
			PrevBlock:  blks[12:44],
			MerkleRoot: blks[44:76],
			Timestamp:  blks[76:80],
			Bits:       blks[80:84],
			Nonce:      blks[84:88],
		}

		parseBlockTransactions := BlockTransactionsBytes{
			TxCount: blks[88:90],
			TxId:    blks[90:blockSize],
		}

		parseBlock := ParseBlockBytes{
			MagicNumber: blks[:4],
			Size:        blks[4:8],
			BlockHeader: parseBlockHeader,
			Tx:          parseBlockTransactions,
		}
		return parseBlock, nil
	} else {
		errMsg := errors.New("can not slice bytes in order to read the size of the next block, index out of bounds")
		return ParseBlockBytes{}, errMsg
	}
}

/*
ParseBlockStr function will parse a single block at a time and return strings of big-endian numbers.
*/
func ParseBlockStr(blks []byte) (ParseBlockString, error) {
	if len(blks) >= 8 {
		blockSize, err := ParseBlockSize(blks)
		if err != nil {
			return ParseBlockString{}, err
		}

		parseBlockHeader := BlockHeaderString{
			Version:    ByteSwap(fmt.Sprintf("%X", blks[8:12])),
			PrevBlock:  ByteSwap(fmt.Sprintf("%X", blks[12:44])),
			MerkleRoot: ByteSwap(fmt.Sprintf("%X", blks[44:76])),
			Timestamp:  ByteSwap(fmt.Sprintf("%X", blks[76:80])),
			Bits:       ByteSwap(fmt.Sprintf("%X", blks[80:84])),
			Nonce:      ByteSwap(fmt.Sprintf("%X", blks[84:88])),
		}

		parseBlockTransactions := BlockTransactionsBytes{
			TxCount: blks[88:90],
			TxId:    blks[90:blockSize],
		}

		parseBlock := ParseBlockString{
			MagicNumber: ByteSwap(fmt.Sprintf("%X", blks[:4])),
			Size:        ByteSwap(fmt.Sprintf("%X", blks[4:8])),
			BlockHeader: parseBlockHeader,
			Tx:          parseBlockTransactions,
		}
		return parseBlock, nil
	} else {
		errMsg := errors.New("can not slice bytes in order to read the size of the next block, index out of bounds")
		return ParseBlockString{}, errMsg
	}
}

/*
ParseBlock function will parse a single block at a time and return strings or ints of big-endian numbers.
*/
func ParseBlock(blks []byte) (BlockData, error) {
	if len(blks) >= 8 {
		blockSize, err := ParseBlockSize(blks)
		if err != nil {
			return BlockData{}, err
		}

		s, err := strconv.ParseInt(ByteSwap(fmt.Sprintf("%X", blks[4:8])), 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap size bytes in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		v, err := strconv.ParseInt(ByteSwap(fmt.Sprintf("%X", blks[8:12])), 16, 16)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap version bytes in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		t, err := strconv.ParseInt(ByteSwap(fmt.Sprintf("%X", blks[76:80])), 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap time bytes in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		n, err := strconv.ParseInt(ByteSwap(fmt.Sprintf("%X", blks[84:88])), 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap nonce bytes in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		parseBlockHeader := BlockHeaderData{
			Version:       v,
			PrevBlock:     ByteSwap(fmt.Sprintf("%X", blks[12:44])),
			MerkleRoot:    ByteSwap(fmt.Sprintf("%X", blks[44:76])),
			TimestampUnix: t,
			Bits:          ByteSwap(fmt.Sprintf("%X", blks[80:84])),
			Nonce:         n,
		}

		parseBlockTransactions := BlockTransactionsBytes{
			TxCount: blks[88:90],
			TxId:    blks[90:blockSize],
		}

		parseBlock := BlockData{
			MagicNumber: ByteSwap(fmt.Sprintf("%X", blks[:4])),
			Size:        s,
			BlockHeader: parseBlockHeader,
			Tx:          parseBlockTransactions,
		}
		return parseBlock, nil
	} else {
		errMsg := errors.New("can not slice bytes in order to read the size of the next block, index out of bounds")
		return BlockData{}, errMsg
	}
}
