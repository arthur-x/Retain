package src

import (
	"io"
	"io/fs"
	"log"
	"os"
)

func upload(file fs.DirEntry, client RPCClient, blockStoreAddr string) (newHashes []string) {
	var succ bool
	f, err := os.Open(ConcatPath(client.BaseDir, file.Name()))
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	buf := make([]byte, client.BlockSize)
	for {
		nbytes, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Panicln(err)
			}
			break
		}
		newBlock := Block{BlockData: buf[:nbytes], BlockSize: int32(nbytes)}
		newHashes = append(newHashes, GetBlockHashString(newBlock.BlockData))
		if err = client.PutBlock(&newBlock, blockStoreAddr, &succ); err != nil {
			log.Panicln(err)
		}
	}
	return newHashes
}

func download(remoteMetaData *FileMetaData, client RPCClient, blockStoreAddr string) {
	remoteHashes := remoteMetaData.BlockHashList
	if len(remoteHashes) == 1 && remoteHashes[0] == "0" {
		if err := os.RemoveAll(ConcatPath(client.BaseDir, remoteMetaData.Filename)); err != nil {
			log.Panicln(err)
		}
		return
	}
	f, err := os.Create(ConcatPath(client.BaseDir, remoteMetaData.Filename))
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	for _, hash := range remoteHashes {
		remoteBlock := Block{}
		if err := client.GetBlock(hash, blockStoreAddr, &remoteBlock); err != nil {
			log.Panicln(err)
		}
		if _, err := f.Write(remoteBlock.BlockData); err != nil {
			log.Panicln(err)
		}
	}
}

func ClientSync(client RPCClient) {
	var blockStoreAddr string
	var Version int32
	if err := client.GetBlockStoreAddr(&blockStoreAddr); err != nil {
		log.Panicln(err)
	}
	localFileMetaMap, err := LoadMetaFromMetaFile(client.BaseDir)
	if err != nil {
		log.Panicln(err)
	}
	files, err := os.ReadDir(client.BaseDir)
	if err != nil {
		log.Panicln(err)
	}
	localFiles := make(map[string]int)
	for _, file := range files {
		if file.Name() == DEFAULT_META_FILENAME {
			continue
		}
		localFiles[file.Name()] = 1
		changed := false
		localMetaData, prs := localFileMetaMap[file.Name()]
		newHashes := upload(file, client, blockStoreAddr)
		if !prs || len(newHashes) != len(localMetaData.BlockHashList) {
			changed = true
		} else {
			for i, hash := range newHashes {
				if localMetaData.BlockHashList[i] != hash {
					changed = true
					break
				}
			}
		}
		if !changed {
			continue
		}
		newMetaData := FileMetaData{Filename: file.Name(), BlockHashList: newHashes}
		if !prs {
			newMetaData.Version = 1
		} else {
			newMetaData.Version = localMetaData.Version + 1
		}
		if err = client.UpdateFile(&newMetaData, &Version); err != nil {
			log.Panicln(err)
		}
		if Version != -1 {
			localFileMetaMap[file.Name()] = &newMetaData
		}
	}
	for filename, localMetaData := range localFileMetaMap {
		if _, prs := localFiles[filename]; !prs && (len(localMetaData.BlockHashList) != 1 || localMetaData.BlockHashList[0] != "0") {
			newMetaData := FileMetaData{Filename: filename, BlockHashList: []string{"0"}}
			newMetaData.Version = localMetaData.Version + 1
			if err = client.UpdateFile(&newMetaData, &Version); err != nil {
				log.Panicln(err)
			}
			if Version != -1 {
				localFileMetaMap[filename] = &newMetaData
			}
		}
	}
	remoteFileMetaMap := make(map[string]*FileMetaData)
	if err := client.GetFileInfoMap(&remoteFileMetaMap); err != nil {
		log.Panicln(err)
	}
	for filename, remoteMetaData := range remoteFileMetaMap {
		if localMetaData, prs := localFileMetaMap[filename]; !prs || localMetaData.Version < remoteMetaData.Version {
			localFileMetaMap[filename] = remoteMetaData
			download(remoteMetaData, client, blockStoreAddr)
		}
	}
	WriteMetaFile(localFileMetaMap, client.BaseDir)
}
