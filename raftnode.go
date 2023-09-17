package raftnodes

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shaj13/raft"
	"github.com/shaj13/raft/transport"
	"github.com/shaj13/raft/transport/raftgrpc"
	"google.golang.org/grpc"
)

type RaftNode[DataType any] struct {
	node     *dataNode[DataType]
	raftnode *raft.Node
	running  bool
}

type RaftNodeConfig struct {
	NodeAddr Address
	ApiAddr  Address
	Join     *Address
}

func New[DataType any]() (*RaftNode[DataType], error) {

	return &RaftNode[DataType]{
		node: newDataNode[DataType](),
	}, nil
}

func (rn *RaftNode[DataType]) Set(key string, value *DataType) {
	data, err := json.Marshal(dataItem[DataType]{Key: key, Value: *value})
	if err != nil {
		panic("value can't be marshaled")
	}
	rn.node.Apply(data)
}

func (rn *RaftNode[DataType]) Get(key string) *DataType {
	value, exists := rn.node.data[key]
	if !exists {
		return nil
	}
	return &value
}

func (rn *RaftNode[DataType]) Exists(key string) bool {
	return rn.Get(key) != nil
}

func (rn *RaftNode[DataType]) IsRunning() bool {
	return rn.running
}

func (rn *RaftNode[DataType]) Start(config RaftNodeConfig) error {

	if rn.running {
		return errors.New("node is already running")
	}

	var (
		opts      []raft.Option
		startOpts []raft.StartOption
	)

	addr := config.NodeAddr.String()
	api := config.ApiAddr.String()
	router := mux.NewRouter()
	state := "/tmp/raft" //TODO: add to configuration with default value

	startOpts = append(startOpts, raft.WithAddress(addr))
	opts = append(opts, raft.WithStateDIR(state))

	if config.Join != nil {
		opt := raft.WithFallback(
			raft.WithJoin(config.Join.String(), time.Second),
			raft.WithRestart(),
		)
		startOpts = append(startOpts, opt)
	} else {
		opt := raft.WithFallback(
			raft.WithInitCluster(),
			raft.WithRestart(),
		)
		startOpts = append(startOpts, opt)
	}

	raftgrpc.Register(
		raftgrpc.WithDialOptions(grpc.WithInsecure()),
	)
	rn.raftnode = raft.NewNode(rn.node, transport.GRPC, opts...)
	server := grpc.NewServer()
	raftgrpc.RegisterHandler(server, rn.raftnode.Handler())

	go func() {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}

		err = server.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := rn.raftnode.Start(startOpts...)
		if err != nil && err != raft.ErrNodeStopped {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(api, router); err != nil {
		log.Fatal(err)
	}
	rn.running = true

	return nil
}
