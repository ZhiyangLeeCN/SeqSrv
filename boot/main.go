package main

import (
	"flag"
	"fmt"
	"github.com/ZhiyangLeeCN/SeqSrv/gen"
	"net/http"
	"os"
	"strconv"
)

// Flag Names
var (
	nmHelp       = "help"
	nmAddr       = "bind-addr"
	nmNodeId     = "node-id"
	nmSeqFile    = "seq-file"
	nmInitCurSeq = "init-cur-seq"
	nmInitMaxSeq = "init-max-seq"
	nmStep       = "step"
)

var (
	help       = flagBoolean(nmHelp, false, "print help")
	addr       = flag.String(nmAddr, "0.0.0.0:8181", "server listen addr.")
	nodeId     = flag.Uint64(nmNodeId, 1, "seq service node id.")
	seqFile    = flag.String(nmSeqFile, "max_seq", "seq persist file path.")
	initCurSeq = flag.Uint64(nmInitCurSeq, 0, "init seq value.")
	initMaxSeq = flag.Uint64(nmInitMaxSeq, 1000, "init max seq value.")
	step       = flag.Uint64(nmStep, 1000, "max seq inc step value.")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	node := initSeqNode(*seqFile, *nodeId, *initCurSeq, *initMaxSeq, *step)
	runNodeServer(*addr, node)
}

func initSeqNode(seqFile string, nodeId uint64, initCurSeq uint64, initMaxSeq uint64, step uint64) *gen.SeqNode {
	if seqFile == "" {
		fmt.Println("seq-file is required!")
		os.Exit(-1)
	}

	sf, err := gen.NewSeqFile(seqFile, initCurSeq, initMaxSeq, step)
	if err != nil {
		fmt.Println("NewSeqFile error.", err)
		os.Exit(-1)
	}

	node, err := gen.NewSeqNode(sf, nodeId)
	if err != nil {
		fmt.Println("NewSeqNode error.", err)
		os.Exit(-1)
	}

	return node
}

func runNodeServer(addr string, node *gen.SeqNode) {
	http.HandleFunc("/seq/next", func(writer http.ResponseWriter, request *http.Request) {
		nextId, err := node.NextID()
		if err != nil {
			writer.WriteHeader(500)
			_, _ = fmt.Fprintf(writer, "gen error.")
		}

		nextIdStr := strconv.FormatUint(nextId, 10)
		writer.WriteHeader(200)
		_, _ = fmt.Fprintf(writer, nextIdStr)
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("seq service listen error.", err)
		os.Exit(-1)
	}
}

func flagBoolean(name string, defaultVal bool, usage string) *bool {
	if !defaultVal {
		usage = fmt.Sprintf("%s (default false)", usage)
		return flag.Bool(name, defaultVal, usage)
	}
	return flag.Bool(name, defaultVal, usage)
}
