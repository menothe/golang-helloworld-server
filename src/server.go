package main

import (
    "fmt"
    "net/http"
    "os"
    "encoding/json"
    "src/data"
    "log"
    "strconv"
    "io/ioutil"
    "github.com/gorilla/mux"
)

func helloWorld(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hello World")
}

func returnAllBoards(w http.ResponseWriter, r *http.Request){
  fmt.Println("Endpoint Hit: returnAllArticles")
  json.NewEncoder(w).Encode(data.Skateboards)
}

func returnSingleBoard(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: returnSingleBoard")

  vars := mux.Vars(r)
  key := vars["id"]

  for _, board := range data.Skateboards {
    if strconv.Itoa(board.Id) == key {
        json.NewEncoder(w).Encode(board)
    }
  }
}

func createNewBoard(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: createNewBoard")

  // get the body of our POST request
  // unmarshal this into a new Skateboard struct
  // append this to our Skateboards array.
  reqBody, _ := ioutil.ReadAll(r.Body)
  var board data.Skateboard
  json.Unmarshal(reqBody, &board)
  // update our global Skateboards array to include
  // our new Skateboard
  data.Skateboards = append(data.Skateboards, board)
  json.NewEncoder(w).Encode(board)
}

func deleteBoard(w http.ResponseWriter, r *http.Request) {
    // once again, we will need to parse the path parameters
    vars := mux.Vars(r)
    // we will need to extract the `id` of the article we
    // wish to delete
    id := vars["id"]

    reqBody, _ := ioutil.ReadAll(r.Body)
    var board data.Skateboard
    json.Unmarshal(reqBody, &board)
    // we then need to loop through all our articles
    for index, board := range data.Skateboards {
        // if our id path parameter matches one of our
        // articles
        if strconv.Itoa(board.Id) == id {
            // updates our Articles array to remove the
            // article
            data.Skateboards = append(data.Skateboards[:index], data.Skateboards[index+1:]...)
        }
    }
}

func updateBoard(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]

  reqBody, _ := ioutil.ReadAll(r.Body)
  var updatedBoard data.Skateboard
  json.Unmarshal(reqBody, &updatedBoard)

  for index, board := range data.Skateboards {

    if strconv.Itoa(board.Id) == id {
      data.Skateboards[index] = updatedBoard
    }
  }
}

func handleRequests() {
  myRouter := mux.NewRouter().StrictSlash(true)

  myRouter.HandleFunc("/", helloWorld)
  myRouter.HandleFunc("/skateboards", returnAllBoards)
  myRouter.HandleFunc("/skateboard", createNewBoard).Methods("POST")
  myRouter.HandleFunc("/skateboard/{id}", deleteBoard).Methods("DELETE")
  myRouter.HandleFunc("/skateboard/{id}", updateBoard).Methods("PUT")
  myRouter.HandleFunc("/skateboard/{id}", returnSingleBoard)

  log.Fatal(http.ListenAndServe(getPort(), myRouter))
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func main() {
  handleRequests()
}
