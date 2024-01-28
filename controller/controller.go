package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/abhis2007/YOUTUECONNECTOR/config"
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

func callUpload(snippet string) {
	fmt.Println("callupload")
	// title := "My Video Title"
	// description := "Description of my video."
	// privacyStatus := "public" // Set to "public" for a public video

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
	// writer.WriteField("snippet", fmt.Sprintf(`{"title":"%s","description":"%s","privacyStatus":"%s"}`, title, description, privacyStatus))

	writer.WriteField("snippet", snippet)
	fmt.Println(snippet)
	// return

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

// type Snippet struct {
// 	title       string   `json:"title"`
// 	description string   `json:"description"`
// 	tags        []string `json:"TagInput"`
// 	categoryId  string   `json:"Category"`
// }

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
	var d formValue
	json.NewDecoder(r.Body).Decode(&d)
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
	fmt.Println(snippet)
	callUpload(snippet)
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
