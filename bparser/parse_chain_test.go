package bparser_test

import (
	"fmt"
	"testing"

	"github.com/davidhintelmann/blockchain/bparser"
)

const (
	genesisBlock   = "f9beb4d91d0100000100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
	blk00000Height = 0
)

var (
	geneisBlockDec = []byte{249, 190, 180, 217, 29, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 59, 163, 237, 253, 122, 123, 18, 178, 122, 199, 44, 62, 103, 118, 143, 97, 127, 200, 27, 195, 136, 138, 81, 50, 58, 159, 184, 170, 75, 30, 94, 74, 41, 171, 95, 73, 255, 255, 0, 29, 29, 172, 43, 124, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 77, 4, 255, 255, 0, 29, 1, 4, 69, 84, 104, 101, 32, 84, 105, 109, 101, 115, 32, 48, 51, 47, 74, 97, 110, 47, 50, 48, 48, 57, 32, 67, 104, 97, 110, 99, 101, 108, 108, 111, 114, 32, 111, 110, 32, 98, 114, 105, 110, 107, 32, 111, 102, 32, 115, 101, 99, 111, 110, 100, 32, 98, 97, 105, 108, 111, 117, 116, 32, 102, 111, 114, 32, 98, 97, 110, 107, 115, 255, 255, 255, 255, 1, 0, 242, 5, 42, 1, 0, 0, 0, 67, 65, 4, 103, 138, 253, 176, 254, 85, 72, 39, 25, 103, 241, 166, 113, 48, 183, 16, 92, 214, 168, 40, 224, 57, 9, 166, 121, 98, 224, 234, 31, 97, 222, 182, 73, 246, 188, 63, 76, 239, 56, 196, 243, 85, 4, 229, 30, 193, 18, 222, 92, 56, 77, 247, 186, 11, 141, 87, 138, 76, 112, 43, 107, 241, 29, 95, 172, 0, 0, 0, 0, 249, 190}
)

func TestParseGenesis(t *testing.T) {
	// bparser.ParseBlocks(geneisBlockDec, 0, blk00000Height, []byte{0})
	magicNumber := bparser.ParseMagicNumber(geneisBlockDec)
	mn := fmt.Sprintf("%X", magicNumber)
	if mn != "F9BEB4D9" {
		t.Errorf("Expected magic number to equal %s, but got %s\n", "F9BEB4D9", mn)
	}

	blockSize, err := bparser.ParseBlockSize(geneisBlockDec)
	if err != nil {
		t.Errorf("could not parse int, error: %v\n", err)
	}

	if blockSize != 285 {
		t.Errorf("Expected block size to equal %s, but got %d\n", "285", blockSize)
	}
}

func ExampleByteSwap() {
	str := bparser.ByteSwap("6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000")
	fmt.Println(str)
	// Output:
	// 000000000019D6689C085AE165831E934FF763AE46A2A6C172B3F1B60A8CE26F
}
