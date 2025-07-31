package jwttokengen

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"tasks/common"
	"tasks/helpers"
	"tasks/tomlutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

//Generate JWT Token
//Step 1 : Generate JWT Secret Key (Random string only knows by server)

func Get_JWT_Secret_Key(pDebug *helpers.HelperStruct) string {
	pDebug.Log(helpers.Statement, "Get_JWT_Secret_Key(+)")
	return fmt.Sprintf(`%v`, tomlutil.ReadTomlConfig("./toml/config.toml").(map[string]any)["JWT_Secret_Key"])
}

func Generate_JWT_Secret_Key(pDebug *helpers.HelperStruct) error {
	pDebug.Log(helpers.Statement, "Generate_JWT_Secret_Key(+)")
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
	pDebug.Log(helpers.Statement, "Generate_JWT_Secret_Key(-)")
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

func GenerateRegisterClaims(pDebug *helpers.HelperStruct, pPurpose string) jwt.RegisteredClaims {
	pDebug.Log(helpers.Statement, "GenerateRegisterClaims(+)")

	lConfigInterface := tomlutil.ReadTomlConfig("toml/config.toml")
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
	pDebug.Log(helpers.Statement, "GenerateRegisterClaims(-)")

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
func GenerateUserDetails(pDebug *helpers.HelperStruct, pUsername, pPwd string) UserDetails {

	var lUserDetails = UserDetails{
		UserID:   pUsername,
		Password: pPwd,
	}
	return lUserDetails
}

func CreateToken(pDebug *helpers.HelperStruct, pclaims Claims) *jwt.Token {
	pDebug.Log(helpers.Statement, "CreateToken(+)")
	// Create token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, pclaims)
}

func SignGenToken(pDebug *helpers.HelperStruct, pToken *jwt.Token, jwtKey []byte) (string, error) {
	pDebug.Log(helpers.Statement, "SignGenToken(+) , jwtKey : ", jwtKey)

	var tokenString string
	var lErr error

	h, err := json.Marshal(pToken.Header)
	if err != nil {
		return "", err
	}

	c, err := json.Marshal(pToken.Claims)
	if err != nil {
		return "", err
	}

	log.Println("h : ", string(h))
	log.Println("c : ", string(c))
	log.Println("token : ", pToken.EncodeSegment(h)+"."+pToken.EncodeSegment(c))
	log.Println("jwttoken : ", common.DoMarshall(pToken))

	// Sign token
	tokenString, lErr = pToken.SignedString(jwtKey)
	if lErr != nil {
		log.Println("JWT token sign error")
		return tokenString, lErr
	}
	pDebug.Log(helpers.Statement, "SignGenToken(-)")
	return tokenString, nil
}

func GenerateJWT(pDebug *helpers.HelperStruct, u, p string) (string, error) {
	pDebug.Log(helpers.Statement, "GenerateJWT(+)")

	var lJWT_Token string
	var lErr error
	var lDecodedsecretBytes []byte

	//Step 1 : Generate JWT Secret Key
	lJWT_Secret_Key := Get_JWT_Secret_Key(pDebug)

	//Step 2 : Generate Claims Struct
	var lCliams Claims

	//Step 2.1 Custom Claims
	lCliams.UserDetails = GenerateUserDetails(pDebug, u, p)

	//Step 2.2 Register Cliams
	lCliams.RegisteredClaims = GenerateRegisterClaims(pDebug, USER_TOKEN)

	//Step 3 : Generate Token (Unsigned token)
	lJWTTokenWithoutSigned := CreateToken(pDebug, lCliams)

	//Step 3.1 Decode JWT Secret Key
	lDecodedsecretBytes, lErr = base64.StdEncoding.DecodeString(lJWT_Secret_Key)
	if lErr != nil {
		return lJWT_Token, lErr
	}

	//Step 4 : Signing Unsigned Token
	lJWT_Token, lErr = SignGenToken(pDebug, lJWTTokenWithoutSigned, lDecodedsecretBytes)
	if lErr != nil {
		return lJWT_Token, lErr
	}

	log.Println("toekn : ", lJWT_Token)
	// lJWT_Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoia2FydGhpY2siLCJlbWFpbCI6IiIsInJvbGUiOiIiLCJwd2QiOiJiZXN0QDEyMyIsImlzcyI6InRhc2tzIiwic3ViIjoidXNlcl90b2tlbiIsImV4cCI6MTc1MTcyODA1MCwiaWF0IjoxNzUxNzI0NjkwLCJqdGkiOiJjMzhmNDU1MGI5N2U0OGQxOTUyOWZjMThmODA0ZTU5MCJ9"
	pDebug.Log(helpers.Statement, "GenerateJWT(-)")
	return lJWT_Token, nil
}
