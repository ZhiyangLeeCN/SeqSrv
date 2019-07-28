package gen

import (
	"sync"
)

type SeqNode struct {
	id       uint64
	seqFile  *SeqFile
	syncLock *sync.Mutex
}

func NewSeqNode(seqFile *SeqFile, nodeId uint64) (*SeqNode, error) {
	sn := SeqNode{}
	sn.id = nodeId
	sn.seqFile = seqFile
	sn.syncLock = &sync.Mutex{}
	return &sn, nil
}

func (sn *SeqNode) NextID() (uint64, error) {
	defer sn.syncLock.Unlock()

	sn.syncLock.Lock()
	newSeq, err := sn.seqFile.NextSeq(1)
	if err != nil {
		return 0, err
	}

	r := newSeq<<10 | sn.id
	return r, nil
}
