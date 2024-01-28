package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/abhis2007/YOUTUECONNECTOR/config"
	"github.com/abhis2007/YOUTUECONNECTOR/routes"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

type SnippetPayload struct {
	CategoryId  string `json:"categoryId"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type StatusPayload struct {
	PrivacyStatus string `json:"privacyStatus"`
}

type BodyPayload struct {
	Snippet SnippetPayload `json:"snippet"`
	Status  StatusPayload  `json:"status"`
}

var templates *template.Template

// executes automstically
func init() {
	templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))
}

// uploadVideo is more understable code base
func main() {
	//http.Handle("/", http.StripPrefix("/videos/", http.FileServer(http.Dir("./videos"))))
	//log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./videos/"))))

	Snippetdata := SnippetPayload{
		CategoryId:  "Keyboard, Mouse, MousePad, Laptop, Table",
		Description: "This is sample video uploaded from google utube API, for content like - Keyboard, Mouse, MousePad, Laptop, Table",
		Title:       "First sample video",
	}
	StatusData := StatusPayload{
		PrivacyStatus: "Unlisted",
	}
	BodyPayloadData := BodyPayload{
		Snippet: Snippetdata,
		Status:  StatusData,
	}
	body, _ := json.Marshal(BodyPayloadData)
	// fmt.Println(string(body))
	if false {
		callUpload()
		getVideo()

		uploadThumbnail()
		config.InitConfigurations()
		createAuthToken()
		uploadVideo(string(body))
		uploadObjectIntoBucket()
	}
	if !false {
		routerInstance := mux.NewRouter()

		// Serve static files (images in this case)
		routerInstance.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
		routes.RouterConfiguration(routerInstance)
		http.Handle("/", routerInstance)
		// createAuthToken()

		fmt.Println("Server started at port : 8080")

		//test()
		log.Fatal(http.ListenAndServe(":8080", routerInstance))
	}
	// test()
	// uploadObjectIntoBucket2()
	//updateObject()

}

// Update the objec inside the bucket
func updateObject() {
	objectUrl := config.OBJECT_URL
	file, err := os.ReadFile(config.VIDEO_PATH2)
	if err != nil {
		log.Fatalf("Error in parsing the file: %v", err)
		return
	}
	request, err := http.NewRequest("PUT", objectUrl, bytes.NewBuffer(file))
	if err != nil {
		fmt.Printf("Error creating PUT request: %v\n", err)
		return
	}
	//do the authentication JWT
	token, _ := generateJWTToken()
	bearerToken := "Bearer " + token
	fmt.Println(bearerToken)

	request.Header.Set("Content-Type", "video/mp4")
	request.Header.Set("Authorization", bearerToken)

	//create a client
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error in sending the put request: %v\n", err)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Print the response status and body
	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Body:", string(responseBody))
}

// fINAL AND TESTED FUNCTION WORKING OK FOR UPLOAD THE OBJECT INTO THE CLOUD STORAGE BUCKET
func test() {
	// Load the service account JSON key file
	serviceAccountData, err := os.ReadFile(config.ServiceAccountPath)
	if err != nil {
		log.Fatalf("Error reading service account JSON: %v", err)
	}

	// Create a JWT Config from the service account JSON
	configToken, err := google.JWTConfigFromJSON(serviceAccountData, storage.DevstorageFullControlScope)
	if err != nil {
		log.Fatalf("Error creating JWT Config: %v", err)
	}

	// Create an HTTP client with OAuth2 authentication
	client := configToken.Client(context.Background())

	// Set headers, including the Authorization header with the JWT token
	key, val := configToken.TokenSource(context.Background()).Token()
	if val != nil {
		fmt.Println(err)
	}
	fmt.Println(key.AccessToken)

	// Add the video file
	filePath := config.VIDEO_PATH

	baseFilePart := filepath.Base(filePath)
	lists := strings.Split(baseFilePart, ".")
	if len(lists) <= 0 {
		log.Fatalf("File name doesnot seems to be of a media type\n")
		return
	}
	var fileName string = lists[0]

	file, err := os.ReadFile(config.VIDEO_PATH)

	if err != nil {
		log.Fatalf("Error reading object data: %v", err)
	}

	//Form teh endpoints
	url := fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/ytc-media-storage/o?uploadType=media&name=%s", fileName)
	fmt.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(file)))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+key.AccessToken)
	req.Header.Set("Content-Type", "video/mp4")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Check if the request was successful (status code 2xx)
	if resp.StatusCode/100 != 2 {
		log.Fatalf("Error: %s", responseBody)
	}

	fmt.Println("Upload successful. Response:", string(responseBody))
}

func uploadObjectIntoBucket() {
	//url := fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/ytc-media-storage/o?uploadType=media&name=%s", "pets%2Fdog.mp4")
	url := "https://storage.googleapis.com/upload/storage/v1/b/ytc-media-storage/o?uploadType=multipart"
	fmt.Println(url)

	// Add the video file
	// file, err := ioutil.ReadFile(config.VIDEO_PATH)
	// file, err := os.ReadFile(config.VIDEO_PATH)
	file, err := os.Open(config.VIDEO_PATH)
	if err != nil {
		log.Fatalf("Error reading object data: %v", err)
	}
	defer file.Close()

	// create a multipart
	buffer := bytes.Buffer{}
	writer := multipart.NewWriter(&buffer)

	//Add the Video in it
	part, err := writer.CreateFormFile("SampleVideoFile", "pets%2Fdog.mp4")
	if err != nil {
		log.Fatalf("Error creating the formFile %v", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("Error copying the formFile part %v", err)
		return
	}
	//close the writer
	writer.Close()

	//create the http request
	// request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(file)))
	request, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	//do the authentication JWT
	token, _ := generateJWTToken()
	bearerToken := "Bearer " + token
	fmt.Println(bearerToken)

	request.Header.Set("Authorization", bearerToken)
	// request.Header.Set("Content-Type", "video/mp4")
	request.Header.Set("Content-Type", "multipart")

	//make the post request
	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer response.Body.Close()

	// //read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Print the response status and body
	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Body:", string(responseBody))

}

func uploadObjectIntoBucket2() {
	url := fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/ytc-media-storage/o?uploadType=multipart&name=%s", "pets%2Fdog2.mp4")
	// url := "https://storage.googleapis.com/upload/storage/v1/b/ytc-media-storage/o?uploadType=multipart"
	fmt.Println(url)

	file, err := os.Open(config.VIDEO_PATH)
	if err != nil {
		log.Fatalf("Error reading object data: %v", err)
	}

	buffer := bytes.Buffer{}
	writer := multipart.NewWriter(&buffer)

	// Add JSON metadata part
	metadataPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"application/json; charset=UTF-8"},
	})

	if err != nil {
		log.Fatalf("Error creating the metadata part: %v", err)
		return
	}

	// Replace OBJECT_METADATA with your actual metadata
	_, err = metadataPart.Write([]byte(`{"Samplekey": "Samplevalue"}`))
	if err != nil {
		log.Fatalf("Error writing metadata: %v", err)
		return
	}

	// Add video file part
	filePart, err := writer.CreateFormFile("file", "SampleFile")

	if err != nil {
		log.Fatalf("Error creating the file part: %v", err)
		return
	}

	_, err = io.Copy(filePart, file)
	if err != nil {
		log.Fatalf("Error copying the file part: %v", err)
		return
	}

	// Close the writer
	writer.Close()

	request, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	token, _ := generateJWTToken()
	bearerToken := "Bearer " + token
	fmt.Println(bearerToken)

	request.Header.Set("Authorization", bearerToken)
	request.Header.Set("Content-Type", "multipart/related")
	request.Header.Set("Content-Length", strconv.Itoa(buffer.Len())) // Set the Content-Length header

	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Body:", string(responseBody))
}

func generateJWTToken() (string, error) {
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
	serviceAccountData, err := os.ReadFile(config.ServiceAccountPath)
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

func createAuthToken() {
	ClientId, ClientSecret := config.LoadClientAndSecretKey("C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\YoutubeConnector\\config\\client_secret.json")
	fmt.Println(ClientId)
	fmt.Println(ClientSecret)
	var (
		clientID     = ClientId
		clientSecret = ClientSecret
		redirectURL  = config.REDIRECT_URL
		scopes       = []string{"https://www.googleapis.com/auth/youtube.upload"}
	)

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}

	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {

		// Extract the authorization code from the URL query parameters
		code := r.URL.Query().Get("code")

		// Exchange authorization code for an access token
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			fmt.Printf("Error exchanging code for token: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		//Use the access token for API requests
		res, _ := json.Marshal(token)
		w.Write(res)
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusOK)

		fmt.Printf("Access Token: %v\n", token.AccessToken)

		// if token.Expiry.Before(time.Now()) {
		// 	// token is expired
		// 	newToken, err := config.TokenSource(context.Background(), token).Token()
		// 	if err != nil {
		// 		fmt.Printf("Error refreshing token: %v\n", err)
		// 		return
		// 	}

		// 	token = newToken
		// 	fmt.Println(newToken)
		// }

		// Handle the response or redirect as needed
		fmt.Fprintf(w, "Authorization successful! You can close this window.")

	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func callUpload() {
	title := "My Video Title"
	description := "Description of my video."
	privacyStatus := "public" // Set to "public" for a public video

	// YouTube API endpoint
	apiEndpoint := config.ROOT_URL + "/" + "upload/youtube/v3/videos"

	// Path to the video file to upload
	videoFilePath := "C:\\Users\\AKumar22\\Downloads\\testvd.mp4"

	// YouTube API key or OAuth 2.0 access token (replace with your actual key or token)
	apiKeyOrAccessToken := config.OAUTH_TOKEN_KR8799

	// Create a new multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add video metadata as form fields
	writer.WriteField("snippet", fmt.Sprintf(`{"title":"%s","description":"%s","privacyStatus":"%s"}`, title, description, privacyStatus))

	// Add the video file
	file, err := os.Open(videoFilePath)
	if err != nil {
		fmt.Printf("Error opening video file: %v\n", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("videoFile", file.Name())
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		return
	}

	// Close the multipart writer
	writer.Close()

	// Create a POST request with the multipart body
	request, err := http.NewRequest("POST", apiEndpoint, body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Set the content type header for multipart requests
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Add the YouTube API key or OAuth 2.0 access token to the request
	request.Header.Set("Authorization", "Bearer "+apiKeyOrAccessToken)

	// Perform the request
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Read and print the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println("Response:", string(responseBody))
}

func uploadVideo(jsonPayloadData string) {
	// title := "Sample First Video"
	// description := "Keyboard | Mouse | MousePad"
	// privacyStatus := "public" // Set to "public" for a public video

	// url := config.ROOT_URL + config.UTUBE_END_ENDPOINT + config.UPLOAD_ENDPOINT
	url := config.ROOT_URL + "/" + "upload/youtube/v3/videos"

	url += "?part=snippet%2Cstatus&key=" + config.API_Key

	// "https://youtube.googleapis.com/youtube/v3/videos?part=snippet%2Cstatus&key="
	// url += "&key=" + config.API_Key
	fmt.Println("Calling", url)
	fmt.Println(jsonPayloadData)

	// Create a new buffer to store the multipart/form-data payload
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	// Add the JSON payload as a form field
	jsonPart, err := writer.CreateFormField("json")
	if err != nil {
		fmt.Printf("Error creating JSON form field: %v\n", err)
		return
	}
	jsonPart.Write([]byte(jsonPayloadData))

	// Create a map representing the "snippet" field
	// snippetData := map[string]string{
	// 	"title":         title,
	// 	"description":   description,
	// 	"privacyStatus": privacyStatus,
	// }

	/*
		// Convert the map to a JSON string
		snippetJSON, err := json.Marshal(snippetData)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		// Use the JSON string in the writer.WriteField function
		writer.WriteField("snippet", string(snippetJSON))
	*/

	// Add the video file
	file, err := os.Open(config.VIDEO_PATH)
	if err != nil {
		fmt.Println("Error opening the video file", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("videoFile", file.Name())
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		return
	}

	// Close the multipart writer
	writer.Close()

	request, err := http.NewRequest(config.POST_METHOD, url, buffer)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	bearer_token := "Bearer " + config.OAUTH_TOKEN

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", bearer_token)

	// Make the POST request

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Print the response status and body
	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Body:", string(responseBody))
}

func getVideo() {
	authenticationMode := config.AUTH_MODE
	rootURL := config.ROOT_URL
	apiKey := config.API_Key

	videoId := "8wx_bxtBQQ0"
	var getVideoEndPoint string = rootURL + config.UTUBE_END_ENDPOINT

	getVideoEndPoint += "?id=" + videoId
	switch authenticationMode {
	case "API_TOKEN":
		getVideoEndPoint += "&key=" + apiKey
	}
	getVideoEndPoint += "&part=snippet,contentDetails,statistics,status"

	fmt.Println(getVideoEndPoint)
	res, err := http.Get(getVideoEndPoint)

	if err != nil {
		fmt.Println(err)
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	// jsondata, err := json.Marshal(body)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Println("Status ", res.StatusCode)

	fmt.Println(res)
}

func uploadThumbnail() {

	// Replace with the video ID for which you want to set the thumbnail
	videoID := "6Qvf-yrzFdE"

	// Replace with the path to the image file you want to use as the thumbnail
	imageFilePath := "C:\\Users\\AKumar22\\Pictures\\Screenshots\\sc.png"

	// Open the image file
	imageFile, err := os.Open(imageFilePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return
	}
	defer imageFile.Close()

	// Create a buffer to store the image file data
	imageBuffer := &bytes.Buffer{}
	_, err = imageBuffer.ReadFrom(imageFile)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return
	}

	// Create a new form data writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the videoId parameter
	_ = writer.WriteField("videoId", videoID)

	// Add the image file to the form data
	part, err := writer.CreateFormFile("media", "thumbnail.jpg")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = part.Write(imageBuffer.Bytes())
	if err != nil {
		fmt.Println("Error writing image file data:", err)
		return
	}

	// Close the form data writer
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing form data writer:", err)
		return
	}

	// Create the POST request
	// url := fmt.Sprintf("https://www.googleapis.com/upload/youtube/v3/thumbnails/set?videoId=%s&key=%s", videoID, apiKey)

	url := fmt.Sprintf("https://www.googleapis.com/upload/youtube/v3/thumbnails/set?videoId=%s", videoID)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type header
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", "Bearer "+config.OAUTH_TOKEN_KR4DRI)

	// Make the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response
	fmt.Println("Response:", string(responseBody))
}
