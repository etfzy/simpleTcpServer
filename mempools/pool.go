package mempools

import (
	mem "github.com/etfzy/mempool/base"
)

type MemPools interface {
	GetProtoMems(num uint64) *mem.Buffer[byte]
	GetContentMems(num uint64) *mem.Buffer[byte]
	PutProtoMems(m *mem.Buffer[byte])
	PutContentMems(m *mem.Buffer[byte])
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

func (mp *memPools) GetProtoMems(num uint64) *mem.Buffer[byte] {
	if mp.protoMem == nil {
		return mem.NewBuffer[byte](int(num))
	}

	return mp.protoMem.Get(num)
}

func (mp *memPools) GetContentMems(num uint64) *mem.Buffer[byte] {
	if mp.contentMem == nil {
		return mem.NewBuffer[byte](int(num))
	}

	return mp.contentMem.Get(num)
}

func (mp *memPools) PutProtoMems(m *mem.Buffer[byte]) {
	if mp.protoMem == nil || m == nil {
		return
	}
	mp.protoMem.PutBack(m)
	return
}

func (mp *memPools) PutContentMems(m *mem.Buffer[byte]) {
	if mp.contentMem == nil || m == nil {
		return
	}
	mp.contentMem.PutBack(m)
	return
}
