package src

import (
	context "context"
	"sync"

	"google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MetaStore struct {
	mu             sync.Mutex
	FileMetaMap    map[string]*FileMetaData
	BlockStoreAddr string
	UnimplementedMetaStoreServer
}

func (m *MetaStore) GetFileInfoMap(ctx context.Context, _ *emptypb.Empty) (*FileInfoMap, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return proto.Clone(&FileInfoMap{FileInfoMap: m.FileMetaMap}).(*FileInfoMap), nil
}

func (m *MetaStore) UpdateFile(ctx context.Context, fileMetaData *FileMetaData) (*Version, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := fileMetaData.Version
	if info, prs := m.FileMetaMap[fileMetaData.Filename]; prs && (v != info.Version+1) {
		return &Version{Version: -1}, nil
	}
	m.FileMetaMap[fileMetaData.Filename] = fileMetaData
	return &Version{Version: v}, nil
}

func (m *MetaStore) GetBlockStoreAddr(ctx context.Context, _ *emptypb.Empty) (*BlockStoreAddr, error) {
	return &BlockStoreAddr{Addr: m.BlockStoreAddr}, nil
}

// This line guarantees all method for MetaStore are implemented
var _ MetaStoreInterface = new(MetaStore)

func NewMetaStore(blockStoreAddr string) *MetaStore {
	return &MetaStore{
		FileMetaMap:    map[string]*FileMetaData{},
		BlockStoreAddr: blockStoreAddr,
	}
}
