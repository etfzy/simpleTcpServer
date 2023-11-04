package mempools

import mem "github.com/etfzy/mempool/base"

type MemPools interface {
	GetProtoMems(num uint64) *[]byte
	GetContentMems(num uint64) *[]byte
	PutProtoMems(m *[]byte)
	PutContentMems(m *[]byte)
}

type memPools struct {
	protoMem   mem.LevelsMemPool[byte]
	contentMem mem.LevelsMemPool[byte]
}

func CreateMems(protoLens, contentLens []uint64) MemPools {
	memp := &memPools{}

	memp.protoMem = mem.NewMemPool[byte](protoLens)

	if contentLens != nil {
		memp.contentMem = mem.NewMemPool[byte](contentLens)
	}

	return memp
}

func (mp *memPools) GetProtoMems(num uint64) *[]byte {
	if mp.protoMem == nil {
		temp := make([]byte, num)
		return &temp
	}

	return mp.protoMem.Get(num)
}

func (mp *memPools) GetContentMems(num uint64) *[]byte {
	if mp.contentMem == nil {
		temp := make([]byte, num)
		return &temp
	}

	return mp.contentMem.Get(num)
}

func (mp *memPools) PutProtoMems(m *[]byte) {
	if mp.protoMem == nil {
		return
	}
	mp.protoMem.PutBack(m)
	return
}

func (mp *memPools) PutContentMems(m *[]byte) {
	if mp.contentMem == nil {
		return
	}
	mp.contentMem.PutBack(m)
	return
}
