/*
package bparser is for parsing bitcoin blockchain data directly from .dat files
*/
package bparser

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/davidhintelmann/blockchain/bparser/ttmpl"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
ParseBlocks function will parse an entire .dat bitcoin-core file (input in the form of bytes), and output a text file.
*/
func ParseBlocks(blks []byte, block_height_start int, block_height_end int, input_remainder []byte) error {
	// blocks := make(map[int]BlockStructure)
	// var blockRemainder byte
	if input_remainder[0] != 0 {
		blks = append(input_remainder, blks...)
	}

	// split bytes on magic number 'f9beb4d9'
	blocks := bytes.Split(blks, []byte{249, 190, 180, 217})[1:]

	for i, b := range blocks {
		// fmt.Printf("parsing block number: %d\n", i)
		blk := append([]byte{249, 190, 180, 217}, b...)
		block, err := ParseBlock(blk)
		if err != nil {
			errMsg := fmt.Sprintf("could not parse block, error: %v\n", err)
			return errors.New(errMsg)
		}
		if i == 0 {
			ttmpl.PrintBlock(block)
			// fmt.Printf("Block Number: %d\nMagic Number: %v\nBlock Size  : %v\n", i, block.Magic, block.Size)
			// blockTimeStamp := time.Unix(block.Header.TimestampUnix, 0)
			// fmt.Printf("Version     : %v\nBlock Hash  : %v\nPrev Block  : %v\nMerkle Root : %v\nTimestamp   : %v\nBits        : %v\nNonce       : %v\n", block.Header.Version, block.Header.BlockHash, block.Header.PrevBlock, block.Header.MerkleRoot, blockTimeStamp, block.Header.Bits, block.Header.Nonce)
			// fmt.Printf("Transactions: %d\n\n", block.Tx.TxCount)
		} else if block.Tx.TxCount > 1 {
			fmt.Printf("Block Number: %d\nMagic Number: %v\nBlock Size  : %v\n", i, block.Magic, block.Size)
			blockTimeStamp := time.Unix(block.Header.TimestampUnix, 0)
			fmt.Printf("Version     : %v\nBlock Hash  : %v\nPrev Block  : %v\nMerkle Root : %v\nTimestamp   : %v\nBits        : %v\nNonce       : %v\n", block.Header.Version, block.Header.BlockHash, block.Header.PrevBlock, block.Header.MerkleRoot, blockTimeStamp, block.Header.Bits, block.Header.Nonce)
			fmt.Printf("Transactions: %d\n\n", block.Tx.TxCount)
			break
		}
	}
	p := message.NewPrinter(language.English)
	p.Printf("parsed %d blocks\n", len(blocks))
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
	Magic  string
	Size   int64
	Header BlockHeaderData
	Tx     BlockTransactionsData
}

type BlockHeaderData struct {
	Version       int64
	BlockHash     string
	PrevBlock     string
	MerkleRoot    string
	TimestampUnix int64
	Bits          string
	Nonce         int64
}

type BlockTransactionsData struct {
	TxCount int64
	TxId    []byte
}

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
ParseBlockSize function takes in 4 bytes and returns an int which is length of the incoming block, in number of bytes.
First the little-endian magic number is converted into big-endian, by swapping the bytes, and finally converted into an int.
*/
func ParseBlockSize(blkSize []byte) (int64, error) {
	if len(blkSize) == 4 {
		parseBlockSize := ParseBlockSizeBytes{Size: blkSize}
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
		blockSize, err := ParseBlockSizeFunc(blks)
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
		blockSize, err := ParseBlockSizeFunc(blks)
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
	if len(blks) >= 4 {
		blockSize, err := ParseBlockSize(blks[4:8])
		if err != nil {
			errMsg := fmt.Sprintf("can not parse block size in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		// input the slice of block bytes which are apart of block header
		parseBlockHeader, err := parseBlockHeader(blks[8:88])
		if err != nil {
			errMsg := fmt.Sprintf("can not parse header in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		// input the slice of block bytes which are apart of block transactions
		parseBlockTransactions, err := parseBlockTransactions(blks[88 : blockSize+8])
		if err != nil {
			errMsg := fmt.Sprintf("can not parse transactions in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		parseBlock := BlockData{
			Magic:  ByteSwap(fmt.Sprintf("%X", blks[:4])),
			Size:   blockSize,
			Header: parseBlockHeader,
			Tx:     parseBlockTransactions,
		}
		return parseBlock, nil
	} else {
		errMsg := errors.New("can not slice bytes to read the size of the next block, index out of bounds. block being parsed is too small")
		return BlockData{}, errMsg
	}
}

/*
parseBlockHeader function is used in ParseBlock function to parse the header block in dat file from bitcoin-core.
*/
func parseBlockHeader(blkHeader []byte) (BlockHeaderData, error) {
	v, err := strconv.ParseInt(ByteSwap(fmt.Sprintf("%X", blkHeader[:4])), 16, 16)
	if err != nil {
		errMsg := fmt.Sprintf("can not swap version bytes in ParseBlock() function.\nerror: %v\n", err)
		return BlockHeaderData{}, errors.New(errMsg)
	}

	t, err := strconv.ParseInt(ByteSwap(fmt.Sprintf("%X", blkHeader[68:72])), 16, 64)
	if err != nil {
		errMsg := fmt.Sprintf("can not swap time bytes in ParseBlock() function.\nerror: %v\n", err)
		return BlockHeaderData{}, errors.New(errMsg)
	}

	n, err := strconv.ParseInt(ByteSwap(fmt.Sprintf("%X", blkHeader[76:80])), 16, 64)
	if err != nil {
		errMsg := fmt.Sprintf("can not swap nonce bytes in ParseBlock() function.\nerror: %v\n", err)
		return BlockHeaderData{}, errors.New(errMsg)
	}

	blockHash := sha256.New()
	blockHash.Write(blkHeader)
	sha256_single := blockHash.Sum(nil)
	blockHash = sha256.New()
	blockHash.Write(sha256_single)
	sha256_double := blockHash.Sum(nil)
	// fmt.Printf("Block Header: %x\n", blkHeader)
	// fmt.Printf("Block SHA256 single: %x\n", sha256_single)
	// fmt.Printf("Block SHA256 double: %x\n", sha256_double)
	slices.Reverse(sha256_double)
	// fmt.Printf("Block SHA256 swap  : %x\n", sha256_double)
	prevBlock := blkHeader[4:36]
	slices.Reverse(prevBlock)

	blockHeaderData := BlockHeaderData{
		Version:       v,
		BlockHash:     fmt.Sprintf("%X", sha256_double),
		PrevBlock:     fmt.Sprintf("%X", prevBlock),
		MerkleRoot:    ByteSwap(fmt.Sprintf("%X", blkHeader[36:68])),
		TimestampUnix: t,
		Bits:          ByteSwap(fmt.Sprintf("%X", blkHeader[72:76])),
		Nonce:         n,
	}

	return blockHeaderData, nil
}

/*
parseBlockTransactions function is used in ParseBlock function to parse the transaction block in dat file from bitcoin-core.
*/
func parseBlockTransactions(blkTransactions []byte) (BlockTransactionsData, error) {
	// size of incoming transaction block as hexidecimal
	// for more info read https://learnmeabitcoin.com/technical/general/compact-size/
	// transactionSizeHex := fmt.Sprintf("%x", 255)
	txCount, err := ParseTransactionBlockSize(blkTransactions[:9])
	if err != nil || txCount < 0 {
		errMsg := fmt.Sprintf("can not parse transaction block size in parseBlockTransactions() function.\nerror: %v\n", err)
		return BlockTransactionsData{}, errors.New(errMsg)
	}

	blockTransactionsData := BlockTransactionsData{
		TxCount: txCount,
		TxId:    blkTransactions[1:],
	}

	// fmt.Printf("TRANSACTIONS: %v\n", txCount)

	// version := blockTransactionsData.TxId[:4]
	// marker := blockTransactionsData.TxId[4:5]
	// flag := blockTransactionsData.TxId[5:6]

	// variable size for inputCount, learn more at https://learnmeabitcoin.com/technical/general/compact-size/
	// inputCount := blockTransactionsData.TxId[6:7]
	// inputs :=

	// variable size for outputCount, learn more at https://learnmeabitcoin.com/technical/general/compact-size/
	// outputCount := blockTransactionsData.TxId[6:7]

	// fmt.Printf("TV: %v\n", fmt.Sprintf("%X", version))

	return blockTransactionsData, nil
}

/*
parseTransactionBlockSize function is used in parseBlockTransactions function to parse the tx count in the transaction block.
*/
func ParseTransactionBlockSize(blkTranSize []byte) (int64, error) {
	leadingByte := strings.ToUpper(fmt.Sprintf("%x", blkTranSize[0]))
	// fmt.Printf("leading byte: %v\n", leadingByte)
	if leadingByte <= "FC" && leadingByte > "0" {
		txCount, err := strconv.ParseInt(leadingByte, 16, 16)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return -1, errors.New(errMsg)
		}
		return txCount, nil
	} else if leadingByte == "FD" {
		number := ByteSwap(fmt.Sprintf("%x", blkTranSize[1:3]))
		txCount, err := strconv.ParseInt(number, 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return -1, errors.New(errMsg)
		}
		return txCount, nil
	} else if leadingByte == "FE" {
		number := ByteSwap(fmt.Sprintf("%x", blkTranSize[1:5]))
		txCount, err := strconv.ParseInt(number, 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return -1, errors.New(errMsg)
		}
		return txCount, nil
	} else if leadingByte == "FF" {
		number := ByteSwap(fmt.Sprintf("%x", blkTranSize[1:9]))
		txCount, err := strconv.ParseInt(number, 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return -1, errors.New(errMsg)
		}
		return txCount, nil
	} else {
		errMsg := fmt.Sprintln("did not expect to return outside of if statement in parseTransactionBlockSize() function. leading byte does not match anything in table from https://learnmeabitcoin.com/technical/general/compact-size/#structure")
		return int64(-1), errors.New(errMsg)
	}
}
