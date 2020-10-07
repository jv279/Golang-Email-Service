package main
	import (
		"encoding/json"
		"github.com/gorilla/mux"
		"log"
		"net/http"
	)
// Init Variables
	var addressBook map[string]string
	var serverPort string
// Handle Requests 
	func handleRequests() {
		router := mux.NewRouter().StrictSlash( true )
		router.HandleFunc( "/Bluebook/{address}", lookUp ).Methods( "GET" )
		log.Fatal( http.ListenAndServe( serverPort, router ) )
	}

// Create the addressbook with key names and value of the port 
	func createBook(){
		addressBook = make(map[string]string)
		addressBook["Bob"] = "8888"
		addressBook["Billy"] = "8888"
		
	}

//Lookup sends the port back as a response
	func lookUp( w http.ResponseWriter, r *http.Request ) {
			vars := mux.Vars( r )
			address := vars[ "address" ]
			if _, ok := addressBook[ address ]; ok {
				w.WriteHeader( http.StatusOK )
				if enc, err := json.Marshal( addressBook[ address ] ); err == nil {
					w.Write( []byte( enc ) )
				} else {
					w.WriteHeader( http.StatusInternalServerError )
					}
			} else {
				w.WriteHeader( http.StatusNotFound )
				}
		}

	func main() {
		serverPort = ":9090"
		createBook()
		handleRequests()
	}
