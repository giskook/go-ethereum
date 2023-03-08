package snapshot

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
	"math/big"
	"time"
)

type Retriever interface {
	DecodeAccount(bz []byte) (nonce uint64, balance *big.Int, root common.Hash, codeHash []byte)
}

func NewCustom(diskdb ethdb.KeyValueStore, triedb *trie.Database, cache int, root common.Hash, async bool, rebuild bool, recovery bool, retriever Retriever) (*Tree, error) {
	// Create a new, empty snapshot tree
	snap := &Tree{
		diskdb:    diskdb,
		triedb:    triedb,
		cache:     cache,
		layers:    make(map[common.Hash]snapshot),
		Retriever: retriever,
	}
	if !async {
		defer snap.waitBuild()
	}
	// Attempt to load a previously persisted snapshot and rebuild one if failed
	head, disabled, err := loadSnapshot(diskdb, triedb, cache, root, recovery)
	if disabled {
		log.Warn("Snapshot maintenance disabled (syncing)")
		return snap, nil
	}
	if err != nil {
		if rebuild {
			log.Warn("Failed to load snapshot, regenerating", "err", err)
			snap.Rebuild(root)
			return snap, nil
		}
		return nil, err // Bail out the error, don't rebuild automatically.
	}
	// Existing snapshot loaded, seed all the layers
	for head != nil {
		snap.layers[head.Root()] = head
		head = head.Parent()
	}
	return snap, nil
}

func generateSnapshotCustom(diskdb ethdb.KeyValueStore, triedb *trie.Database, cache int, root common.Hash, retriever Retriever) *diskLayer {
	// Create a new disk layer with an initialized state marker at zero
	var (
		stats     = &generatorStats{start: time.Now()}
		batch     = diskdb.NewBatch()
		genMarker = []byte{} // Initialized but empty!
	)
	rawdb.WriteSnapshotRoot(batch, root)
	journalProgress(batch, genMarker, stats)
	if err := batch.Write(); err != nil {
		log.Crit("Failed to write initialized state marker", "err", err)
	}
	base := &diskLayer{
		diskdb:     diskdb,
		triedb:     triedb,
		root:       root,
		cache:      fastcache.New(cache * 1024 * 1024),
		genMarker:  genMarker,
		genPending: make(chan struct{}),
		genAbort:   make(chan chan *generatorStats),
		Retriever:  retriever,
	}
	go base.generate(stats)
	log.Debug("Start snapshot generation", "root", root)
	return base
}
