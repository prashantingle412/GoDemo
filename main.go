package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var err error

type Booking struct {
	Id      int    `json:"id"`
	User    string `json:"user"`
	Members int    `json:"members"`
	City    string `json:"city"`
}
type User struct {
	gorm.Model
	Name string `json:"name"`
}
type Profile struct {
	gorm.Model
	User User
	Type string `json:"type"`
}

func AddProfile(w http.ResponseWriter, r *http.Request) {
	db.SingularTable(true)
	db.AutoMigrate(&Profile{})
	var profile Profile
	reqBodyErr := json.NewDecoder(r.Body).Decode(&profile)
	if reqBodyErr != nil {
		fmt.Println("error in decoding ", reqBodyErr)
	}
	db.Create(&Profile{Type: profile.Type})
	json.NewEncoder(w).Encode(profile)

}
func createNewBooking(w http.ResponseWriter, r *http.Request) {
	var booking Booking
	reqBodyErr := json.NewDecoder(r.Body).Decode(&booking)
	if reqBodyErr != nil {
		fmt.Println("error in decoding ", reqBodyErr)
	}
	db.SingularTable(true)
	db.Create(&Booking{Id: booking.Id, User: booking.User, Members: booking.Members, City: booking.City})
	fmt.Println("Endpoint Hit: Creating New Booking")
	json.NewEncoder(w).Encode(booking)
}
func GetAll(w http.ResponseWriter, r *http.Request) {
	var booking []Booking
	db.SingularTable(true)
	db.Find(&booking)
	json.NewEncoder(w).Encode(booking)
}
func Update(w http.ResponseWriter, r *http.Request) {
	var booking Booking
	db.SingularTable(true)
	reqBodyErr := json.NewDecoder(r.Body).Decode(&booking)
	if reqBodyErr != nil {
		fmt.Println("error in decoding ", reqBodyErr)
	}
	db.Model(&booking).Update(map[string]interface{}{"user": "newone", "members": 4})
	db.Model(&booking).Update("user", booking.User)
	json.NewEncoder(w).Encode(booking)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	db.SingularTable(true)
	db.AutoMigrate(&User{})
	var user User
	reqBodyErr := json.NewDecoder(r.Body).Decode(&user)
	if reqBodyErr != nil {
		fmt.Println("error in decoding json array", reqBodyErr)
	}
	db.Create(&User{Name: user.Name})
	json.NewEncoder(w).Encode(user)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	e := os.Getenv("ENV_NAME")
	db, err = gorm.Open(e, "root:root@tcp(127.0.0.1:3306)/Football?charset=utf8&parseTime=True")
	// db,err = gorm.Open(env_name,)
	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}
	r := mux.NewRouter()
	r.HandleFunc("/", GetAll)
	r.HandleFunc("/edit", Update)
	r.HandleFunc("/new", createNewBooking)
	r.HandleFunc("/createuser", CreateUser)
	r.HandleFunc("/addprofile", AddProfile)
	http.ListenAndServe(":4447", r)
}
