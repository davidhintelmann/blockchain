/*
package bparser is for parsing bitcoin blockchain data directly from .dat files
*/
package bparser

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// const (
// 	segWitHeight = 481_824
// )

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
ByteSwapStr function will take a 64 character hexidecimal string and swap the bytes and return a string (all upper case).

This is for converting from LittleEndian to BigEndian and vice-versa.

# Example

	genesisBlock := "6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000"
	ByteSwap(genesisBlock)

returns "000000000019D6689C085AE165831E934FF763AE46A2A6C172B3F1B60A8CE26F"
*/
func ByteSwapStr(hash string) string {
	var hashHex []string
	for i := 0; i < len(hash); i += 2 {
		hashHex = append(hashHex, hash[i:i+2])
	}
	slices.Reverse(hashHex)

	return strings.ToUpper(strings.Join(hashHex, ""))
}

/*
Insert Byte Swap Func
*/
func ByteSwap(hash []byte) string {
	hashString := fmt.Sprintf("%X", hash)
	var hashHex []string
	for i := 0; i < len(hashString); i += 2 {
		hashHex = append(hashHex, hashString[i:i+2])
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
func ParseBlocks(blks []byte, block_height_start int, block_height_end int, input_remainder []byte) (int, error) {
	if input_remainder[0] != 0 {
		blks = append(input_remainder, blks...)
	}

	// split bytes on magic number 'f9beb4d9'
	blocks := bytes.Split(blks, []byte{249, 190, 180, 217})[1:]

	// var block BlockData
	for i, b := range blocks {
		// fmt.Printf("parsing block number: %d\n", i)
		blk := append([]byte{249, 190, 180, 217}, b...)

		// parse block
		block, err := ParseBlock(blk, i)
		if err != nil {
			errMsg := fmt.Sprintf("could not parse block, error: %v\n", err)
			return -1, errors.New(errMsg)
		}

		if i < block_height_start {
			continue
		} else if i >= block_height_end {
			return i, nil
		} else if i == 0 {
			// must be run from main.go
			printBlock(block, "../bparser/block.tmpl")
		} else if i == 95414 { // else if block.Tx.TxCount > 2
			// must be run from main.go
			printBlock(block, "../bparser/block.tmpl")
			fmt.Println()
			return i, nil
		}
	}

	return -1, errors.New("did no expect to exit loop, ParseBlocks() function")
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
	Version    []byte
	PrevBlock  []byte
	MerkleRoot []byte
	Timestamp  []byte
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

// structs for block strings/int/time data in big-endian format
type BlockData struct {
	BlockNumber int
	Magic       string
	Size        int64
	Header      BlockHeaderData
	Tx          BlockTransactionsData
}

type BlockHeaderData struct {
	Version       int64
	BlockHash     string
	PrevBlock     string
	MerkleRoot    string
	TimestampUnix int64
	Timestamp     time.Time
	Bits          string
	Nonce         int64
}

type BlockTransactionsData struct {
	TxCount int64
	Tx      TxData
}

type TxData struct {
	Version     int64
	InputCount  int64
	Inputs      []TxInputs
	OutputCount int64
	Outputs     []TxOutputs
	Locktime    []byte
}

type TxInputs struct {
	TxId          string
	Vout          string
	ScriptSigSize int64
	ScriptSig     string
	Sequence      string
}

type TxOutputs struct {
	Amount           []byte
	ScriptPubKeySize int64
	ScriptPubKey     []byte
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
	blockSize2 := ByteSwapStr(bs)
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
	blockSize2 := ByteSwapStr(bs)
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
			Version:    ByteSwapStr(fmt.Sprintf("%X", blks[8:12])),
			PrevBlock:  ByteSwapStr(fmt.Sprintf("%X", blks[12:44])),
			MerkleRoot: ByteSwapStr(fmt.Sprintf("%X", blks[44:76])),
			Timestamp:  ByteSwapStr(fmt.Sprintf("%X", blks[76:80])),
			Bits:       ByteSwapStr(fmt.Sprintf("%X", blks[80:84])),
			Nonce:      ByteSwapStr(fmt.Sprintf("%X", blks[84:88])),
		}

		parseBlockTransactions := BlockTransactionsBytes{
			TxCount: blks[88:90],
			TxId:    blks[90:blockSize],
		}

		parseBlock := ParseBlockString{
			MagicNumber: ByteSwapStr(fmt.Sprintf("%X", blks[:4])),
			Size:        ByteSwapStr(fmt.Sprintf("%X", blks[4:8])),
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
func ParseBlock(blk []byte, blockNum int) (BlockData, error) {
	// var blockSize int64
	if len(blk) >= 4 {
		blockSize, err := ParseBlockSize(blk[4:8])
		if err != nil {
			errMsg := fmt.Sprintf("can not parse block size in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		} else if len(blk)-8 != int(blockSize) {
			errMsg := fmt.Sprintf("len of blk (%d) does not equal block size (%d).\nerror: %v\n", len(blk), int(blockSize), err)
			return BlockData{}, errors.New(errMsg)
		}

		// input the slice of block bytes which are apart of block header
		parseBlockHeader, err := parseBlockHeader(blk[8:88])
		if err != nil {
			errMsg := fmt.Sprintf("can not parse header in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		// input the slice of block bytes which are apart of block transactions
		parseBlockTransactions, err := parseBlockTransactions(blk[88:])
		if err != nil {
			errMsg := fmt.Sprintf("can not parse transactions in ParseBlock() function.\nerror: %v\n", err)
			return BlockData{}, errors.New(errMsg)
		}

		parseBlock := BlockData{
			BlockNumber: blockNum,
			Magic:       ByteSwapStr(fmt.Sprintf("%X", blk[:4])),
			Size:        blockSize,
			Header:      parseBlockHeader,
			Tx:          parseBlockTransactions,
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
	v, err := strconv.ParseInt(ByteSwapStr(fmt.Sprintf("%X", blkHeader[:4])), 16, 16)
	if err != nil {
		errMsg := fmt.Sprintf("can not swap version bytes in ParseBlock() function.\nerror: %v\n", err)
		return BlockHeaderData{}, errors.New(errMsg)
	}

	t, err := strconv.ParseInt(ByteSwapStr(fmt.Sprintf("%X", blkHeader[68:72])), 16, 64)
	if err != nil {
		errMsg := fmt.Sprintf("can not swap time bytes in ParseBlock() function.\nerror: %v\n", err)
		return BlockHeaderData{}, errors.New(errMsg)
	}

	n, err := strconv.ParseInt(ByteSwapStr(fmt.Sprintf("%X", blkHeader[76:80])), 16, 64)
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
	slices.Reverse(sha256_double)

	prevBlock := blkHeader[4:36]
	slices.Reverse(prevBlock)

	blockHeaderData := BlockHeaderData{
		Version:       v,
		BlockHash:     fmt.Sprintf("%X", sha256_double),
		PrevBlock:     fmt.Sprintf("%X", prevBlock),
		MerkleRoot:    ByteSwapStr(fmt.Sprintf("%X", blkHeader[36:68])),
		TimestampUnix: t,
		Timestamp:     time.Unix(t, 0),
		Bits:          ByteSwapStr(fmt.Sprintf("%X", blkHeader[72:76])),
		Nonce:         n,
	}

	return blockHeaderData, nil
}

/*
parseBlockTransactions function is used in ParseBlock function to parse the transaction block in dat file from bitcoin-core.
*/
func parseBlockTransactions(blkTransactions []byte) (BlockTransactionsData, error) {
	// size of incoming transaction block
	txCount, pad, err := ParseTransactionBlockSize(blkTransactions[:9])
	if err != nil || txCount < 0 {
		errMsg := fmt.Sprintf("can not parse transaction block size in parseBlockTransactions() function.\nerror: %v\n", err)
		return BlockTransactionsData{}, errors.New(errMsg)
	}

	txData, err := ParseBlockTx(blkTransactions, pad)
	if err != nil {
		errMsg := fmt.Sprintf("can not swap parse block tx in ParseBlock() function.\nerror: %v\n", err)
		return BlockTransactionsData{}, errors.New(errMsg)
	}

	blockTransactionsData := BlockTransactionsData{
		TxCount: txCount,
		Tx:      txData,
	}

	return blockTransactionsData, nil
}

/*
ParseBlockTx function is used in parseBlockTransactions function to parse the individual tx in transaction block.
*/
func ParseBlockTx(blkTransactions []byte, pad int) (TxData, error) {
	// parse version number for block transaction
	v, err := strconv.ParseInt(ByteSwapStr(fmt.Sprintf("%X", blkTransactions[pad:pad+4])), 16, 16)
	if err != nil {
		errMsg := fmt.Sprintf("can not swap version bytes in ParseBlock() function.\nerror: %v\n", err)
		return TxData{}, errors.New(errMsg)
	}

	// variable size for inputCount
	inputCount, txInputPad, err := ParseTransactionBlockSize(blkTransactions[pad+4 : pad+18])
	if err != nil || inputCount < 0 {
		errMsg := fmt.Sprintf("can not parse transaction input count in parseBlockTransactions() function, segwit transaction?\nerror: %v\n", err)
		return TxData{}, errors.New(errMsg)
	}

	var txInputs []TxInputs
	var blkPad int = pad + txInputPad + 4
	for i := 0; i < int(inputCount); i++ {
		blkTx := blkTransactions[blkPad:]
		txId := blkTx[:32]
		vOut := blkTx[32:36]
		// variable size
		scriptSigSize, scriptPad, err := ParseTransactionBlockSize(blkTx[36:46])
		if err != nil || inputCount < 0 {
			errMsg := fmt.Sprintf("can not parse transaction input script size in parseBlockTransactions() function, segwit transaction?\nerror: %v\n", err)
			return TxData{}, errors.New(errMsg)
		}

		scriptPad = scriptPad + 36
		scriptPadEnd := blkPad + scriptPad + 4
		scriptSig := blkTx[scriptPad : scriptPad+int(scriptSigSize)]
		sequence := blkTx[scriptPad+int(scriptSigSize) : scriptPad+4+int(scriptSigSize)]
		blkPad = scriptPadEnd + int(scriptSigSize)

		txInput := TxInputs{
			TxId:          fmt.Sprintf("%X", txId),
			Vout:          fmt.Sprintf("%X", vOut),
			ScriptSigSize: scriptSigSize,
			ScriptSig:     fmt.Sprintf("%X", scriptSig),
			Sequence:      fmt.Sprintf("%X", sequence),
		}
		txInputs = append(txInputs, txInput)
	}

	// variable size for outputCount
	outputCount, txOutputPad, err := ParseTransactionBlockSize(blkTransactions[blkPad : blkPad+10]) // blockTransactionsData.Tx[4:14]
	if err != nil || outputCount < 0 {
		errMsg := fmt.Sprintf("can not parse transaction output count in parseBlockTransactions() function, segwit transaction?\nerror: %v\n", err)
		return TxData{}, errors.New(errMsg)
	}

	var txOutputs []TxOutputs
	blkPadOutput := blkPad + txOutputPad
	for i := 0; i < int(outputCount); i++ {
		blkTx := blkTransactions[blkPadOutput:]
		amount := blkTx[:8]
		// variable size
		scriptPubKeySize, scriptPad, err := ParseTransactionBlockSize(blkTx[8:18]) // blockTransactionsData.Tx[4:14]
		if err != nil || outputCount < 0 {
			errMsg := fmt.Sprintf("can not parse transaction output script size in parseBlockTransactions() function, segwit transaction?\nerror: %v\n", err)
			return TxData{}, errors.New(errMsg)
		}
		scriptPubKey := blkTx[scriptPad+8 : scriptPad+8+int(scriptPubKeySize)]
		blkPadOutput = blkPadOutput + scriptPad + 8 + int(scriptPubKeySize)

		if int(scriptPubKeySize) != len(scriptPubKey) {
			errMsg := fmt.Sprintf("scriptPubKeySize does not equal length of scriptPubKey bytes in parseBlockTransactions() function\nerror: %v\n", err)
			return TxData{}, errors.New(errMsg)
		}

		txOutput := TxOutputs{
			Amount:           amount,
			ScriptPubKeySize: scriptPubKeySize,
			ScriptPubKey:     scriptPubKey,
		}
		txOutputs = append(txOutputs, txOutput)
	}

	txData := TxData{
		Version:     v,
		InputCount:  inputCount,
		Inputs:      txInputs,
		OutputCount: outputCount,
		Outputs:     txOutputs,
		Locktime:    blkTransactions[len(blkTransactions)-4:],
	}

	return txData, nil
}

/*
parseTransactionBlockSize function is used in parseBlockTransactions function to parse the tx count in the transaction block.

Leading byte determines how many more bytes to read to figure out tx count; leading byte is defined as blkTranSize[0].

- if leading byte = 'fc' then parse next byte

- if leading byte = 'fd' then parse next two bytes

- if leading byte = 'fe' then parse next four bytes

- if leading byte = 'ff' then parse next eight bytes

for more info read https://learnmeabitcoin.com/technical/general/compact-size/
*/
func ParseTransactionBlockSize(blkTranSize []byte) (int64, int, error) {
	leadingByte := strings.ToUpper(fmt.Sprintf("%x", blkTranSize[0]))
	// fmt.Printf("leading byte: %v\n", leadingByte)
	if leadingByte <= "FC" && leadingByte >= "0" {
		txCount, err := strconv.ParseInt(leadingByte, 16, 16)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return -1, -1, errors.New(errMsg)
		}
		return txCount, 1, nil
	} else if leadingByte == "FD" {
		number := ByteSwapStr(fmt.Sprintf("%x", blkTranSize[1:3]))
		txCount, err := strconv.ParseInt(number, 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return int64(-1), -1, errors.New(errMsg)
		}
		return txCount, 2, nil
	} else if leadingByte == "FE" {
		number := ByteSwapStr(fmt.Sprintf("%x", blkTranSize[1:5]))
		txCount, err := strconv.ParseInt(number, 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return int64(-1), -1, errors.New(errMsg)
		}
		return txCount, 4, nil
	} else if leadingByte == "FF" {
		number := ByteSwapStr(fmt.Sprintf("%x", blkTranSize[1:9]))
		txCount, err := strconv.ParseInt(number, 16, 64)
		if err != nil {
			errMsg := fmt.Sprintf("can not swap TxCount bytes in parseTransactionBlockSize() function.\nerror: %v\n", err)
			return int64(-1), -1, errors.New(errMsg)
		}
		return txCount, 8, nil
	} else {
		errMsg := fmt.Sprintln("did not expect to return outside of if statement in parseTransactionBlockSize() function. leading byte does not match anything in table from https://learnmeabitcoin.com/technical/general/compact-size/#structure")
		return int64(-1), -1, errors.New(errMsg)
	}
}

// Output a single blocks details to the terminal.
// Used in ParseBlocks function.
func printBlock(block BlockData, tmplFile string) {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		log.Fatalf("can not find go lang template file\nerror: %v\n", err)
	}
	err = tmpl.Execute(os.Stdout, block)
	if err != nil {
		log.Fatalf("can not execute go lang template\nerror: %v\n", err)
	}
}
