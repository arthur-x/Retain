package src

import (
	context "context"
	"sync"
)

type BlockStore struct {
	mu       sync.Mutex
	BlockMap map[string]*Block
	UnimplementedBlockStoreServer
}

func (bs *BlockStore) GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if b, prs := bs.BlockMap[blockHash.Hash]; prs {
		return &Block{BlockData: b.BlockData, BlockSize: b.BlockSize}, nil
	}
	return &Block{}, nil
}

func (bs *BlockStore) PutBlock(ctx context.Context, block *Block) (*Success, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	hashString := GetBlockHashString(block.BlockData)
	if _, prs := bs.BlockMap[hashString]; prs {
		return &Success{Flag: false}, nil
	}
	bs.BlockMap[hashString] = block
	return &Success{Flag: true}, nil
}

// Given a list of hashes “in”, returns a list containing the
// subset of in that are stored in the key-value store
func (bs *BlockStore) HasBlocks(ctx context.Context, blockHashesIn *BlockHashes) (*BlockHashes, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	blockHashesOut := &BlockHashes{}
	for _, hashString := range blockHashesIn.Hashes {
		if _, prs := bs.BlockMap[hashString]; prs {
			blockHashesOut.Hashes = append(blockHashesOut.Hashes, hashString)
		}
	}
	return blockHashesOut, nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)

func NewBlockStore() *BlockStore {
	return &BlockStore{
		BlockMap: map[string]*Block{},
	}
}
