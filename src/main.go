package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
)

////////////////////////////////////////////////////////////////

type Block struct {
	Stage         uint64
	PrevBlockHash []byte
	Data          []byte
	Nonce         []byte
	Hash          []byte
}

func NewBlock(stage uint64, prevBlockHash []byte, data string) *Block {
	block := &Block{
		Stage:         stage,
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		Nonce:         []byte{},
		Hash:          []byte{},
	}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	stage := make([]byte, 8)
	binary.LittleEndian.PutUint64(stage, b.Stage)
	headers := bytes.Join([][]byte{
		stage,
		b.PrevBlockHash,
		b.Data,
		b.Nonce,
	}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func (b *Block) Prove(bc *Blockchain) {
	b.Nonce = make([]byte, 8)
	hashNum := big.NewInt(0)
	for {
		rand.Read(b.Nonce)
		b.SetHash()
		hashNum.SetBytes(b.Hash)
		if bc.Target.Cmp(hashNum) > 0 {
			return
		}
	}
}

////////////////////////////////////////////////////////////////

type Stage struct {
	Blocks []*Block
}

func (s *Stage) containsHash(hash []byte) bool {
	var same bool
	for _, element := range s.Blocks {
		same = true
		if !bytes.Equal(hash, element.Hash) {
			same = false
		}
		if same {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////

type Blockchain struct {
	Stages []*Stage
	Target *big.Int
}

func NewBlockchain(targetBits int) *Blockchain {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &Blockchain{[]*Stage{&Stage{[]*Block{NewBlock(0, []byte{}, "GENBLOCK")}}}, target}
}

func (bc *Blockchain) PrintBlockchain() {
	for index, element1 := range bc.Stages {
		fmt.Printf("Stage %d\n", index)
		for _, element2 := range element1.Blocks {
			fmt.Printf("\tPrevBlockHash: %s, Data: %s, Hash: %s\n", string(element2.PrevBlockHash), element2.Data, string(element2.Hash))
		}
	}
}

func (bc *Blockchain) AddBlock(b *Block) bool {
	b.SetHash()
	hashNum := big.NewInt(0)
	hashNum.SetBytes(b.Hash)
	if bc.Target.Cmp(hashNum) > 0 && b.Stage > 0 {
		if b.Stage == uint64(len(bc.Stages)) {
			bc.Stages = append(bc.Stages, &Stage{[]*Block{}})
		}
		if bc.Stages[b.Stage-1].containsHash(b.PrevBlockHash) {
			bc.Stages[b.Stage].Blocks = append(bc.Stages[b.Stage].Blocks, b)
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////

func main() {
	x := NewBlockchain(26)
	b := NewBlock(1, x.Stages[0].Blocks[0].Hash, "Here is some data!")
	b.Prove(x)
	wasAdded := x.AddBlock(b)
	fmt.Println(wasAdded)
	x.PrintBlockchain()
}
