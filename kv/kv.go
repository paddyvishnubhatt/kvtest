package kv

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var myraft *raft.Raft

const KV_STORE_NAME = "k-v-store"
const DATA_DIR = "/tmp/raft_dir"
const LOG_FILE = "logs.dat"
const SFILE = "stable.dat"
const SEPARATOR = "_"

type KV struct {
	Done chan bool
	mtx  sync.RWMutex
	kv   map[string]string
}

func (kv *KV) InitRaft(serverAddr string, id string, bs bool) {
	fmt.Printf("In KV.InitRaft %v %v %v\n", serverAddr, id, bs)
	ctx := context.Background()
	r, tm, err := kv.SetupRaft(ctx, id, serverAddr, bs)
	myraft = r
	if err != nil {
		fmt.Printf("Error creating Raft %v", err)
	}
	s := grpc.NewServer()
	tm.Register(s)
	leaderhealth.Setup(r, s, []string{"Example"})
	raftadmin.Register(s, r)
	reflection.Register(s)
	_, port, err := net.SplitHostPort(serverAddr)
	if err != nil {
		log.Fatalf("failed to parse local address (%q): %v", serverAddr, err)
	}
	sock, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(sock); err != nil {
		log.Fatalf("failed to serve: %v", err)
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

func (kv *KV) AddVoter(voter string, id string) error {
	fmt.Println("KV.AddVoter " + voter + " " + id)
	f := myraft.AddVoter(raft.ServerID(id), raft.ServerAddress(voter), 0, time.Second)
	if err := f.Error(); err != nil {
		return rafterrors.MarkRetriable(err)
	}
	return nil
}

func (kv *KV) SetupRaft(ctx context.Context, myID, myAddress string, bs bool) (*raft.Raft, *transport.Manager, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)
	baseDir := filepath.Join(DATA_DIR, myID)

	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, LOG_FILE))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, "logs.dat", err)
	}

	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, SFILE))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, "stable.dat", err)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf(`raft.NewFileSnapshotStore: %v`, err)
	}

	tm := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithInsecure()})

	r, err := raft.NewRaft(c, kv, ldb, sdb, fss, tm.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if bs {
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
	}

	kv.kv = make(map[string]string)

	return r, tm, nil
}
