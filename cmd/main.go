package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/davidhintelmann/blockchain/bparser"
)

const (
	gensisBlock     = "6FE28C0AB6F1B372C1A6A246AE63F74F931E8365E15A089C68D6190000000000"
	gensisBlockSwap = "000000000019D6689C085AE165831E934FF763AE46A2A6C172B3F1B60A8CE26F"
	blocksFilePath  = "C:\\Users\\david\\OneDrive\\Documents\\code\\python\\Blockchain\\Bitcoin\\data\\bitcoin_data\\"
	blk00000Height  = 119_965
)

func main() {
	fmt.Printf("Gensis Block Hash: %s\nGensis Block Swap: %s\n", gensisBlock, gensisBlockSwap)

	// h := sha256.New()
	// h.Write(genesisByte)
	// bs := h.Sum(nil)
	// fmt.Printf("Gensis Blk SHA256: %x\n", bs)

	fmt.Println(bparser.ByteSwap(gensisBlock))
	fmt.Println(bparser.ByteSwapArray(gensisBlock))

	fmt.Println(filepath.Dir(blocksFilePath))

	globStart := time.Now()
	matches, err := bparser.GlobDat(blocksFilePath)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	fmt.Printf("duration of glob: %v\n", time.Since(globStart))

	fmt.Println(len(matches))
	fmt.Println(matches[0])

	fileStart := time.Now()
	file, err := os.Open(matches[0])
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer file.Close()

	readAll, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	fmt.Printf("duration of file read: %v\n", time.Since(fileStart))

	fmt.Println(readAll[:295])
	blkFileLen := len(readAll)
	p := message.NewPrinter(language.English)
	p.Printf("block file length: %d\n", blkFileLen)
	fmt.Println()
	fmt.Println(blk00000Height)
	bparser.ParseBlocks(readAll, 0, blk00000Height, []byte{0})
}
