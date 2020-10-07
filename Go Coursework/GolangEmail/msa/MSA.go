package main
	import (
		"github.com/gorilla/mux"
		"log"
		"net/http"
		"encoding/json"
		"strconv"
		"fmt"
		"io/ioutil"
	)

// Structures
	type Email struct {
		To string
		From string
		Body string
	}

	type Sentmail struct {
		To string
		Body string
	}

// Init Variables
	var outbox map[ string ] []Email
	var inbox map[ string ] []Email
	var serverPort string

// Handle Requests
	func handleRequests() {
		router := mux.NewRouter().StrictSlash( true )
		router.HandleFunc( "/MSA/Outbox/{user}", SendToOutbox ).Methods( "POST" )
		router.HandleFunc( "/MSA/Inbox/{user}", SendToInbox ).Methods( "POST" )
		router.HandleFunc( "/MSA/{box}/{user}/{ID}", ReadEmails ).Methods( "GET" )
		router.HandleFunc( "/MSA/{box}/{user}/{ID}", DeleteFromBox ).Methods( "DELETE" )
		router.HandleFunc( "/MSA/{box}/{user}", ListEmails ).Methods( "GET" )

		log.Fatal( http.ListenAndServe( serverPort, router ) )
	}
	

// Returns the variable inbox or outbox based on the string arguement
	func inOrOut(box string) map[ string ] []Email{
		if box == "Inbox"{
			return inbox
		}
		if box == "Outbox"{
			return outbox
		}
		return nil
	}

// Lists all the emails in whichever box is specified
	func ListEmails( w http.ResponseWriter, r *http.Request ) {
		vars := mux.Vars( r )
		user := vars[ "user" ]
		box := inOrOut(vars["box"])		
		if email, ok := box[ user ]; ok {
			w.WriteHeader( http.StatusOK )
			if len(email) !=0{
				if enc, err := json.Marshal( email ); err == nil {
					w.Write( []byte( enc ) )
				}else {
					w.WriteHeader( http.StatusInternalServerError )
				}
			}else{
				fmt.Println("Nothing in ", vars["box"], " for user :", user)
			}
		} else {
			w.WriteHeader( http.StatusNotFound )
		} 
	}


// Reads all emails in the inbox or outbox
	func ReadEmails( w http.ResponseWriter, r *http.Request ) {
	// Get Variables
		vars := mux.Vars( r )
		user := vars[ "user" ]
		box := inOrOut(vars["box"])
	// The ID retrieved is a string and gets converted here into an integer
		StrID := vars[ "ID" ]
		var ID int
		IntID, err := strconv.Atoi(StrID)
		if err != nil {
			fmt.Println(w, "String to Int error %s ",  err)
		} else {
		ID = IntID
		}

	
		if email, ok := box[ user ]; ok {
				w.WriteHeader( http.StatusOK )
				
				if len(email) !=0{
					if (len(box[user])-1)< (ID) { 
						fmt.Println("It is not on the archives, it doesn't exist") 
					}else{
					if enc, err := json.Marshal( email[ID] ); err == nil {
						w.Write( []byte( enc ) )
					} else {
						w.WriteHeader( http.StatusInternalServerError )
					}}
				}else{
					fmt.Println("Nothing in ", vars["box"], " for user :", user)
				}
			} else {
				w.WriteHeader( http.StatusNotFound )
	} }
// Deletes a specified email from the inbox or outbox
		func DeleteFromBox( w http.ResponseWriter, r *http.Request){
			vars := mux.Vars( r )
			user := vars[ "user" ]
			box := inOrOut(vars["box"])
		// The ID retrieved is a string and gets converted here into an integer
			StrID := vars[ "ID" ]
			var ID int
			IntID, err := strconv.Atoi(StrID)
			if err != nil {
				fmt.Println(w, "String to Int error %s ",  err)
			} else {
			ID = IntID
			}

			if email, ok := box[ user ]; ok {
				if len(email) !=0{
				w.WriteHeader( http.StatusNoContent )
				box[ user ] = append(box[ user ][:ID], box[ user ][ID+1:]...)
				}else{
					fmt.Println("Nothing in ", vars["box"], " for user :", user)
				}
			  } else {
				w.WriteHeader( http.StatusNotFound )
			  }      
		}

// Sends incoming Email into the inbox 
		func SendToInbox( w http.ResponseWriter, r *http.Request ) {
			vars := mux.Vars( r )
			user := vars[ "user" ]
			decoder := json.NewDecoder( r.Body )
			var email Email
			if err := decoder.Decode( &email ); err == nil {
				w.WriteHeader( http.StatusCreated )
				inbox[ user ] = append(inbox[ user ], email)
			} else {
				w.WriteHeader( http.StatusBadRequest )
			}
		}

// Puts the request in the outbox after changing it into an "Email" structure from a "Sentmail" structure
		func SendToOutbox( w http.ResponseWriter, r *http.Request ) {
			vars := mux.Vars( r )
			user := vars[ "user" ]
			struser := string(user) 
			decoder := json.NewDecoder( r.Body )
		// sendusers adds the name of the user of the outbox to MTA so that it knows which outbox to check when sending Emails
			sendUsers(struser)
			var sentmail Sentmail 
			var email Email
		
			if err := decoder.Decode( &sentmail ); err == nil {
			// Creating the email structure from the Sentmail Structure
				email = Email{
					To: sentmail.To,
					From : user,
					Body :sentmail.Body,
				}
				w.WriteHeader( http.StatusCreated )
				outbox[ user ] = append(outbox[ user ], email)
		
			} else {
				w.WriteHeader( http.StatusBadRequest )
			}
		}



// Sends the user string to the MTA
		func sendUsers(user string){
			url := "http://192.168.1.7:8989/sendUser/"+user
			client := &http.Client {}
			
				if req, err1 := http.NewRequest( "POST", url, nil );
					err1 == nil {
					if resp, err2 := client.Do( req );
						err2 == nil {
						if _, err3 := ioutil.ReadAll( resp.Body );
							err3 == nil { 
								//nothing
						} else {
						fmt.Println( "POST failed with %s\n", err3 )
						}
					} else {
					fmt.Println( "POST failed with %s\n", err2 )
					}
				} else {
				fmt.Println( "POST failed with %s\n", err1 )
				}

		}


		func main() {
			serverPort = ":8888"
			inbox = make( map[ string ] []Email )
			outbox = make( map[ string ] []Email )
			handleRequests()
		}

