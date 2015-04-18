package main

import (
    "net/http"
	"github.com/julienschmidt/httprouter"
    "github.com/jbenet/go-ipfs/core"
	"path"
	"html/template"
	"fmt"
	"log"
	"github.com/go-kerala/kerala"
	peer "github.com/jbenet/go-ipfs/p2p/peer"

	
	
	
)

type IPFSHandler struct {
	node  *core.IpfsNode
}

type DemoPage struct {
    Title  string
    Author string
	Tweet []string 
	isMine bool
	Balance float64
}

type PeerList struct {
	Allpeers []string
	
}

func main() {
	
	node, err := kerala.StartNode()
	if err != nil {
		panic(err)
	}
	

	//[2] Define routes 
    router := httprouter.New()
	//Route 1 Home (profile)
    router.GET("/", TextInput(node))	
	//Route 2 Discover page
	router.GET("/discover", displayUsers(node))
	//Route 3 Other user profiles
    router.GET("/profile/:name", TextInput(node))	
	//Route 4 Add text to IPFS
    router.POST("/textsubmitted", addTexttoIPFS(node))
	router.GET("/test", Test())
	
    //[3] link resources
	router.ServeFiles("/resources/*filepath", http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.Handle("/", router)
	
	//[4] Start server
	fmt.Println("serving at 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
	
	
	
	
	
}


func displayUsers(node *core.IpfsNode) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		
		//get peers
		
		peers := node.Peerstore.Peers()
		data := make([]string, len(peers))
		for i := range data {
		    data[i] = peer.IDB58Encode(peers[i])
		}
		fmt.Println("the peers are %s", data)
		
		
		
		
	    demoheader := PeerList{data}
        fp := path.Join("templates", "discover.html")
        tmpl, err := template.ParseFiles(fp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if err := tmpl.Execute(w, demoheader); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
		
	}
}


func Test() httprouter.Handle {

return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
	hash := kerala.Pay("1000","1HihKUXo6UEjJzm4DZ9oQFPu2uVc9YK9Wh", "akSjSW57xhGp86K6JFXXroACfRCw7SPv637", "10", "AHthB6AQHaSS9VffkfMqTKTxVV43Dgst36", "L1jftH241t2rhQSTrru9Vd2QumX4VuGsPhVfSPvibc4TYU4aGdaa" )
	fmt.Println(hash)

		
}
}


func TextInput(node *core.IpfsNode) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {		

		
		
		fmt.Println("textinput is even called %s", ps.ByName("name"))
		var userID = ps.ByName("name")
		balance := kerala.GetBalance("1HihKUXo6UEjJzm4DZ9oQFPu2uVc9YK9Wh")
		
		
		
		//if its the homepage, just pull the node ID from output.html
		if userID == "" {
			tweetArray, err := kerala.GetStrings(node, "")
			if err != nil {
				panic(err)
			}
			fmt.Println("the tweet array is %s", tweetArray)
			if tweetArray == nil {
				fmt.Println("tweetarray is nil")
			    demoheader := DemoPage{"Decentralized Twitter", "Siraj", nil, true, balance }
		        fp := path.Join("templates", "index.html")
		        tmpl, err := template.ParseFiles(fp)
		        if err != nil {
		            http.Error(w, err.Error(), http.StatusInternalServerError)
		            return
		        }
		        if err := tmpl.Execute(w, demoheader); err != nil {
		            http.Error(w, err.Error(), http.StatusInternalServerError)
		        }
			} else {
				fmt.Println("tweetarray is not nil")
		
		    demoheader := DemoPage{"Decentralized Twitter", "Siraj", tweetArray, true, balance}
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
		  
		} else {
			
			//else pull it from the URL
			tweetArray, err := kerala.GetStrings(node, userID)
			if err != nil {
				panic(err)
			}
			fmt.Println("the tweet array is %s", tweetArray)
			if tweetArray == nil {
				fmt.Println("tweetarray is nil")
			    demoheader := DemoPage{"Decentralized Twitter", "Siraj", nil, false,balance}
		        fp := path.Join("templates", "index.html")
		        tmpl, err := template.ParseFiles(fp)
		        if err != nil {
		            http.Error(w, err.Error(), http.StatusInternalServerError)
		            return
		        }
		        if err := tmpl.Execute(w, demoheader); err != nil {
		            http.Error(w, err.Error(), http.StatusInternalServerError)
		        }
			} else {
				fmt.Println("tweetarray is not nil")
		
		    demoheader := DemoPage{"Decentralized Twitter", "Siraj", tweetArray, false,balance}
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
			 
			
		}
	

	 
	
}
	 
	
}

func addTexttoIPFS(node *core.IpfsNode) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    r.ParseForm()
	fmt.Println("input text is:", r.Form["sometext"])
	var userInput = r.Form["sometext"]
    Key, err := kerala.AddString(node, userInput[0])
	if err != nil {
		panic(err)
	}
	
	fmt.Println("the key", Key)
   // w.Write([]byte(Key))
	
	
	
}
}











 
	 
	 
