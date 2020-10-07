package main
	import (
		"encoding/json"
		"github.com/gorilla/mux"
		"log"
		"bytes"
		"fmt"
		"net/http"
		"io/ioutil"
		"time"
	)

// Structures
	type Email struct {
		To string
		From string
		Body string
	}
// Init Variables
	var outboxListUsers []string
	var serverPort string
// Handle Requests 
	func handleRequests() {
		router := mux.NewRouter().StrictSlash( true )
		router.HandleFunc( "/sendUser/{user}", addUser ).Methods( "POST" )
		log.Fatal( http.ListenAndServe( serverPort, router ) )
	}

// adds the user send from the MSA to the list outboxListUsers.
	func addUser( w http.ResponseWriter, r *http.Request ) {
		vars := mux.Vars( r )
		user := vars[ "user" ]
		outboxListUsers = append(outboxListUsers , user)
		fmt.Println( "added user ", user, " to ", outboxListUsers)
		
	}
//Runs in the background for a ticker in order to check the outboxListUsers periodically 
	func backgroundTask() {
		//a tick occurs every ten seconds
		ticker := time.NewTicker(10 * time.Second)
		for _ = range ticker.C {

			if (len(outboxListUsers) > 0) {
				//checkOutbox is called if theres something in the outboxListUsers
				fmt.Println( "Outboxes Pending: ", outboxListUsers)
				checkOutbox(outboxListUsers)
				//the first user is removed like a queue for the next tick
				outboxListUsers = outboxListUsers[1:]

			}else{
				fmt.Println( "Outboxes Empty")
			}
		}
	}
	
// first email in the outbox of the user in the front of the list is sent 
	func checkOutbox(outboxListUsers []string){
		//first user's outbox in the list is checked 
		fromUser := outboxListUsers[0]
		//lookup uses the bluebook to return the port of the users address (in this case address is "bob" or "billy")
		fromPort := lookUp(fromUser)
		//URL is created for the request to get the email
		fromUrl := "http://192.168.1.8:"+fromPort+"/MSA/Outbox/"+ fromUser +"/0"
		email := getEmail(fromUrl)
		//The email is then sent and deleted from the outbox
		sendEmail(email)
		deleteEmail(fromUrl)
		fmt.Println( "Sent and Deleted Email From ", fromUser  )
		
		}



// Looks up the port of the user/address in the bluebook
	func lookUp(user string) string{
		client := &http.Client {}
	//port is used in urls so needs to be in string format e.g. "8888" and not 8888
		var port string
		var jsonBlob []byte
		var url string
		url = "http://192.168.1.6:9090/Bluebook/"+user
		if req, err1 := http.NewRequest( "GET", url, nil );
			err1 == nil {
			if resp, err2 := client.Do( req );
				err2 == nil {
				if body, err3 := ioutil.ReadAll( resp.Body );
					err3 == nil {
					jsonBlob = body
					// unmarshaling of the json body into port format
					err := json.Unmarshal(jsonBlob, &port)	
					if err != nil {
						fmt.Println("error:", err)
					}
				} else {
				fmt.Println( "POST failed with %s\n", err3 )
			}
			} else {
			fmt.Println( "POST failed with %s\n", err2 )
			}
		} else {
		fmt.Println( "POST failed with %s\n", err1 )
		}
		return port
	}
	
// Retrieves the email from the MSA outbox and returns it
	func getEmail(url string) Email{
		client := &http.Client {}
		var email Email
		var jsonBlob []byte

		if req, err1 := http.NewRequest( "GET", url, nil );
			err1 == nil {
				
			if resp, err2 := client.Do( req );
				err2 == nil {

				if body, err3 := ioutil.ReadAll( resp.Body );
					err3 == nil {
					jsonBlob = body
					
					err := json.Unmarshal(jsonBlob, &email)	
					if err != nil {
						fmt.Println("error:", err)
					}
				} else {
				fmt.Println( "POST failed with %s\n", err3 )
			}

			} else {
			fmt.Println( "POST failed with %s\n", err2 )
			}

		} else {
		fmt.Println( "POST failed with %s\n", err1 )
		}
		return email
	}

// Sends the email to the MSA at the right address
	func sendEmail(email Email){
	//Building the URL from the email given 
		to:= email.To
		port:= lookUp(to)
		url := "http://192.168.1.8:"+port+"/MSA/Inbox/"+to
		client := &http.Client {}		
		if enc, err := json.Marshal( email );
			err == nil {
			if req, err1 := http.NewRequest( "POST", url, bytes.NewBuffer( enc ) );
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
		} else {
		fmt.Println( "POST failed with %s\n", err )
		}
	}


// Deletes the email that has been send, the URL will end in /0 as its the first email in the outbox that has been sent
	func deleteEmail(url string){
		
		client := &http.Client {}
		
		if req, err := http.NewRequest( "DELETE", url, nil );
			err == nil {
			if _, err1 := client.Do( req );
				err1 == nil {
				// nothing
				
			} else {
				fmt.Println( "DELETE failed with %s\n", err1 )
			}
		} else {
			fmt.Println( "DELETE failed with %s\n", err )
		
		}
		

	}

	func main() {
		go backgroundTask()
		serverPort = ":8989"
		handleRequests()
	}
