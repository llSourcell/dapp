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


//struct to handle the IPFS core node
type IPFSHandler struct {
	node  *core.IpfsNode
}


//struct to pass vars to front-end
type DemoPage struct {
    Title  string
    Author string
	Tweet []string 
	isMine bool
	Balance float64
}

//struct for list of all peers in IPFS 
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
	
    //[3] link resources
	router.ServeFiles("/resources/*filepath", http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.Handle("/", router)
	
	//[4] Start server
	fmt.Println("serving at 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
	
	
	
	
	
}

//Called on the discovery page. It will create the list of
//all IPFS peers currently online
func displayUsers(node *core.IpfsNode) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		
		//get peers
		peers := node.Peerstore.Peers()
		data := make([]string, len(peers))
		for i := range data {
		    data[i] = peer.IDB58Encode(peers[i])
		}
		fmt.Println("the peers are %s", data)
		
        //send the peer list to the front end template 
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

//Called on all profile pages. Fills the profile page with tweets for the relevant user.
func TextInput(node *core.IpfsNode) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {		

		var userID = ps.ByName("name")
		fmt.Println("BITCH",)
		
		
		//get your current balance
		balance := kerala.GetBalance("1HihKUXo6UEjJzm4DZ9oQFPu2uVc9YK9Wh")
		

		//if its the homepage, just pull the node ID from output.html
		
		//[1] If its your home profile page
		if userID == "" {
			pointsTo, err := node.Namesys.Resolve(node.Context(), node.Identity.Pretty())
			tweetArray, err := kerala.GetStrings(node, pointsTo.B58String())
			if err != nil {
				panic(err)
			}
			fmt.Println("the tweet array is %s", tweetArray)
			//[1A] If no tweets, send nil
			if tweetArray == nil {
				fmt.Println("tweetarray is nil")
			    demoheader := DemoPage{"Decentralized Twitter", "SR", nil, true, balance }
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
				//[1B] If tweets, send tweet array
		    demoheader := DemoPage{"Decentralized Twitter", "SR", tweetArray, true, balance}
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
			
			//[2] If its another profile
			//Pull from IPNS
			pointsTo, err := node.Namesys.Resolve(node.Context(), userID)
			//[2A] If nil, send nil
			if err != nil {
				fmt.Println("ERROR")
				fmt.Println("tweetarray is nil")
			    demoheader := DemoPage{"Decentralized Twitter", "SR", nil, false,balance}
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
				//[2B] Else send tweetarray 
				fmt.Println("RESOLVED")
				//else pull it from the URL
				tweetArray, err := kerala.GetStrings(node, pointsTo.B58String())
				if err != nil {
					panic(err)
				}
				hash := kerala.Pay("1000","1HihKUXo6UEjJzm4DZ9oQFPu2uVc9YK9Wh", "akSjSW57xhGp86K6JFXXroACfRCw7SPv637", "10", "AHthB6AQHaSS9VffkfMqTKTxVV43Dgst36", "L1jftH241t2rhQSTrru9Vd2QumX4VuGsPhVfSPvibc4TYU4aGdaa" )
			 	fmt.Println(hash)
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

//Called when user submits text to IPFS.
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
	
	
	
}
}











 
	 
	 
