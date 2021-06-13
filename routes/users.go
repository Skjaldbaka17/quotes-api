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

var requestsPerHour = map[string]int{"free": 100, "basic": 1000, "lilleBoy": 100000, "GOD": -1}
var TIERS = []string{"free", "basic", "lilleBoy", "GOD"}

// swagger:route POST /users USERS SignUp
// Create A user to get a free ApiKey
// responses:
//	200: userResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// CreateUsers handles post requests to create a user and an accompanying ApiKey
func CreateUser(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.UserRequest
	if err := handlers.GetUserRequestBody(rw, r, &requestBody); err != nil {
		return
	}

	if err := handlers.ValidateUserInformation(rw, r, &requestBody); err != nil {
		return
	}

	apiKey, _ := uuid.NewRandom()
	passHash, _ := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	apiKeyHash := hash(apiKey.String())
	requestBody.Tier = "lilleBoy"

	user := structs.User{Name: requestBody.Name, ApiKeyHash: apiKeyHash, Tier: requestBody.Tier, Email: requestBody.Email, PasswordHash: string(passHash)}

	result := handlers.Db.Table("users").Select("name", "api_key_hash", "tier", "email", "password_hash").Create(&user)

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

	json.NewEncoder(rw).Encode(structs.User{Id: user.Id, ApiKey: apiKey.String()})
}

func Login(rw http.ResponseWriter, r *http.Request) {
	var requestBody structs.UserRequest
	if err := handlers.GetUserRequestBody(rw, r, &requestBody); err != nil {
		return
	}
}

func UpgradeTier(rw http.ResponseWriter, r *http.Request) {}

func DowngradeTier(rw http.ResponseWriter, r *http.Request) {}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {}

func hash(stringToHash string) string {
	bv := []byte(stringToHash)
	hasher := sha256.New()
	hasher.Write(bv)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
