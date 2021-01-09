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
    jwt "github.com/dgrijalva/jwt-go"
    "github.com/joho/godotenv"
    "time"
    "strings"
)

//Secret key to uniquely sign the token
var key []byte

//Credential User's login information
type Credential struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

//Token jwt standard claim object
type Token struct {
  Username string `json:"username"`
  jwt.StandardClaims
}

//dummy local db instance as a key value pair
var userdb = map[string]string {
  "user1": "password123",
}

//assign the secret key to key variable on program's first run
func init() {
  err := godotenv.Load()

  if err != nil {
    log.Fatal("Error loading .env file")
  }

  key = []byte(os.Getenv("SECRET_KEY"))
}


//login user login function
func login(w http.ResponseWriter, r *http.Request) {
  //create a Credentials object
  var creds Credential
  //decode json to struct
  err := json.NewDecoder(r.Body).Decode(&creds)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  //verify if user exists or not
  userPassword, ok := userdb[creds.Username]

  //if user exists, verify the password
  if !ok || userPassword != creds.Password {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  //Create a token object and add the Username and StandardClaims
  var tokenClaim = Token {
    Username: creds.Username,
    StandardClaims: jwt.StandardClaims{
      //Enter the expiration in milliseconds
      ExpiresAt: time.Now.Add(10 * time.Minute).Unix(),
    },
  }

  //Create a new claim with HS256 algorithm and token claim
  token := jwt.NewWithClaims(jwt.SigningMethodHs256, tokenClaim)

  tokenString, err := token.SignedString(key)

  if err != nil {
    log.Fatal(err)
  }

  json.NewEncoder(w).Encode(tokenString)
}

func helloWorld(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hello Alex")
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

  //main CRUD endpoints
  myRouter.HandleFunc("/", helloWorld)
  myRouter.HandleFunc("/skateboards", returnAllBoards)
  myRouter.HandleFunc("/skateboard", createNewBoard).Methods("POST")
  myRouter.HandleFunc("/skateboard/{id}", deleteBoard).Methods("DELETE")
  myRouter.HandleFunc("/skateboard/{id}", updateBoard).Methods("PUT")
  myRouter.HandleFunc("/skateboard/{id}", returnSingleBoard)

  //user authentication endpoints
  myRouter.HandleFunc("/login", login).Methods("POST")
  myRouter.HandleFunc("/dashboard", dashboard).Methods("GET")

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
