package gen

import (
	"errors"
	"github.com/ZhiyangLeeCN/SeqSrv/common/allutil"
	"github.com/ZhiyangLeeCN/SeqSrv/common/atomic"
	"strconv"
)

var (
	ReadPersistError = errors.New("read persist file error.")
)

type SeqFile struct {
	fileName string
	curSeq   atomic.Uint64
	maxSeq   uint64
	step     uint64
}

func NewSeqFile(fileName string, initCurSeq uint64, initMaxSeq uint64, step uint64) (*SeqFile, error) {
	sf := SeqFile{}
	sf.fileName = fileName
	sf.curSeq = atomic.MakeUint64(initCurSeq)
	sf.maxSeq = initMaxSeq
	sf.step = step

	err := sf.load()
	if err != nil {
		return nil, err
	}

	return &sf, nil
}

func (sf *SeqFile) load() error {

	readMaxSeq, err := sf.ReadPersistMaxSeq()
	if err != nil {
		_, err := sf.updateAndPersistMaxSeq(sf.maxSeq)
		if err != nil {
			return nil
		}

		return nil
	} else {
		sf.maxSeq = readMaxSeq
		sf.curSeq.Store(readMaxSeq)
		return nil
	}
}

func (sf *SeqFile) updateAndPersistMaxSeq(newMaxSeq uint64) (uint64, error) {
	sf.maxSeq = newMaxSeq
	err := allutil.String2File(sf.GetMaxSeqStr(), sf.fileName)
	if err != nil {
		return 0, err
	}
	return newMaxSeq, nil
}

func (sf *SeqFile) GetMaxSeq() uint64 {
	return sf.maxSeq
}

func (sf *SeqFile) GetMaxSeqStr() string {
	return strconv.FormatUint(sf.maxSeq, 10)
}

func (sf *SeqFile) ReadPersistMaxSeq() (uint64, error) {
	content := allutil.File2String(sf.fileName)
	if content == "" {
		return 0, ReadPersistError
	}

	readMaxSeq, err := strconv.ParseUint(content, 10, 64)
	if err != nil {
		return 0, nil
	}

	return readMaxSeq, nil
}

func (sf *SeqFile) CheckMaxSeqIsNeedPersist() bool {
	curSeq := sf.curSeq.Load()
	return (curSeq + 1) >= sf.maxSeq
}

func (sf *SeqFile) NextSeq(delta uint64) (uint64, error) {
	if sf.CheckMaxSeqIsNeedPersist() {
		newMaxSeq := sf.maxSeq + sf.step
		_, err := sf.updateAndPersistMaxSeq(newMaxSeq)
		if err != nil {
			return 0, err
		}
	}
	return sf.curSeq.Add(delta), nil
}
