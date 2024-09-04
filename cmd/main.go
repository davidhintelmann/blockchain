package main

import (
	"fmt"

	"github.com/davidhintelmann/blockchain/bparser"
)

const (
	gensisBlock     = "6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000"
	gensisBlockSwap = "000000000019D6689C085AE165831E934FF763AE46A2A6C172B3F1B60A8CE26F"
	blocksFilePath  = "C:\\Users\\david\\OneDrive\\Documents\\code\\python\\Blockchain\\Bitcoin\\data\\bitcoin_data"
)

func main() {
	fmt.Printf("Gensis Block Hash: %s\nGensis Block Swap: %s\n", gensisBlock, gensisBlockSwap)

	// h := sha256.New()
	// h.Write(genesisByte)
	// bs := h.Sum(nil)
	// fmt.Printf("Gensis Blk SHA256: %x\n", bs)

	fmt.Println(bparser.ByteSwap(gensisBlock))
	fmt.Println(bparser.ByteSwapArray(gensisBlock))
}
