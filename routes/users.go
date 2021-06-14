package routes

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/Skjaldbaka17/quotes-api/handlers"
	"github.com/Skjaldbaka17/quotes-api/structs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// swagger:route POST /users/signup USERS SignUp
// Create A user to get a free ApiKey
// responses:
//	200: userResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// CreateUsers handles post requests to create a user and an accompanying ApiKey
func CreateUser(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.UserApiModel
	if err := handlers.GetUserRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	if err := handlers.ValidateUserInformation(rw, r, &requestBody); err != nil {
		return
	}

	uuid, _ := uuid.NewRandom()
	apiKey := uuid.String()
	passHash, _ := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	requestBody.Tier = handlers.TIERS[0]
	log.Println("The created apikey:", apiKey)
	user := structs.UserDBModel{Name: requestBody.Name, ApiKey: apiKey, Tier: requestBody.Tier, Email: requestBody.Email, PasswordHash: string(passHash)}

	result := handlers.Db.Table("users").Select("name", "api_key", "tier", "email", "password_hash").Create(&user)

	//Error handle
	if result.Error != nil {
		m1 := regexp.MustCompile(`duplicate key value violates unique constraint "users_email_key"`)
		if m1.Match([]byte(result.Error.Error())) {
			rw.WriteHeader(http.StatusBadRequest)
			log.Printf("Got error when creating user, constraint error: %s", result.Error)
			json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "This email is taken."})
			return
		}
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got error when creating user: %s", result.Error)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	} else if user.Id <= 0 {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Got no id when creating user: %s", result.Error)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
		return
	}

	json.NewEncoder(rw).Encode(structs.UserResponse{Id: user.Id, ApiKey: user.ApiKey})
}

// swagger:route POST /users/login USERS Login
// Login to get the api_key for the user
// responses:
//	200: userResponse
//  400: incorrectBodyStructureResponse
//  401: incorrectCredentialsResponse
//  500: internalServerErrorResponse

// Login handles post requests to login to a user and receive his ApiKey
func Login(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.UserApiModel
	if err := handlers.GetUserRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	var user structs.UserDBModel
	if err := handlers.Db.Table("users").Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Printf("Got error when login/fetching user: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "No user with the given email address. Maybe try WindsOfWinterWillNeverBeFinished@WeAreAllSinnersBeforeTheSeven.com"})
		return
	}

	//Compare passwords / Check correct password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(requestBody.Password)); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		log.Printf("Got error when comparing passwords in login: %s", err)
		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "Credentials not correct. Shame. Shame. Shame is the name of the game."})
		return
	}

	json.NewEncoder(rw).Encode(structs.UserResponse{ApiKey: user.ApiKey})
}

func UpgradeTier(rw http.ResponseWriter, r *http.Request) {}

func DowngradeTier(rw http.ResponseWriter, r *http.Request) {}

func UpdateUser(rw http.ResponseWriter, r *http.Request) {}

func ReplaceApiKey(rw http.ResponseWriter, r *http.Request) {}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {}

func hash(stringToHash string) string {
	bv := []byte(stringToHash)
	hasher := sha256.New()
	hasher.Write(bv)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
