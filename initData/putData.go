package initdata

import (
	"GoFinalProject/structure"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const (
	dbPathUsers = "./dataBase/users.json"
	dbPathGames = "./dataBase/games.json"
)

func PutInitData() {
	file, err := os.Create(dbPathUsers)
	if err != nil {
		log.Fatalf("Cannot create file,\nError:%v", err)
	}
	defer file.Close()

	password := "HeyYo2001"
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	password2 := "GoodShit07"
	hash2, _ := bcrypt.GenerateFromPassword([]byte(password2), bcrypt.DefaultCost)

	Users := []structure.User{
		{
			Email:    "Bruh@mail.ru",
			Username: "BruhBruhovski",
			Password: string(hash),
		},
		{
			Email:    "IITUsh@gmail.com",
			Username: "Kendrihno",
			Password: string(hash2),
		},
	}

	byteSlice, _ := json.MarshalIndent(Users, "", " ")
	err = ioutil.WriteFile(dbPathUsers, byteSlice, 0644)
	if err != nil {
		log.Fatalf("Cannot write to file,\nError:%v", err)
	}

	fileForGames, err := os.Create(dbPathGames)
	if err != nil {
		log.Fatal(err)
	}
	defer fileForGames.Close()

	games := []structure.GamePage{
		{
			Id:        "1",
			Title:     "Prototype",
			Desc:      "You are the Prototype, Alex Mercer, a man without memory armed with amazing shape-shifting abilities, hunting your way to the heart of the conspiracy which created you; making those responsible pay.",
			Rating:    "Not rated",
			Developer: "Ctytek Software",
			Genre:     "Slasher, open world",
			VideoUrl:  "https://www.youtube.com/embed/Nc3XptLacMM?controls=0",
			PhotoUrl:  "https://upload.wikimedia.org/wikipedia/ru/c/c0/Prototype_Cover_Art.jpeg",
			Comments:  []structure.Comment{},
		},
		{
			Id:        "2",
			Title:     "Borderlands 2",
			Desc:      "A new era of shoot and loot is about to begin. Play as one of four new vault hunters facing off against a massive new world of creatures, psychos and the evil mastermind, Handsome Jack. Make new friends, arm them with a bazillion weapons and fight alongside them in 4 player co-op on a relentless quest for revenge and redemption across the undiscovered and unpredictable living planet.",
			Rating:    "Not rated",
			Developer: "Gearbox software",
			Genre:     "FPS, RPG, open world",
			VideoUrl:  "https://www.youtube.com/watch?v=5TW0wJTFLiw",
			PhotoUrl:  "https://media.redadn.es/imagenes/borderlands-2_152890.jpg",
			Comments:  []structure.Comment{},
		},
	}

	byteSliceForGames, _ := json.MarshalIndent(games, "", " ")
	err = ioutil.WriteFile(dbPathGames, byteSliceForGames, 0644)
	if err != nil {
		log.Fatalf("Cannot write to file,\nError:%v", err)
	}
}
