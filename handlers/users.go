package handlers

import (
	"encoding/json"
	"fmt"
	"go-gorm-pg/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
)

var (
	Response models.Response
	user     *models.User
	// err      error
)

func Home(w http.ResponseWriter, r *http.Request) {

}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		db := models.DBCon()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		usr, _ := GetUser(user, db)
		if usr != nil {
			Response.Status = "Failed"
			Response.Message = "User already registered, please login"
			Response.Data = nil
			JSON(w, http.StatusBadRequest, Response)
			return
		}

		// user.Firstname = r.FormValue("firstname")
		user.Firstname = strings.TrimSpace(user.Firstname)
		// user.Lastname = r.FormValue("lastname")
		user.Lastname = strings.TrimSpace(user.Lastname)
		// user.Email = r.FormValue("email")
		user.Email = strings.TrimSpace(user.Email)
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			log.Fatal("Invalid Email")
		}
		// password := r.FormValue("password")
		// password = strings.TrimSpace(password)
		password := strings.TrimSpace(user.Password)
		hashedpassword, err := HashPassword(password)
		if err != nil {
			log.Fatal("Error hashing password:", err)
			// tpl.ExecuteTemplate(w, "Register", err)
		}
		user.Password = string(hashedpassword)

		fmt.Printf("%s, %s, %s\n", user.Firstname, user.Lastname, user.Email)

		err = db.Debug().Create(&user).Error
		if err != nil {
			log.Fatal("Error when inserting:", err)
		}

		log.Println("=> Inserted: First Name: " + user.Firstname + " | Last Name: " + user.Lastname)

		w.WriteHeader(http.StatusCreated)
		Response.Data["user"] = user
		Response.Message = "User created successfully"
		Response.Status = "Success"
		JSON(w, http.StatusOK, Response)

		// http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	} else if r.Method == "GET" {
		Response.Status = "Failed"
		Response.Message = "Method not allowed"
		Response.Data = nil
		JSON(w, http.StatusBadRequest, Response)
		// tpl.ExecuteTemplate(w, "register.html", nil)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		db := models.DBCon()

		user := &models.User{}
		body, err := ioutil.ReadAll(r.Body) // read user input from request
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		// Prepare(user) // strip the text of white spaces

		err = Validate(user, "login") // fields(email, password) are validated
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		usr, err := GetUser(user, db)
		if err != nil {
			ERROR(w, http.StatusInternalServerError, err)
			return
		}

		// "Email": "bmadukoma@gmail.com",
		// "Password": "blessed10"
		if usr == nil { // user is not registered
			Response.Status = "Failed"
			Response.Message = "Login failed because you have not registered, please sign up"
			Response.Data = nil
			JSON(w, http.StatusBadRequest, Response)
			return
		}

		err = CheckPasswordHash(user.Password, usr.Password)
		if err != nil {
			Response.Status = "Failed"
			Response.Message = "Wrong Password! Please try again!"
			Response.Data = nil
			JSON(w, http.StatusForbidden, Response)
			return
		}
		token, err := EncodeAuthToken(usr.ID)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		m := make(map[string]interface{})
		m["token"] = token
		m["user_id"] = usr.ID
		m["user"] = usr
		Response.Data = m // putting the token inside the data
		Response.Message = "Successfully logged in!"
		Response.Status = "Success"
		JSON(w, http.StatusOK, Response)
		return

	} else {
		Response.Status = "Failed"
		Response.Message = "Method not allowed"
		Response.Data = nil
		JSON(w, http.StatusMethodNotAllowed, Response)
	}
}
