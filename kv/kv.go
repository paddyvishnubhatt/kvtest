package kv

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
)

var once sync.Once
var myraft *raft.Raft

const KV_STORE_NAME = "k-v-store"
const LOG_FILE = "logs.dat"
const SFILE = "stable.dat"
const SEPARATOR = "_"

type KV struct {
	Done chan bool
	mtx  sync.RWMutex
	kv   map[string]string
}

func (kv *KV) InitRaft(serverAddr string) {
	fmt.Println("In KV.InitRaft " + serverAddr)
	ctx := context.Background()
	err := kv.SetupRaft(ctx, "node1", serverAddr)
	if err != nil {
		fmt.Printf("Error creating Raft %v", err)
	}

	fmt.Println("Done KV.InitRaft")
}

func (kv *KV) Apply(l *raft.Log) interface{} {
	fmt.Println("In KV.Apply " + string(l.Data))
	kv.mtx.Lock()
	defer kv.mtx.Unlock()
	w := string(l.Data)
	s := strings.Split(w, SEPARATOR)
	kv.kv[s[0]] = s[1]
	return nil
}

func (kv *KV) Restore(r io.ReadCloser) error {
	fmt.Println("In KV.Restore TBD")

	return nil
}

func (kv *KV) Snapshot() (raft.FSMSnapshot, error) {
	fmt.Println("In KV.Snapshot TBD")

	return nil, nil
}

func (kv *KV) Get(key string) string {
	fmt.Println("KV.Getting " + key)
	kv.mtx.RLock()
	defer kv.mtx.RUnlock()
	return kv.kv[key]
}

func (kv *KV) Put(key string, val string) error {
	fmt.Println("KV.Putting " + key + " " + val)
	f := myraft.Apply([]byte(key+SEPARATOR+val), time.Second)
	if err := f.Error(); err != nil {
		return rafterrors.MarkRetriable(err)
	}
	return nil
}

func (kv *KV) SetupRaft(ctx context.Context, myID, myAddress string) error {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)

	ldb, err := boltdb.NewBoltStore(LOG_FILE)
	if err != nil {
		return fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, "logs.dat", err)
	}

	sdb, err := boltdb.NewBoltStore(SFILE)
	if err != nil {
		return fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, "stable.dat", err)
	}

	fss, err := raft.NewFileSnapshotStore(".", 3, os.Stderr)
	if err != nil {
		return fmt.Errorf(`raft.NewFileSnapshotStore: %v`, err)
	}

	tm := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithInsecure()})

	r, err := raft.NewRaft(c, kv, ldb, sdb, fss, tm.Transport())
	if err != nil {
		return fmt.Errorf("raft.NewRaft: %v", err)
	}

	once.Do(func() {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(myID),
					Address:  raft.ServerAddress(myAddress),
				},
			},
		}
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
		kv.kv = make(map[string]string)
	})

	myraft = r

	return nil
}
