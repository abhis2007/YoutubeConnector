package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/abhis2007/YOUTUECONNECTOR/config"
	"github.com/abhis2007/YOUTUECONNECTOR/routes"
	_ "github.com/denisenkom/go-mssqldb"
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
	config.DbInit()
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
		//uploadObjectOnGCS("C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\sample.mp4")
		//createAuthToken()
		//downloadObject()
		//test2("")

		//insertdata("to_tony_20240202_234305")
		//deleteObjectOnGCS()
		fmt.Println("Server started at port : 8080")

		// test()
		log.Fatal(http.ListenAndServe(":8080", routerInstance))

	}

	// test()
	// uploadObjectIntoBucket2()
	//updateObject()

}

// Tested code
func deleteObjectOnGCS() {
	url := "https://storage.googleapis.com/storage/v1/b/ytc-media-storage/o/to_tony_20240204_075052"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}
	serviceAccountData, err := os.ReadFile(config.ServiceAccountPath)
	if err != nil {
		log.Fatalf("Error extacting the service acc: %v", err)
	}
	configToken, err := google.JWTConfigFromJSON(serviceAccountData, storage.DevstorageFullControlScope)
	if err != nil {
		log.Fatalf("Error creating JWT Config: %v", err)
	}
	key, val := configToken.TokenSource(context.Background()).Token()
	if val != nil {
		fmt.Println(err)
	}

	// Set headers, including the Authorization header with the JWT token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key.AccessToken)

	// Create an HTTP client with OAuth2 authentication
	//client := configToken.Client(context.Background())

	client := &http.Client{}
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
}

/*
//Still need to know can we do the override of the content inside the bucket object.

func updateVideoOnGCS() {
	fmt.Println("Called UpdateVideoOnGCS")
	videoPath := "C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\sample.mp4"
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
	// key, val := configToken.TokenSource(context.Background()).Token()
	// if val != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(key.AccessToken)

	// baseFilePart := filepath.Base(videoPath)
	// originalFileName := strings.TrimSuffix(baseFilePart, ".mp4")
	// fmt.Println(originalFileName)

	file, err := os.ReadFile(videoPath)
	file = file

	if err != nil {
		log.Fatalf("Error reading object data: %v", err)
	}

	//Form the endpoints
	// objectId := "to_tony_20240204_065217"
	// url := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/ytc-media-storage/o/%s", objectId)
	url := "https://storage.googleapis.com/storage/v1/b/ytc-media-storage/o/to_tony_20240204_065217/rewriteTo/b/ytc-media-storage/o/to_tony_20240204_065217"
	fmt.Println(url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	// req.Header.Set("Authorization", "Bearer "+key.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)

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
*/

// Probably in future this should go off - then we should be using config.DBInit() function rather than below setupDb()
func setupDb() {
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
	config.DB = db

	// Create the table
	createTable := config.ObjectTblCreation
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(config.UserTableCreation, "userTbl")
	if err != nil {
		log.Fatal(err)
	}
}

// Tested code - Now we can authenticate this with using the service account.
func downloadObject() {

	//extract the token from service account
	tokenKey, err := generateJWTToken()

	if err != nil {
		log.Fatalf("Error extacting the service acc: %v", err)
	}

	bucketId := "ytc-media-storage"
	url := "https://storage.googleapis.com/" + bucketId + "/"
	url += "abhis2007" + "/to_tony_20240203_175457"

	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)

	// Set headers, including the Authorization header with the JWT token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenKey)

	if err != nil {
		fmt.Println("Error in making the request")
		return
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(res)
		return
	}

	if res.StatusCode/100 != 2 {
		fmt.Println(res)
		return
	}

	file, err := os.Create("C:\\Users\\AKumar22\\Desktop\\StudyContents\\GoLang\\YoutubeConnector\\static\\videos\\today.mp4")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println(err)
	}
}

// Update the object inside the bucket
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
// func test(videoPath string) {
func uploadObjectOnGCS(videoPath string) {

	//extract the token from service account
	tokenKey, err := generateJWTToken()
	if err != nil {
		log.Fatalf("Error creating JWT Config: %v", err)
	}

	// Add the video file
	// filePath := config.VIDEO_PATH
	filePath := videoPath

	baseFilePart := filepath.Base(filePath)
	lists := strings.Split(baseFilePart, ".")
	if len(lists) <= 0 {
		log.Fatalf("File name doesnot seems to be of a media type\n")
		return
	}
	var fileName string = lists[0]

	file, err := os.ReadFile(filePath)

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

	req.Header.Set("Authorization", "Bearer "+tokenKey)
	req.Header.Set("Content-Type", "video/mp4")

	// Make the request
	client := http.Client{}
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

func generateJWTToken() (string, error) {
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
		//scopes       = []string{"https://www.googleapis.com/auth/youtube.upload"}
		scopes = []string{
			"https://www.googleapis.com/auth/youtube",
			"https://www.googleapis.com/auth/youtube.force-ssl",
			"https://www.googleapis.com/auth/cloud-platform",
		}
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
