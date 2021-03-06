package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"src/data"
	"strconv"
	"strings"
	"time"
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
var userdb = map[string]string{
	"user1": "password123",
}

//assign the secret key to key variable on program's first run
func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key = []byte(os.Getenv("SECRET_KEY"))
}

func signup(w http.ResponseWriter, r *http.Request) {
	//create a Credentials object
	var creds Credential
	//decode json to struct
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userPassword, ok := userdb[creds.Username]

	//if user exists, verify the password
	if ok || userPassword == creds.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Account already exists for this user"))
		return
	}

	userdb[creds.Username] = creds.Password

	json.NewEncoder(w).Encode(fmt.Sprintf("%s successfully signed up. Please login", creds.Username))
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
		w.Write([]byte("Unrecognized credentials. Must sign up first"))
		return
	}

	//Create a token object and add the Username and StandardClaims
	var tokenClaim = Token{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			//Enter the expiration in milliseconds
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	//Create a new claim with HS256 algorithm and token claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)

	tokenString, err := token.SignedString(key)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(tokenString)
}

//dashboard User's personalized dashboard
func dashboard(w http.ResponseWriter, r *http.Request) {
	//get the bearer token from the request handler

	bearerToken := r.Header.Get("Authorization")

	//validate token, it will return Token and error
	token, err := ValidateToken(bearerToken)

	if err != nil {
		//check if Error is Signature Invalid Error

		if err == jwt.ErrSignatureInvalid {
			//return the Unauthorized Status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Return the bad request for any other error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		//return the Unauthorized status for expired token
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Type cast the claims to *Token type
	user := token.Claims.(*Token)

	//send the username Dashboard message
	json.NewEncoder(w).Encode(fmt.Sprintf("%s Dashboard", user.Username))
}

// ValidateToken validates the token with the secret key and returns the object
func ValidateToken(bearerToken string) (*jwt.Token, error) {

	//format the token string
	tokenString := strings.Split(bearerToken, " ")[1]

	//Parse the token with tokenObj
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	//return token and err
	return token, err
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Alex")
}

func returnAllBoards(w http.ResponseWriter, r *http.Request) {
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
	myRouter.HandleFunc("/signup", signup).Methods("POST")
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
