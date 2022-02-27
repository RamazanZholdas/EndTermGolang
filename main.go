package main

import (
	initdata "GoFinalProject/initData"
	"GoFinalProject/loggs"
	"GoFinalProject/structure"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbPathUsers = "./dataBase/users.json"
	dbPathGames = "./dataBase/games.json"
	port        = ":8080"
)

type UniversalStruct struct {
	Auth structure.Auth
	Game []structure.GamePage
}

var (
	isAuthenticated = structure.Auth{}
	games           = []structure.GamePage{}
	doubleSexual    = UniversalStruct{}
	emailRegex      = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func init() {
	initdata.PutInitData()

	file, err := os.Open(dbPathGames)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteSlice, &games)

	isAuthenticated = structure.Auth{
		IsAuth:   false,
		Username: "",
	}

	doubleSexual = UniversalStruct{
		Auth: isAuthenticated,
		Game: games,
	}
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func main() {
	file := loggs.SetupLogOutput()
	defer file.Close()

	router := gin.New()
	router.Use(gin.Recovery(), loggs.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/mainPage", mainPage)
	router.GET("/signUpPage", signUpPage)
	router.POST("/signUpPageAuth", signUpPageAuth)
	router.GET("/signInPage", signInPage)
	router.POST("/signInPageAuth", signInPageAuth)
	router.GET("/gamePage/:id", gamePage)
	router.POST("/gamePage/:id/postComment", gamePageWithComment)
	router.GET("/logoutPage", logoutPage)
	router.GET("/adminPanel", adminPanel)

	router.Run(port)
}

func adminPanel(c *gin.Context) {
	file, err := os.Open(dbPathUsers)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	users := []structure.User{}
	byteSlice, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteSlice, &users)

	allDb := struct {
		Users []structure.User
		Games []structure.GamePage
	}{
		Users: users,
		Games: doubleSexual.Game,
	}

	c.HTML(http.StatusOK, "adminPanel.html", allDb)
}

func mainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", doubleSexual)
}

func signUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signUp.html", nil)
}

func logoutPage(c *gin.Context) {
	doubleSexual.Auth.IsAuth = false
	doubleSexual.Auth.Username = ""
	c.HTML(http.StatusOK, "index.html", doubleSexual)
}

func signUpPageAuth(c *gin.Context) {
	c.Request.ParseForm()

	username := c.Request.FormValue("username")

	firstCheck := true
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			firstCheck = false
		}
	}

	secondCheck := false
	if len(username) >= 5 && len(username) <= 35 {
		secondCheck = true
	}

	fmt.Println("first check:", firstCheck, "second check:", secondCheck)
	if !firstCheck || !secondCheck {
		c.HTML(http.StatusSeeOther, "signUp.html", "Your username must contain only letters and digits no bullshit like spaces and symbols")
		return
	}

	email := c.Request.FormValue("email")

	if !isEmailValid(email) {
		c.HTML(http.StatusSeeOther, "signUp.html", "Your email is not valid")
		return
	}

	password := c.Request.FormValue("password")

	var haveLowercase, haveUppercase, haveNumber, enoughLength, haveNoSpaces bool
	haveNoSpaces = true
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			haveLowercase = true
		case unicode.IsUpper(char):
			haveUppercase = true
		case unicode.IsNumber(char):
			haveNumber = true
		case len(password) > 5 && len(password) < 30:
			enoughLength = true
		case unicode.IsSpace(int32(char)):
			haveNoSpaces = false
		}
	}

	fmt.Println("havelowercase:", haveLowercase, "haveUppercase", haveUppercase, "haveNumber", haveNumber, "enoughLenght", enoughLength, "haveNoSpaces", haveNoSpaces)
	if !haveLowercase || !haveUppercase || !haveNumber || !enoughLength || !haveNoSpaces {
		c.HTML(http.StatusSeeOther, "signUp.html", "U idiot, pls fix your password")
		return
	}

	passwordAgain := c.Request.FormValue("passwordAgain")
	if passwordAgain != password {
		c.HTML(http.StatusSeeOther, "signUp.html", "your password doesnt match with password that you wrote again")
		return
	}

	file, err := os.Open(dbPathUsers)
	if err != nil {
		log.Fatalf("Can't open a file,\nError:%v", err)
	}
	defer file.Close()

	users := []structure.User{}
	byteSlice, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Can't read file,\nError:%v", err)
	}
	err = json.Unmarshal(byteSlice, &users)
	if err != nil {
		log.Fatal(err)
	}

	for index := range users {
		if users[index].Username == username {
			c.HTML(http.StatusSeeOther, "signUp.html", "Username already taken")
			return
		}
		if users[index].Email == email {
			c.HTML(http.StatusSeeOther, "signUp.html", "User that uses this email exists")
			return
		}
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	users = append(users, structure.User{
		Email:    email,
		Username: username,
		Password: string(hash),
	})

	byteSlice, err = json.MarshalIndent(users, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dbPathUsers, byteSlice, 0644)
	if err != nil {
		log.Fatal(err)
	}

	doubleSexual.Auth.IsAuth = true
	doubleSexual.Auth.Username = username
	c.HTML(http.StatusOK, "index.html", doubleSexual)
}

func signInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signIn.html", nil)
}

func signInPageAuth(c *gin.Context) {
	c.Request.ParseForm()

	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	file, err := os.Open(dbPathUsers)
	if err != nil {
		log.Fatalf("Can't open a file,\nError:%v", err)
	}
	defer file.Close()

	users := []structure.User{}
	byteSlice, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Can't read file,\nError:%v", err)
	}
	err = json.Unmarshal(byteSlice, &users)
	if err != nil {
		log.Fatal(err)
	}

	var firstCheck, secondCheck bool
	var i int

	for index := range users {
		if users[index].Email == email {
			firstCheck = true
			if err = bcrypt.CompareHashAndPassword([]byte(users[index].Password), []byte(password)); err == nil {
				secondCheck = true
				i = index
			}
		}
		if firstCheck && secondCheck {
			break
		}
	}

	if firstCheck && secondCheck {
		/*isAuthenticated.IsAuth = true
		isAuthenticated.Username = users[i].Username*/
		doubleSexual.Auth.IsAuth = true
		doubleSexual.Auth.Username = users[i].Username
		c.HTML(http.StatusOK, "index.html", doubleSexual)
	} else {
		c.HTML(http.StatusSeeOther, "signIn.html", "Wrong password or email")
		return
	}
}

func gamePage(c *gin.Context) {
	id := c.Param("id")
	games := []structure.GamePage{}

	file, err := os.Open(dbPathGames)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteSlice, &games)

	for index := range games {
		if games[index].Id == id {
			authWithGame := struct {
				Game    structure.GamePage
				Authent structure.Auth
			}{
				Game: games[index],
				Authent: structure.Auth{
					IsAuth:   doubleSexual.Auth.IsAuth,
					Username: doubleSexual.Auth.Username,
				},
			}

			c.HTML(http.StatusOK, "gamePage.html", authWithGame)
			return
		}
	}

	c.String(http.StatusNotFound, "This page does not exits %s", c.Request.URL.Path)
}

func gamePageWithComment(c *gin.Context) {
	c.Request.ParseForm()
	id := c.Param("id")
	games := []structure.GamePage{}

	file, err := os.Open(dbPathGames)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteSlice, &games)

	for index := range games {
		if games[index].Id == id {
			games[index].Comments = append(games[index].Comments, structure.Comment{
				Username:   doubleSexual.Auth.Username,
				UserRating: c.Request.FormValue("rating"),
				Review:     c.Request.FormValue("review"),
			})

			doubleSexual.Game[index] = games[index]

			fmt.Println("Rating:", c.Request.FormValue("rating"))
			avgRate := fmt.Sprint(averageRating(index))
			games[index].Rating = avgRate

			doubleSexual.Game[index] = games[index]

			bytesliceWithGame, err := json.MarshalIndent(games, "", " ")
			if err != nil {
				log.Fatal(err)
			}

			ioutil.WriteFile(dbPathGames, bytesliceWithGame, 0644)

			authWithGame := struct {
				Game    structure.GamePage
				Authent structure.Auth
			}{
				Game: games[index],
				Authent: structure.Auth{
					IsAuth:   doubleSexual.Auth.IsAuth,
					Username: doubleSexual.Auth.Username,
				},
			}

			c.HTML(http.StatusOK, "gamePage.html", authWithGame)
			return
		}
	}
}

func averageRating(index int) float64 {
	counter := 0.0
	sum := 0.0
	var temp int
	for i := range doubleSexual.Game[index].Comments {
		fmt.Println(doubleSexual.Game[index].Comments[i].UserRating)
		temp, _ = strconv.Atoi(doubleSexual.Game[index].Comments[i].UserRating)
		sum += float64(temp)
		counter++
	}
	fmt.Println(sum, counter)
	if counter == 0 {
		return 0.0
	}
	return sum / counter
}
