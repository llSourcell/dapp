
package main

import (
    "net/http"
	"github.com/julienschmidt/httprouter"
    "github.com/jbenet/go-ipfs/core"
    "github.com/jbenet/go-ipfs/core/coreunix"
    "github.com/jbenet/go-ipfs/repo/fsrepo"
    "code.google.com/p/go.net/context"
	"path"
	"html/template"
	"fmt"
	"log"
	"io/ioutil"
	//"strings"
	//"os"
	"bytes"
	merkledag "github.com/jbenet/go-ipfs/merkledag"
	u "github.com/jbenet/go-ipfs/util"
	
	
)

//Global IPFS node
var gnode *core.IpfsNode

func main() {
	
	//[1] init IPFS node
    builder := core.NewNodeBuilder().Online()
    r := fsrepo.At("~/.go-ipfs")
    if err := r.Open(); err != nil {
       panic(err)
    }
    builder.SetRepo(r)
    // Make our 'master' context and defer cancelling it
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    node, err := builder.Build(ctx)
    if err != nil {
        panic(err)
    }
    gnode = node
    fmt.Printf("I am peer: %s\n", node.Identity)
	
	//[2] Define routes 
    router := httprouter.New()
    router.GET("/", TextInput)
    router.POST("/textsubmitted", addTexttoIPFS)
	
	//[3] Start server
    log.Fatal(http.ListenAndServe(":8080", router))
}

type DemoPage struct {
    Title  string
    Author string
	Tweet string 
}

//[1] input page
func TextInput(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
	//pull tweet from IPFS and display
 reader, err := coreunix.Cat(gnode, "QmdWDjrgPmnDdu4UZeunmtp5Aa8i6HpxpeCi3Tor4BNaAH")
 buf := new(bytes.Buffer)
 buf.ReadFrom(reader)
 s := buf.String()
 fmt.Printf("here it is %s", s)
 
f,err := ioutil.ReadFile("output.html")
var Key = u.B58KeyDecode(string(f))
   var tweetArray = ResolveAllInOrder(gnode,Key)
 
    demoheader := DemoPage{"HTTP -> IPFS Add demo", "Siraj", tweetArray[0]}

     fp := path.Join("templates", "index.html")
     tmpl, err := template.ParseFiles(fp)
     if err != nil {
         http.Error(w, err.Error(), http.StatusInternalServerError)
         return
     }

     if err := tmpl.Execute(w, demoheader); err != nil {
         http.Error(w, err.Error(), http.StatusInternalServerError)
     }
	 
	 
	
	 
	 
	 
}

//[2] text submitted page
func addTexttoIPFS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//[1] Get new tweet
    r.ParseForm()
	fmt.Println("input text is:", r.Form["sometext"])
	var userInput = r.Form["sometext"]
	
	//[2] Check if key is stored locally
	f,err := ioutil.ReadFile("output.html")
 	if err!=nil {
  			fmt.Print( "err in ioutil.ReadFile" )
  		}
	//[3] If key is not stored locally, the user is new
	if string(f) == "" {
		fmt.Printf("local storage is empty")
		//[4] Initialize a MerkleDAG node and key 
		var NewNode * merkledag.Node
		var Key u.Key
		//[5] Fill the node with user input
		NewNode = MakeStringNode(userInput[0])
		//[6] Add the node to IPFS
		Key, _ = gnode.DAG.Add(NewNode)
		//[7] Add the user input key to local store
		fmt.Printf("the key is %s", Key)
		err := ioutil.WriteFile("output.html", []byte(Key.B58String()) , 0777)
		if err != nil {
			panic(err)
		}
		//[8] Write the hash as a response on the page
        w.Write([]byte(Key))
		
		//else the user has already tweeted
		}else {
			//[9] Initialize a new MerkleDAG node and key
			var NewNode * merkledag.Node
			var Key u.Key
			//[10] Fill the node with user input
			NewNode = MakeStringNode(userInput[0])
			//[11] Read the hash from the file
	 		f,err := ioutil.ReadFile("output.html")
			//[12] Convert it into a key
			Key = u.B58KeyDecode(string(f))
			//[13] Get the Old MerkleDAG node and key
			var OldNode * merkledag.Node
			var Key2 u.Key
			OldNode, _ = gnode.DAG.Get(Key)
			//[14]Add a link to the old node
	 		NewNode.AddNodeLink("next", OldNode)
			//[15] Add thew new node to IPFS
			Key2, _ = gnode.DAG.Add(NewNode)
			//[16] Add the new node key to local store (overwrite)
			err2 := ioutil.WriteFile("output.html", []byte(Key2.B58String()) , 0777)
			if err2 != nil {
				panic(err)
			}
			//[17] Write the hash as a response on the page
	        w.Write([]byte(Key2))
		}
		

}

func MakeStringNode(s string) * merkledag.Node {
	n := new(merkledag.Node)
	n.Data = make([]byte, len(s))
	copy(n.Data, s)
	return n;
}

func ResolveAllInOrder(nd * core.IpfsNode, k u.Key) []string {
	var stringArr []string
	var node * merkledag.Node
	node, err := nd.DAG.Get(k)
	if err != nil {
		fmt.Println(err)
		//return
	}

	fmt.Printf("%s ", string(node.Data[:]))
	for ;; {
		var err error

		if len(node.Links) == 0 {
			break;
		}

		node, err = node.Links[0].GetNode(nd.DAG)
		if err != nil {
			fmt.Println(err)
			//return
		}

		fmt.Printf("%s ", string(node.Data[:]))
		stringArr = append(stringArr, string(node.Data[:]))
	}

	fmt.Printf("\n");
	
	return stringArr

}



 
	 
	 
