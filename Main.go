package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var ServerPeerMap = []ID{MINERID, MINER2ID}
var SID = ID{Address: "http://localhost:", Port: "6688"}

//SRD is a server register data global for a server instance
var SRD = ServerRegisterData{ServerID: SID, PeerMapJSON: "", PeerMap: ServerPeerMap}

// Bc ...
var Bc = NewBlockChain()
var mutex = &sync.Mutex{}

// PEERLIST is for the miner, so initialized with empty peer IDs
var PEERLIST = PeerList{SelfID: MINERID, PeerIDs: []ID{}, Length: 0}

// port options 8080 6689

var MINERID = ID{Address: "http://localhost:", Port: "6001"}

var MINER2ID = ID{Address: "http://localhost:", Port: "8001"}

//var target = "000000" // six 0 fairly quick

// var SELFID = []string{"http://localhost:", "6689"}

//var PEERID = []string{"http://localhost:8080"}
var target = "0000000" // seven 0 ... long time

// SELFADDRESS ...
//var SELFADDRESS = "http://localhost:" + SELFID[1]

func main() {
	miner1Setup()
	// registationServerSetup()
}
func miner1Setup() {

	genesis := makeGenesisBlock()
	Bc.Insert(genesis)
	// fmt.Println(Bc.Show())
	//testAsk()
	fmt.Println("LOG: I am a miner")
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
	fmt.Println("Log: registration server has been started up at: ", SID.Port)
	fmt.Println("Log: my  peermapjson is: ", SRD.PeerMap)
	router := NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":"+SID.Port, router))
	}

}
