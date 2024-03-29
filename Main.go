package main

import (
	"log"
	"net/http"
	"os"
)

// MINERID is globabl for miner ID - ideally port will become OS.Arg[1] when launched
var MINERID = makeMinerID()

// ServerPeerMap Initializing with two hardcoded minors
// var ServerPeerMap = []ID{MINERID, miner2id}
var ServerPeerMap = []ID{}

// SID is the server id, which is always at 6688
var SID = ID{Address: "http://localhost:", Port: "6688"}

// SRD is a server register data global for a server instance
var SRD = ServerRegisterData{ServerID: SID, PeerMapJSON: "[]", PeerMap: ServerPeerMap}

// SYNCBC global sync block chain instance containing *BC
var SYNCBC = NewSyncBlockChain()

// PEERLIST is for the miner, so initialized with empty peer IDs
var PEERLIST = PeerList{SelfID: makeMinerID(), PeerIDs: []ID{SID}, Length: 0}

//var target = "000000" // six 0 fairly quick

// TARGET used in block_generation so the hash output must beging with at least target 0's
var TARGET = "00000" // five 0, very quick

func main() {
	if len(os.Args) > 1 {
		minerSetup()
	} else {
		registationServerSetup()
	}

}

func testAsk() {
	askForParent("97611bc0a6e098f600d4c776252ffc16173058b4d5e2ae4a7d336fe18eb7f11326b2ca4e40be4c5572800ca76c6cdf4b65931b297f098f73a256f497c8907736", "1")
}

func makeMinerID() ID {
	if len(os.Args) > 1 {
		return ID{Port: os.Args[1], Address: "http://localhost:"}
	}
	return SID // ID{Port: "6688", Address: "http://localhost:"}
}

func minerSetup() {

	log.Println("LOG: I am a miner")
	log.Println("LOG: My peerlist prior to registration is: ", PEERLIST.PeerIDs)
	router := NewRouter()
	// if len(os.Args) > 1 {
	log.Fatal(http.ListenAndServe(":"+MINERID.Port, router))
	// 	} else {
	// 		log.Fatal(http.ListenAndServe(":"+MINERID.Port, router))
	// 	}
}

func registationServerSetup() {
	genesis := makeGenesisBlock()
	SYNCBC.Insert(genesis)
	SRD.EncodePeerMapToJSON()
	log.Println("Log: this is a registration Node has started up at: ", SID.Port)
	log.Println("Log: my  peermapjson is: ", SRD.PeerMap)
	router := NewRouter()
	// if len(os.Args) > 1 {
	// 	log.Fatal(http.ListenAndServe(":"+MINERID.Port, router))
	// } else {
	log.Fatal(http.ListenAndServe(":"+SID.Port, router))
	// }
}
