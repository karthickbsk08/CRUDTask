package jwttokengen

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"tasks/tomlutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

//Generate JWT Token
//Step 1 : Generate JWT Secret Key (Random string only knows by server)

func Get_JWT_Secret_Key() string {
	return fmt.Sprintf(`%v`, tomlutil.ReadTomlConfig("config.toml").(map[string]any)["JWT_Secret_Key"])
}

func Generate_JWT_Secret_Key() error {

	// Generate 64 random bytes (512-bit key)
	key := make([]byte, 64)

	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	// Encode to Base64 (same format as openssl output)
	base64Key := base64.StdEncoding.EncodeToString(key)

	// Print the key
	fmt.Println("Your secure JWT secret key:")
	fmt.Println(base64Key)

	err = tomlutil.WriteTomlFile("JWT_Secret_Key", base64Key, "toml/config.toml")
	if err != nil {
		return err
	}
	return nil
}

/*
========  Registered Claims (standard, optional fields)============

These are common fields defined in the JWT standard:

| Field | Type        | Purpose                                    |
| ----- | ----------- | ------------------------------------------ |
| `exp` | `ExpiresAt` | Token expiry time                          |
| `iat` | `IssuedAt`  | When the token was issued                  |
| `nbf` | `NotBefore` | Token is not valid before this time        |
| `iss` | `Issuer`    | Who issued the token (e.g., your app name) |
| `sub` | `Subject`   | What the token is about (user ID, etc.)    |
| `aud` | `Audience`  | Who the token is intended for              |
| `jti` | `JWT ID`    | Unique identifier for the token            |
*/

func GenerateRegisterClaims(pPurpose string) jwt.RegisteredClaims {

	lConfigInterface := tomlutil.ReadTomlConfig("config.toml")
	lExpireTime, _ := strconv.Atoi(fmt.Sprintf("%v", lConfigInterface.(map[string]any)["EXPIRE_TIME_JWT_TOKEN"]))
	pAppName := fmt.Sprintf("%v", lConfigInterface.(map[string]any)["JWT_App_Name"])

	session := uuid.NewV4()
	sessionSHA256 := session.String()
	KeyValue := strings.ReplaceAll(sessionSHA256, "-", "")

	var RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(lExpireTime) * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    pAppName,
		Subject:   pPurpose,
		ID:        KeyValue,
	}

	return RegisteredClaims

}

/*
Custom Claims (your own fields)
These are your app’s data — you choose what to put in:

You might include:

user_id – primary key from your DB
email – to identify the user
role – like admin, user, guest
permissions – optional list
name, client_id, etc.
*/
func GenerateUserDetails(pUsername, pPwd string) UserDetails {

	var lUserDetails = UserDetails{
		UserID:   pUsername,
		Password: pPwd,
	}
	return lUserDetails
}

func CreateToken(pclaims Claims) *jwt.Token {
	// Create token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, pclaims)
}

func SignGenToken(pToken *jwt.Token, jwtKey string) (string, error) {
	var tokenString string
	var lErr error

	// Sign token
	tokenString, lErr = pToken.SignedString(jwtKey)
	if lErr != nil {
		log.Println("JWT token sign error")
		return tokenString, lErr
	}

	return tokenString, nil
}

func GenerateJWT(u, p string) (string, error) {

	var lJWT_Token string
	var lErr error

	//Step 1 : Generate JWT Secret Key
	lJWT_Secret_Key := Get_JWT_Secret_Key()

	//Step 2 : Generate Claims Struct
	var lCliams Claims

	//Step 2.1 Custom Claims
	lCliams.UserDetails = GenerateUserDetails(u, p)

	//Step 2.2 Register Cliams
	lCliams.RegisteredClaims = GenerateRegisterClaims(USER_TOKEN)

	//Step 3 : Generate Token (Unsigned token)
	lJWTTokenWithoutSigned := CreateToken(lCliams)

	//Step 4 : Signing Unsigned Token
	lJWT_Token, lErr = SignGenToken(lJWTTokenWithoutSigned, lJWT_Secret_Key)
	if lErr != nil {
		return lJWT_Token, lErr
	}

	return lJWT_Token, nil
}
