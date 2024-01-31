package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abhis2007/YOUTUECONNECTOR/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

func Index(w http.ResponseWriter, r *http.Request) {

	parsed, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = parsed.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upload Endpoint - Post request")
	parsed, err := template.ParseFiles("./templates/upload.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = parsed.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

}

type snippetBody struct {
	Title       string   `json:"title"`
	CategoryId  string   `json:"categoryId"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type statusBody struct {
	PrivacyStatus string `json:"privacyStatus"`
}

type mainBody struct {
	Id      string      `json:"id"`
	Snippet snippetBody `json:"snippet"`
	Status  statusBody  `json:"status"`
}

func updateVideo(bodyArgs string, videoId string) {
	accessToken := config.OAUTH_TOKEN_KR8799
	uploadEndpoint := "https://youtube.googleapis.com/youtube/v3/videos?part=snippet&part=status"
	// fmt.Println(uploadEndpoint)

	//Create the request
	request, err := http.NewRequest("PUT", uploadEndpoint, strings.NewReader(bodyArgs))

	if err != nil {
		return
	}

	// Add the YouTube API key or OAuth 2.0 access token to the request
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	//return
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	fmt.Println("full response:", response)
	defer response.Body.Close()

	// Read and print the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println("Response:", string(responseBody))

}

func callUpload(snippet string) {
	fmt.Println("callupload")
	fmt.Println(snippet)

	// YouTube API endpoint
	apiEndpoint := config.ROOT_URL + "/" + "upload/youtube/v3/videos"
	//apiEndpoint += "?part=snippet,status"

	// Path to the video file to upload
	videoFilePath := "C:\\Users\\AKumar22\\Downloads\\testvd.mp4"

	// YouTube API key or OAuth 2.0 access token (replace with your actual key or token)
	apiKeyOrAccessToken := config.OAUTH_TOKEN_KR8799

	// Create a new multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	var fw io.Writer

	//adding the non-video data
	fw, err := writer.CreateFormField("snippet")
	if err != nil {
		fmt.Println("error while writinf the data : ", err)
		return
	}
	// _, err = io.Copy(fw, strings.NewReader(snippet))
	_, err = io.Copy(fw, bytes.NewReader([]byte(snippet)))
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Snippet data copied.")
	}

	// Add the video file
	file, err := os.Open(videoFilePath)
	if err != nil {
		fmt.Printf("Error opening video file: %v\n", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		return
	} else {
		fmt.Println("File data copied.")
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
	//return
	// Perform the request
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	fmt.Println("full response:", response)
	defer response.Body.Close()

	// Read and print the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println("Response:", string(responseBody))
}

type VideoRequest struct {
	FileLocation string `json:"fileLocation"`
}

// fINAL AND TESTED FUNCTION WORKING OK FOR UPLOAD THE OBJECT INTO THE CLOUD STORAGE BUCKET
func test(videoPath string) {
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
	//fmt.Println(key.AccessToken)

	baseFilePart := filepath.Base(videoPath)
	originalFileName := strings.TrimSuffix(baseFilePart, ".mp4")
	fmt.Println(originalFileName)

	file, err := os.ReadFile(videoPath)

	if err != nil {
		log.Fatalf("Error reading object data: %v", err)
	}

	//Form the endpoints
	url := fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/ytc-media-storage/o?uploadType=media&name=%s", originalFileName)
	//fmt.Println(url)

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

func sendError(w http.ResponseWriter, errorMessage string, statusCode int) {
	http.Error(w, errorMessage, statusCode)
	//w.Write([]byte(errorMessage))
}

// tested and working - request will be from the user.
func UploadVideoOnStorageServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	file, handler, err := r.FormFile("userFile")
	if err != nil {
		sendError(w, "", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := isValidFile(handler); err != nil {
		fmt.Println(err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	originalFileName := strings.TrimSuffix(handler.Filename, fileExt)

	uniqueFileName := originalFileName + "_" + time.Now().Format("20060102_150405") + fileExt
	destnationFileLocation := config.STATIC_FILE_PATH + "//" + uniqueFileName
	//fmt.Println(destnationFileLocation)
	dst, err := os.Create(destnationFileLocation)
	if err != nil {
		sendError(w, "Error in accessing the path", http.StatusBadRequest)
		return
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		sendError(w, "Failed to open the file content", http.StatusNotFound)
		return
	}
	test(destnationFileLocation)

}

func isValidFile(filepart *multipart.FileHeader) error {
	file, err := filepart.Open()
	if err != nil {
		return err
	}
	extn := filepath.Ext(filepart.Filename)
	if extn != ".mp4" {
		return fmt.Errorf("invalid file extension. Only .mp4 files are allowed")
	}

	buffer := make([]byte, 512) // Only read the first 512 bytes to detect MIME type
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimeType, "video/") {
		return fmt.Errorf("invalid file type. Only video files are allowedi")
	}

	return nil

}

// Tested code for update the video metada after the successfull upload of the video.
func FetchAndUploadVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	type formValue struct {
		Title           string   `json:"title"`
		Description     string   `json:"description"`
		Category        string   `json:"category"`
		Audience        string   `json:"audience"`
		AgeRestrictions string   `json:"ageRestrictions"`
		TagInput        []string `json:"tagInput"`
		Privacy         string   `json:"privacy"`
	}
	var (
		d       formValue
		videoId = "zHA-YusgSSo"
	)

	json.NewDecoder(r.Body).Decode(&d)
	snippetParam := snippetBody{
		Title:       d.Title,
		CategoryId:  "22",
		Description: d.Description,
		Tags:        d.TagInput,
	}
	statusParam := statusBody{
		PrivacyStatus: d.Privacy,
	}
	bodyParams := mainBody{
		Id:      videoId,
		Snippet: snippetParam,
		Status:  statusParam,
	}
	jsonData, err := json.Marshal(bodyParams)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
	updateVideo(string(jsonData), videoId)
}

func FetchAndUploadVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	type formValue struct {
		Title           string   `json:"title"`
		Description     string   `json:"description"`
		Category        string   `json:"category"`
		Audience        string   `json:"audience"`
		AgeRestrictions string   `json:"ageRestrictions"`
		TagInput        []string `json:"tagInput"`
		Privacy         string   `json:"privacy"`
	}
	// type snippetData struct {
	// 	Snippet formValue `json:snippet`
	// }
	var d formValue

	json.NewDecoder(r.Body).Decode(&d)

	// data := snippetData {
	// 	Snippet: d,
	// }
	fmt.Println(d.Title)
	fmt.Println(d.Description)
	// fmt.Println(d.Category)
	// fmt.Println(d.Audience)
	// fmt.Println(d.AgeRestrictions)
	// fmt.Println(d.TagInput)
	// fmt.Println(d.Privacy)

	// if any input is empty amd is mandatory then returen below some info.
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(d)
	//return
	// snippet := {
	// 	title       string   `json:"title"`
	// 	description string   `json:"description"`
	// 	tags        []string `json:"TagInput"`
	// 	categoryId  string   `json:"Category"`
	// }
	// Convert tags array to JSON
	tagsJSON, err := json.Marshal(d.TagInput)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	snippet := fmt.Sprintf(`{"title":"%s","description":"%s","tags":%s, "categoryId":"%s"}`, d.Title, d.Description, tagsJSON, "22")
	data := fmt.Sprintf(`{"snippet" : %s}`, snippet)
	//fmt.Println(snippet)
	callUpload(data)
	// fmt.Println(snippet)
}

func Videos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("videos")
	// LogEntry represents a log entry.
	type LogEntry struct {
		Level      string
		Message    string
		ResourceId string

		TraceId          string
		SpanId           string
		Commit           string
		ParentResourceId string
	}
	type PageData struct {
		Results []LogEntry
	}
	var logEntries []LogEntry
	for i := 1; i <= 9; i++ {
		logEntries = append(logEntries, LogEntry{
			Level:      "Info",
			Message:    "Log message",
			ResourceId: "123",

			TraceId:          "456",
			SpanId:           "789",
			Commit:           "abc",
			ParentResourceId: "xyz",
		})
	}

	parsed, err := template.ParseFiles("./templates/videos.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = parsed.Execute(w, PageData{Results: logEntries})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func uploadTOGCS() {

}
