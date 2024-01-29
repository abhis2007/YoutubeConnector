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
	"strings"

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

type snippetBody struct {
	Title       string `json:"title"`
	CategoryId  string `json:"categoryId"`
	Description string `json:"description"`
}

type mainBody struct {
	Id      string      `json:"id"`
	Snippet snippetBody `json:"snippet"`
}

func updateVideo(snippet string, videoId string) {
	fmt.Println(videoId)
	accessToken := config.OAUTH_TOKEN_KR8799
	uploadEndpoint := "https://youtube.googleapis.com/youtube/v3/videos?part=snippet"

	fmt.Println(uploadEndpoint)

	snippetParam := snippetBody{
		Title:       "title_new_struct",
		CategoryId:  "22",
		Description: "NeeDesc",
	}
	bodyParams := mainBody{
		Id:      videoId,
		Snippet: snippetParam,
	}
	jsonData, err := json.Marshal(bodyParams)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	
	fmt.Println(string(jsonData))
	

	//Create the request
	request, err := http.NewRequest("PUT", uploadEndpoint, strings.NewReader(string(jsonData)))

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

func FetchAndUploadVideo(w http.ResponseWriter, r *http.Request) {

	type VideoSnippet struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		CategoryId  string   `json:"categoryId"`
		Tags        []string `json:"tags,omitempty"`
	}
	type VideoStatus struct {
		PrivacyStatus string `json:"privacyStatus"`
	}
	type Video struct {
		Snippet *VideoSnippet `json:"snippet"`
		Status  *VideoStatus  `json:"status"`
	}

	title := "Sample Title"
	description := "Sample Description"
	category := "24"
	privacy := "unlisted"

	// Create Video instance
	upload := &Video{
		Snippet: &VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  category,
			Tags:        []string{"tag1", "tag2"},
		},
		Status: &VideoStatus{
			PrivacyStatus: privacy,
		},
	}

	// Marshal the struct into JSON
	jsonData, err := json.Marshal(upload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Print the resulting JSON string
	//fmt.Println(string(jsonData))
	//callUpload(string(jsonData))
	updateVideo(string(jsonData), "zHA-YusgSSo")
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
