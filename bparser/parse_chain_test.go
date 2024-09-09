package bparser_test

import (
	"fmt"
	"testing"

	"github.com/davidhintelmann/blockchain/bparser"
)

// const (
// 	genesisBlock   = "f9beb4d91d0100000100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
// 	blk00000Height = 0
// )

var (
	geneisBlockDec = []byte{249, 190, 180, 217, 29, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 59, 163, 237, 253, 122, 123, 18, 178, 122, 199, 44, 62, 103, 118, 143, 97, 127, 200, 27, 195, 136, 138, 81, 50, 58, 159, 184, 170, 75, 30, 94, 74, 41, 171, 95, 73, 255, 255, 0, 29, 29, 172, 43, 124, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 77, 4, 255, 255, 0, 29, 1, 4, 69, 84, 104, 101, 32, 84, 105, 109, 101, 115, 32, 48, 51, 47, 74, 97, 110, 47, 50, 48, 48, 57, 32, 67, 104, 97, 110, 99, 101, 108, 108, 111, 114, 32, 111, 110, 32, 98, 114, 105, 110, 107, 32, 111, 102, 32, 115, 101, 99, 111, 110, 100, 32, 98, 97, 105, 108, 111, 117, 116, 32, 102, 111, 114, 32, 98, 97, 110, 107, 115, 255, 255, 255, 255, 1, 0, 242, 5, 42, 1, 0, 0, 0, 67, 65, 4, 103, 138, 253, 176, 254, 85, 72, 39, 25, 103, 241, 166, 113, 48, 183, 16, 92, 214, 168, 40, 224, 57, 9, 166, 121, 98, 224, 234, 31, 97, 222, 182, 73, 246, 188, 63, 76, 239, 56, 196, 243, 85, 4, 229, 30, 193, 18, 222, 92, 56, 77, 247, 186, 11, 141, 87, 138, 76, 112, 43, 107, 241, 29, 95, 172, 0, 0, 0, 0}
)

func TestParseGenesis(t *testing.T) {
	// bparser.ParseBlocks(geneisBlockDec, 0, blk00000Height, []byte{0})
	magicNumber := bparser.ParseMagicNumber(geneisBlockDec)
	mn := fmt.Sprintf("%X", magicNumber)
	if mn != "F9BEB4D9" {
		t.Errorf("Expected magic number to equal %s, but got %s\n", "F9BEB4D9", mn)
	}

	// test ParseBlockSizeFunc function
	blockSize, err := bparser.ParseBlockSizeFunc(geneisBlockDec)
	if err != nil {
		t.Errorf("could not parse int, error: %v\n", err)
	}

	if blockSize != 285 {
		t.Errorf("Expected block size to equal %s, but got %d\n", "285", blockSize)
	}

	// test ParseBlockSize function
	blockSize, err = bparser.ParseBlockSize(geneisBlockDec[4:8])
	if err != nil {
		t.Errorf("could not parse int, error: %v\n", err)
	}

	if blockSize != 285 {
		t.Errorf("Expected block size to equal %s, but got %d\n", "285", blockSize)
	}

	// test ParseBlockRaw function
	_, err = bparser.ParseBlockRaw(geneisBlockDec)
	if err != nil {
		t.Errorf("could not parse int, error: %v\n", err)
	}

	// fmt.Printf("Magic Number: %v\nBlock Size: %v\n", blockRaw.MagicNumber, blockRaw.Size)
	// fmt.Printf("Version: %v\nPrev Block: %v\nMerkle Root: %v\nTimestamp: %v\nBits: %v\nNonce: %v\n", blockRaw.BlockHeader.Version, blockRaw.BlockHeader.PrevBlock, blockRaw.BlockHeader.MerkleRoot, blockRaw.BlockHeader.Timestamp, blockRaw.BlockHeader.Bits, blockRaw.BlockHeader.Nonce)

	// test ParseBlockStr function
	_, err = bparser.ParseBlockStr(geneisBlockDec)
	if err != nil {
		t.Errorf("could not parse int, error: %v\n", err)
	}
	// fmt.Printf("Magic Number: %v\nBlock Size: %v\n", block.MagicNumber, block.Size)
	// fmt.Printf("Version: %v\nPrev Block: %v\nMerkle Root: %v\nTimestamp: %v\nBits: %v\nNonce: %v\n", block.BlockHeader.Version, block.BlockHeader.PrevBlock, block.BlockHeader.MerkleRoot, block.BlockHeader.Timestamp, block.BlockHeader.Bits, block.BlockHeader.Nonce)
}

func TestParseBlockSize(t *testing.T) {
	// test ParseBlockSizeRaw
	_, err := bparser.ParseBlockSizeRaw(geneisBlockDec)
	if err != nil {
		t.Errorf("test ParseBlockSizeRaw could not parse block to find the size of the next block, error: %v\n", err)
	}

	_, err = bparser.ParseBlockSizeRaw(geneisBlockDec[:2])
	if err == nil {
		t.Errorf("test ParseBlockSizeRaw expected to fail since it was given a slice which is too small, error: %v\n", err)
	}

	// test ParseBlockSize
	blockSize, err := bparser.ParseBlockSizeFunc(geneisBlockDec)
	if err != nil {
		t.Errorf("test ParseBlockSize could not parse block to find the size of the next block, error: %v\n", err)
	} else if blockSize != 285 {
		t.Errorf("test ParseBlockSize expected block size of 285 but got %d, error: %v\n", blockSize, err)
	}

	_, err = bparser.ParseBlockSize(geneisBlockDec[:2])
	if err == nil {
		t.Errorf("test ParseBlockSize expected to fail since it was given a slice which is too small, error: %v\n", err)
	}
}

func TestParseTransactionBlockSize(t *testing.T) {
	tests := []struct {
		name            string
		transationBytes []byte
		want            int64
	}{
		{
			name:            "genesis block size, leading byte <= FC",
			transationBytes: []byte{01},
			want:            int64(1),
		},
		{
			name:            "single byte - 252 leading byte <= FC",
			transationBytes: []byte{252},
			want:            int64(252),
		},
		{
			name:            "next two bytes - 253, 232, 03 leading byte == FD",
			transationBytes: []byte{253, 232, 03},
			want:            int64(1_000),
		},
		{
			name:            "next four bytes - 254, 160, 134, 01, 00 leading byte == FE",
			transationBytes: []byte{254, 160, 134, 01, 00},
			want:            int64(100_000),
		},
		{
			name:            "next eight bytes - 255, 00, 228, 11, 84, 02, 00, 00, 00 leading byte == FF",
			transationBytes: []byte{255, 00, 228, 11, 84, 02, 00, 00, 00},
			want:            int64(10_000_000_000),
		},
	}

	// test genesis block first
	blockSize, err := bparser.ParseBlockSizeFunc(geneisBlockDec)
	if err != nil {
		t.Errorf("test ParseBlockSize could not parse block to find the size of the next block, error: %v\n", err)
	}

	genesisTransactionBlock := geneisBlockDec[88 : blockSize+8]
	t.Run("parse genesis block from const in test file", func(t *testing.T) {
		got, _, err := bparser.ParseTransactionBlockSize(genesisTransactionBlock)
		if err != nil {
			t.Errorf("test ParseTransactionBlockSize could not parse block to find the size of the next block, error: %v\n", err)
		} else if got != 1 {
			t.Errorf("ParseTransactionBlockSize(blkTranSize:%v) got = %v, want %v", genesisTransactionBlock, got, 1)
		}
	})

	// test leading bytes from the table at https://learnmeabitcoin.com/technical/general/compact-size/#structure
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _ := bparser.ParseTransactionBlockSize(tt.transationBytes)
			if got != tt.want {
				t.Errorf("ParseTransactionBlockSize(blkTranSize:%v) got = %v, want %v", tt.transationBytes, got, tt.want)
			}
		})
	}

}

// func TestParseTransactionBlock(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		input   []byte
// 		padding int
// 		want    bparser.TxData
// 	}{
// 		{
// 			name:    "parse gensis block",
// 			input:   []byte{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 77, 4, 255, 255, 0, 29, 1, 4, 69, 84, 104, 101, 32, 84, 105, 109, 101, 115, 32, 48, 51, 47, 74, 97, 110, 47, 50, 48, 48, 57, 32, 67, 104, 97, 110, 99, 101, 108, 108, 111, 114, 32, 111, 110, 32, 98, 114, 105, 110, 107, 32, 111, 102, 32, 115, 101, 99, 111, 110, 100, 32, 98, 97, 105, 108, 111, 117, 116, 32, 102, 111, 114, 32, 98, 97, 110, 107, 115, 255, 255, 255, 255, 1, 0, 242, 5, 42, 1, 0, 0, 0, 67, 65, 4, 103, 138, 253, 176, 254, 85, 72, 39, 25, 103, 241, 166, 113, 48, 183, 16, 92, 214, 168, 40, 224, 57, 9, 166, 121, 98, 224, 234, 31, 97, 222, 182, 73, 246, 188, 63, 76, 239, 56, 196, 243, 85, 4, 229, 30, 193, 18, 222, 92, 56, 77, 247, 186, 11, 141, 87, 138, 76, 112, 43, 107, 241, 29, 95, 172, 0, 0, 0, 0},
// 			padding: 1,
// 			want: bparser.TxData{
// 				Version:    int64(1),
// 				InputCount: int64(1),
// 				Inputs: []bparser.TxInputs{
// 					{
// 						TxId:          "0000000000000000000000000000000000000000000000000000000000000000",
// 						Vout:          "ffffffff",
// 						ScriptSigSize: int64(77),
// 						ScriptSig:     "04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73",
// 						Sequence:      "ffffffff",
// 					},
// 				},
// 				OutputCount: int64(1),
// 				Outputs: []bparser.TxOutputs{
// 					{
// 						Amount:           []byte{0, 242, 5, 42, 1, 0, 0, 0},
// 						ScriptPubKeySize: int64(43),
// 						ScriptPubKey:     []byte{65, 4, 103, 138, 253, 176, 254, 85, 72, 39, 25, 103, 241, 166, 113, 48, 183, 16, 92, 214, 168, 40, 224, 57, 9, 166, 121, 98, 224, 234, 31, 97, 222, 182, 73, 246, 188, 63, 76, 239, 56, 196, 243, 85, 4, 229, 30, 193, 18, 222, 92, 56, 77, 247, 186, 11, 141, 87, 138, 76, 112, 43, 107, 241, 29, 95, 172},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:    "parse transaction block for block height 672,119",
// 			input:   []byte{1, 0, 0, 0, 1, 59, 165, 212, 161, 9, 141, 155, 79, 44, 51, 107, 193, 189, 157, 137, 26, 146, 136, 243, 179, 89, 182, 137, 30, 118, 132, 21, 248, 36, 42, 30, 59, 1, 0, 0, 0, 107, 72, 48, 69, 2, 33, 0, 241, 77, 54, 196, 153, 187, 17, 32, 238, 11, 31, 180, 251, 105, 111, 28, 42, 42, 114, 222, 121, 224, 245, 29, 210, 143, 46, 224, 29, 161, 180, 246, 2, 32, 15, 35, 36, 53, 92, 213, 223, 136, 187, 39, 77, 166, 240, 141, 247, 93, 114, 12, 193, 143, 190, 225, 8, 69, 220, 206, 46, 253, 14, 141, 79, 166, 1, 33, 3, 161, 115, 190, 132, 127, 152, 90, 10, 217, 7, 87, 107, 209, 97, 144, 108, 177, 197, 85, 203, 128, 242, 80, 131, 34, 139, 23, 83, 88, 69, 184, 186, 255, 255, 255, 255, 2, 215, 37, 3, 0, 0, 0, 0, 0, 25, 118, 169, 20, 65, 160, 218, 69, 116, 194, 64, 156, 150, 113, 176, 36, 245, 207, 103, 118, 106, 249, 119, 134, 136, 172, 14, 73, 9, 0, 0, 0, 0, 0, 23, 169, 20, 203, 205, 60, 129, 136, 102, 212, 187, 36, 245, 189, 212, 99, 222, 23, 150, 52, 158, 121, 40, 135, 0, 0, 0, 0},
// 			padding: 1,
// 			want: bparser.TxData{
// 				Version:    int64(1),
// 				InputCount: int64(1),
// 				Inputs: []bparser.TxInputs{
// 					{
// 						TxId:          "3ba5d4a1098d9b4f2c336bc1bd9d891a9288f3b359b6891e768415f8242a1e3b",
// 						Vout:          "01000000",
// 						ScriptSigSize: int64(107),
// 						ScriptSig:     "483045022100f14d36c499bb1120ee0b1fb4fb696f1c2a2a72de79e0f51dd28f2ee01da1b4f602200f2324355cd5df88bb274da6f08df75d720cc18fbee10845dcce2efd0e8d4fa6012103a173be847f985a0ad907576bd161906cb1c555cb80f25083228b17535845b8ba",
// 						Sequence:      "ffffffff",
// 					},
// 				},
// 				OutputCount: int64(2),
// 				Outputs: []bparser.TxOutputs{
// 					{
// 						Amount:           []byte{215, 37, 3, 0, 0, 0, 0, 0},
// 						ScriptPubKeySize: int64(17),
// 						ScriptPubKey:     []byte{169, 20, 203, 205, 60, 129, 136, 102, 212, 187, 36, 245, 189, 212, 99, 222, 23, 150, 52, 158, 121, 40, 135},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, _ := bparser.ParseBlockTx(tt.input, 1)
// 			gotVersion, gotInputCount, gotOutputCount := got.Version, got.InputCount, got.OutputCount
// 			// gotInputs, gotOutputs := got.Inputs, got.Outputs
// 			if gotVersion != tt.want.Version || gotInputCount != tt.want.InputCount || gotOutputCount != tt.want.OutputCount {
// 				t.Errorf("ParseTransactionBlockSize(blkTranSize:%v, %v) got = %v, want %v", tt.input, tt.padding, got, tt.want)
// 			}
// 		})
// 	}
// }

func ExampleByteSwap() {
	str := bparser.ByteSwap("6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000")
	fmt.Println(str)
	// Output:
	// 000000000019D6689C085AE165831E934FF763AE46A2A6C172B3F1B60A8CE26F
}
