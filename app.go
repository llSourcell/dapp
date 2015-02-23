
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
	"strings"
	"os"
	"bytes"
)
var gnode *core.IpfsNode

type DemoPage struct {
    Title  string
    Author string
	Tweet string 
}


//submit page
func SubmitTextPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
	//pull tweet from IPFS and display
 reader, err := coreunix.Cat(gnode, "QmdWDjrgPmnDdu4UZeunmtp5Aa8i6HpxpeCi3Tor4BNaAH")
 buf := new(bytes.Buffer)
 buf.ReadFrom(reader)
 s := buf.String()
 fmt.Printf("here it is %s", s)
 
 
    demoheader := DemoPage{"HTTP -> IPFS Add demo", "Siraj", s}

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

//posts submit data
func textToIPFS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//get user input text 
    r.ParseForm()
	fmt.Println("input text is:", r.Form["sometext"])
	var userInput = r.Form["sometext"]
	
	//write it to output.html
    err := ioutil.WriteFile("output.html", []byte(userInput[0]) , 0777)
	
	
	//read output.html
 		f,err := ioutil.ReadFile("output.html")
 		if err!=nil {
 			fmt.Print( "err in ioutil.ReadFile" )
 		}
 		node := strings.NewReader(string(f))
	

		//add output.html to IPFS 
     str, err := coreunix.Add(gnode, node)
     if err != nil {
         w.WriteHeader(504)
         w.Write([]byte(err.Error()))
     } else {
		 //IPFS hash
         w.Write([]byte(str))
     }
	 
	 //delete output.html locally
	  err2 := os.Remove("output.html")
	  if err2 != nil {
	  	
		  panic(err2)
	  }
}


func main() {
	
	//init IPFS node
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
	
	//router 
    router := httprouter.New()
    router.GET("/", SubmitTextPage)
    router.POST("/submittext", textToIPFS)
	
    log.Fatal(http.ListenAndServe(":8080", router))
}

	 
	 
	 
