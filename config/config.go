package config

import (
	"encoding/json"
	"fmt"

	"os"
)

var VIDEO_PATH2 = "C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\YoutubeConnector\\testvd.mp4"
var AUTH_MODE = "API_TOKEN"
var ROOT_URL string = "https://www.googleapis.com"
var OBJECT_URL = "https://storage.cloud.google.com/ytc-media-storage/pets_dog.mp4"
var UPLOAD_END_POINT = "upload/youtube/v3/videos"

var UTUBE_END_ENDPOINT string = "/youtube/v3/SampleVideo"

// var UPLOAD_ENDPOINT string = "?part=snippet%2Cstatus"
var API_Key string = "AIzaSyC_5hvxTsU8vijTreOE5zrwAws9XnCH6is"
var OAUTH_TOKEN string = "ya29.a0AfB_byDtDOdgo7Dzc_Pg-URblTXa3VSB5j-KAzkFyvmcO3u9PcCA-t3-8mN21irBWaKnUy9Lh444iZorST4x_uj-f0UnrY8lvUunI7vT8bDhnp3mLw0JhQWzJZ0H2FzVr1e8pnwgUcQRnC9GvmRtmgGWIGngg-QnwsAaCgYKAQYSARESFQHGX2MiX1UfikBfsVQooD-4dNVW5g0170"
var OAUTH_TOKEN_KR8799 string = "ya29.a0AfB_byC1uPyVDG3H6feJhhyfw3vdfsMbNw16X7ve2Z-EIyfmUtJPSLsL6F_Lgo_3yLCQInzKDyhhHmWhPC2cj226xRHjczTq58MjAfE9T6e_oY1s2BtVsp3gQXE5pw1hzZmhznw_A6MU9PJLUDKReuKPcUX0ZtyhGgaCgYKATkSARESFQHGX2MiPeoz-mT7n8b9pZa6jTKfew0169"
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
