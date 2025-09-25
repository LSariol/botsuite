package main

import (
	"fmt"
	"log"

	"github.com/lsariol/botsuite/internal/adapters/twitch/auth"
	"github.com/lsariol/botsuite/internal/app"
)

// Use to Authenticate a user once
func main() {

	setUpNewUser()

}

func setUpNewUser() {
	userToken := "ztmn1rfniu340rfnxv5xy1fdsgiblo"
	var userData auth.UserData
	var success *auth.VerificationSuccess

	deps, err := app.NewDependencies()
	if err != nil {
		log.Fatal("Error loading deps")
	}

	err = auth.RefreshAppAccessToken(&deps.Config.Twitch, deps.HTTP)
	if err != nil {
		log.Fatal(err)
	}

	userData, err = auth.GenerateUserAcessToken(userToken, &deps.Config.Twitch, deps.HTTP)
	if err != nil {
		log.Fatal(err)
	}

	userData, err = auth.RefreshUserAccessToken(userData.UserRefreshToken, &deps.Config.Twitch, deps.HTTP)
	if err != nil {
		log.Fatal(err)
	}

	success, err = auth.ValidateToken(userData.UserAccessToken, deps.HTTP)
	if err != nil {
		log.Fatal(err)
	}

	userData.UserID = success.UserId
	userData.Username = success.Login

	err = userData.Store()
	if err != nil {
		fmt.Println(userData)
		log.Fatal(err)
	}

	fmt.Printf("\n%s has been added to the database.\n", userData.Username)
}
