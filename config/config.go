package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

var DB *sql.DB // Global variable to hold the *sql.DB instance
var STATIC_FILE_PATH = "C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\YoutubeConnector\\static\\videos"
var VIDEO_PATH2 = "C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\YoutubeConnector\\testvd.mp4"
var AUTH_MODE = "API_TOKEN"
var ROOT_URL string = "https://www.googleapis.com"
var OBJECT_URL = "https://storage.cloud.google.com/ytc-media-storage/pets_dog.mp4"
var UPLOAD_END_POINT = "upload/youtube/v3/videos"

var UTUBE_END_ENDPOINT string = "/youtube/v3/SampleVideo"

// var UPLOAD_ENDPOINT string = "?part=snippet%2Cstatus"
var API_Key string = "AIzaSyC_5hvxTsU8vijTreOE5zrwAws9XnCH6is"
var OAUTH_TOKEN string = "ya29.a0AfB_byDtDOdgo7Dzc_Pg-URblTXa3VSB5j-KAzkFyvmcO3u9PcCA-t3-8mN21irBWaKnUy9Lh444iZorST4x_uj-f0UnrY8lvUunI7vT8bDhnp3mLw0JhQWzJZ0H2FzVr1e8pnwgUcQRnC9GvmRtmgGWIGngg-QnwsAaCgYKAQYSARESFQHGX2MiX1UfikBfsVQooD-4dNVW5g0170"
var OAUTH_TOKEN_KR8799 string = "ya29.a0AfB_byCv-BmAK_FyYoVCKA9HwxsY2plLx6YcbmjNnkmQU2PWlVS6QzBVJYR7yno_hVdlGZBIRF1estnb6Z74NXvleqw--m1X6PBoHSZH4Au87jtJNp-pUyVCQ5Vu7rKKqVgPD5ddhIkbxFYUbUXMyahZq-ubOFWArspYaCgYKARISARASFQHGX2MiinTrcahbuysTGBTPd0d7vg0171"
var OAUTH_TOKEN_KR4DRI string = "ya29.a0AfB_byAMUjYzuFfeSjAOttAq9bsmkGldvKgjHVv5eJIQtkbfeCGyHY707MQAifHoPTpdB9HBuA_rKcEUbTaScEPMHPUz2sV2hndGAMM_nJ0qV_b2_k6_2i0B-Z-O8UgO_KmLqU_a_52bXFYcYaRbgGHvaRcw8Ehk77xnaCgYKAcMSARESFQHGX2Mik_rcjhZTiilJZ94HNQaFfA0171"
var VIDEO_PATH string = "C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\SampleVideo.mp4"
var POST_METHOD = "POST"
var GET_METHOD = "GET"
var CLIENT_ID = "397765413570-3thqle7blon88v54bgfuo3l7k5esukrh.apps.googleusercontent.com"
var CLIENT_SECRET = "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"
var REDIRECT_URL = "http://localhost:8080/callback"
var ServiceAccountPath = "C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\YoutubeConnector\\ytconnectormedia-410409-fbba40340683.json"
var BUCKET_NAME = "ytc-media-storage"
var UPLOAD_API_FILESYTEM = "https://storage.googleapis.com/upload/storage/v1/b/" + BUCKET_NAME + "/o?uploadType=media"

func InitConfigurations() {
	os.Setenv("API_Key", API_Key)
	fmt.Println("API_Key", os.Getenv("API_Key"))

	os.Setenv("CLIENT_ID", CLIENT_ID)
	fmt.Println("CLIENT_ID", os.Getenv("CLIENT_ID"))

	os.Setenv("CLIENT_SECRET", CLIENT_SECRET)
	fmt.Println("CLIENT_SECRET", os.Getenv("CLIENT_SECRET"))
}

/*
snippet.title
snippet.description
snippet.tags[]
snippet.categoryId

	type formValue struct {
			Title           string   `json:"title"`
			Description     string   `json:"description"`
			Category        string   `json:"category"`
			Audience        string   `json:"audience"`
			AgeRestrictions string   `json:"ageRestrictions"`
			TagInput        []string `json:"tagInput"`
			Privacy         string   `json:"privacy"`
		}
*/

type WebConfig struct {
	ClientID                string   `json:"client_id"`
	RedirectURIs            []string `json:"redirect_uris"`
	ProjectID               string   `json:"project_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX500CertUrl string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
}

type Configuration struct {
	Web WebConfig `json:"web"`
}

type ServiceAccount struct {
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

func readFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Error while loading the client secret")
		return nil, nil
	}

	return content, nil
}
func LoadClientAndSecretKey(path string) (string, string) {
	content, err := readFile(path)
	if err != nil {
		return "INVALID", "INVALID"
	}
	var config Configuration
	err = json.Unmarshal(content, &config)
	if err != nil {
		return "INVALID", "INVALID"
	}

	return config.Web.ClientID, config.Web.ClientSecret
}

func LoadServiceAccount(path string) (*ServiceAccount, error) {
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}

	var tokenConfig ServiceAccount
	err = json.Unmarshal(content, &tokenConfig)

	if err != nil {
		return nil, err
	}
	return &tokenConfig, nil
}
var USERID string;
var ObjectTblCreation = `IF not EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'objectTbl' )
BEGIN
  CREATE TABLE objectTbl (
		userName varchar(50),
		objectId VARCHAR(50) NOT NULL,
		bucketId VARCHAR(50) NOT NULL
	)
END`

var UserTableCreation = `IF not EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'userTbl')
BEGIN
	CREATE table userTbl(
	  userId varchar(50) primary key,
	  password varchar(50) not null
	)
END`


func DbInit() {
	server := "MF-H59IBOW2THNM"
	port := 1433
	user := "sa"
	password := "1iso*help"
	database := "ytc"

	// Construct the DSN
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)
	fmt.Println(dsn)
	// Open a connection to the SQL Server database
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	// Ping the database to check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the SQL Server database!")
	DB = db

	// Create the table
	createTable := ObjectTblCreation
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(UserTableCreation, "userTbl")
	if err != nil {
		log.Fatal(err)
	}
}


func GenerateJWTToken() (string, error) {
	// CREATING SOME WRONG JWT KEY HENCE USED THE CLIENT LIBRARY
	// serviceAccount, err := config.LoadServiceAccount(config.ServiceAccountPath)

	// if err != nil {
	// 	log.Fatal("Error loading service account:", err)
	// 	return "", err
	// }

	// // Create a new token object with claim
	// token := jwt.New(jwt.SigningMethodRS256)

	// //set the token with claims
	// claims := token.Claims.(jwt.MapClaims)
	// claims["sub"] = serviceAccount.ClientEmail
	// //claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Replace with your desired expiration time
	// fmt.Println(claims["exp"])
	// claims["iss"] = serviceAccount.ClientEmail
	// claims["scope"] = storage.DevstorageFullControlScope
	// claims["aud"] = "https://www.googleapis.com/oauth2/v4/token"

	// // parse the private key
	// signingKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(serviceAccount.PrivateKey))
	// if err != nil {
	// 	log.Fatal("Error parsing the private key ", err)
	// 	return "", err
	// }

	// jwtSignedToken, err := token.SignedString(signingKey)
	// if err != nil {
	// 	log.Fatal("Error in signing the token ", err)
	// 	return "", err
	// }
	// return jwtSignedToken, nil

	// Load the service account JSON key file
	serviceAccountData, err := os.ReadFile(ServiceAccountPath)
	// serviceAccountData, err := ioutil.ReadFile(config.ServiceAccountPath)
	if err != nil {
		log.Fatalf("Error reading service account JSON: %v", err)
	}

	// Create a JWT Config from the service account JSON
	configToken, err := google.JWTConfigFromJSON(serviceAccountData, storage.DevstorageFullControlScope)
	if err != nil {
		log.Fatalf("Error creating JWT Config: %v", err)
	}

	// Set headers, including the Authorization header with the JWT token
	key, val := configToken.TokenSource(context.Background()).Token()
	if val != nil {
		fmt.Println(err)
	}
	//fmt.Println(key.AccessToken)
	return key.AccessToken, nil

}

