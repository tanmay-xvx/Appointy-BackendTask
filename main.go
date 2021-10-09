package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)


type User struct {
	Id uuid.UUID `json:"id" bson:"_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	Id string `json:"id"`
	Caption string `json:"caption"`
	ImageUrl string `json:"imageurl"`
	Timestamp string `json:"timestamp"`
}

type userHandlers struct {
	sync.Mutex
	store map[string]User
}

func (h *userHandlers) users(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case "GET":
			h.get(w,r)
			return
		case "POST":
			h.post(w,r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("not implemented"))
			return
	}
}

func (h *userHandlers) get(w http.ResponseWriter, r *http.Request){
	users := make([]User, len(h.store))
	i:=0
	for _, user := range h.store{
		users[i] = user
	}

	jsonbytes, err := json.Marshal(users);
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonbytes)

	// fmt.Fprint(w,users)
}

func insertUser(u User){
	
	
	clientOptions := options.Client().ApplyURI("mongodb+srv://tanmay:123@cluster0.zx333.mongodb.net/instaClone?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
    	log.Fatal(err)
	}
	ctxm, _ := context.WithTimeout(context.Background(), 15*time.Second)

	col := client.Database("instaClone").Collection("Users")

	u.Id,_ = uuid.New()
	result, insertErr := col.InsertOne(ctxm, u)

	println(result,insertErr)
	println(u.Name,u.Email,u.Password)
	
}

func (h *userHandlers) post(w http.ResponseWriter, r *http.Request){
	bodybytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var usr User
	err = json.Unmarshal(bodybytes, &usr)
	insertUser(usr)
}

func newUserHandlers() *userHandlers {
	uid,_ := uuid.New()
	return &userHandlers{
		store:map[string]User{
			"1":User{
				Id:uid,
				Name:"John Smith",
				Email:"John Smith",
				Password:"123",
			},
		},
	}
}


func handleRequests(){
	uh := newUserHandlers()
	http.HandleFunc("/users",uh.users)
	log.Fatal(http.ListenAndServe(":1000",nil))
}



func main() {
	// insertUser()
	handleRequests()
}