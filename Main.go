package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

// ServerPeerMap Initializing with two hardcoded minors
var ServerPeerMap = []ID{MINERID, miner2id}

// SID is the server id, which is always at 6688
var SID = ID{Address: "http://localhost:", Port: "6688"}

//SRD is a server register data global for a server instance
var SRD = ServerRegisterData{ServerID: SID, PeerMapJSON: "", PeerMap: ServerPeerMap}

// Bc ...
var Bc = NewBlockChain()
var mutex = &sync.Mutex{}

// PEERLIST is for the miner, so initialized with empty peer IDs
var PEERLIST = PeerList{SelfID: MINERID, PeerIDs: []ID{}, Length: 0}

// MINERID is globabl for miner ID - port will become OS.Arg[1] when launched
var MINERID = ID{Address: "http://localhost:", Port: "6001"}
var miner2id = ID{Address: "http://localhost:", Port: "8001"}

//var target = "000000" // six 0 fairly quick
var target = "0000000" // seven 0 ... long time

func main() {
	if len(os.Args) > 1 {
		minerSetup()
	} else {
		registationServerSetup()
	}
}

func minerSetup() {
	genesis := makeGenesisBlock()
	Bc.Insert(genesis)
	log.Println("LOG: I am a miner")
	fmt.Println("LOG: My peerlist prior to registration is: ", PEERLIST.PeerIDs)
	MINERID.Port = os.Args[1]
	router := NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":"+MINERID.Port, router))
	}
}

func registationServerSetup() {
	SRD.EncodePeerMapToJSON()
	fmt.Println("Log: this is a registration server has started up at: ", SID.Port)
	fmt.Println("Log: my  peermapjson is: ", SRD.PeerMap)
	router := NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":"+SID.Port, router))
	}

}
