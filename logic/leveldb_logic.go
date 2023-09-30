package logic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
)

func SaveBlockToDB(block Block, db *leveldb.DB) error {
	blockJson, err := json.Marshal(block)
	if err != nil {
		return err
	}
	key := strconv.Itoa(block.Index)
	err = db.Put([]byte(key), blockJson, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetBlockFromDB(index int, db *leveldb.DB) (Block, error) {
	var block Block
	key := strconv.Itoa(index)
	blockJson, err := db.Get([]byte(key), nil)
	if err != nil {
		return block, err
	}
	err = json.Unmarshal(blockJson, &block)
	if err != nil {
		return block, err
	}
	return block, nil
}

func GetLastBlock(db *leveldb.DB) (Block, error) {
	var lastBlock Block
	iter := db.NewIterator(nil, nil)
	var maxIndex int
	for iter.Next() {
		var currentBlock Block
		err := json.Unmarshal(iter.Value(), &currentBlock)
		if err != nil {
			return Block{}, err
		}
		if currentBlock.Index > maxIndex {
			maxIndex = currentBlock.Index
			lastBlock = currentBlock
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return Block{}, err
	}
	return lastBlock, nil
}

func CountTotalBlocks(db *leveldb.DB) (int, error) {
	iter := db.NewIterator(nil, nil)
	count := 0
	for iter.Next() {
		count++
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetTotalTransactions(db *leveldb.DB, blockIndex int) (int, error) {
	key := []byte(fmt.Sprintf("%d", blockIndex))
	data, err := db.Get(key, nil)
	if err != nil {
		return 0, err
	}
	var block Block
	err = json.Unmarshal(data, &block)
	if err != nil {
		return 0, err
	}
	return len(block.Transactions), nil
}

func CreateGenesis(db *leveldb.DB) {
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	if !iter.First() {
		fmt.Print("\nCreando bloque Genesis.\n\n")
		genesis := Genesis()
		err := SaveBlockToDB(genesis, db)
		if err != nil {
			panic(err)
		}
		Pretty(genesis)
		fmt.Print("\n\nEnter... ")
		var wait int
		fmt.Scanln(&wait)
	}

}


