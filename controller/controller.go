package controller

import (
	"bytes"
	"database/sql"
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
	"github.com/google/uuid"
)

type Session struct {
	UserId string
	Expiry time.Time
}

var Sessions = map[string]Session{}

type Credential struct {
	UserName string
	Password string
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reached")
	parseTemplate, err := template.ParseFiles("./templates/users.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	db := config.DB
	var sqlQuery string

	if config.LoginUserProfile != "Admin" {
		sqlQuery += "select AdminId from adminMapTbl where userId = '" + config.USERID + "'"
	} else {
		sqlQuery = "select userId from adminMapTbl where AdminId = '" + config.USERID + "'"
	}

	fmt.Println(sqlQuery)

	rows, err := db.Query(sqlQuery)
	if err != nil {
		fmt.Println(err)
	}

	var userIds []string
	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			fmt.Println(err)
		}
		userIds = append(userIds, userId)
	}

	err = parseTemplate.Execute(w, userIds)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func AssignTaskToCoUser(w http.ResponseWriter, r *http.Request) {

	if config.LoginUserProfile == "Admin" {
		UploadVideo(w, r)
	}
}

func RenderVideoReviewTemplate(w http.ResponseWriter, r *http.Request) {
	//This should only be trigger at the time of video review not directly.

	type dummy struct {
		arrayData string
	}

	type propertyKeyValuePair struct {
		PropertyName       string
		PropertyValue      string
		propertyDefaultKey []dummy
	}

	type propertyAndKeyvalue struct {
		PropertyKeyValue []propertyKeyValuePair
	}

	query := "select * from metadataFiledTbl where userId = '" + config.USERID + "' AND AdminId = '" + "Abhis2007'"
	fmt.Println(query)
	db := config.DB
	row, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	propertyKeyValues := []propertyKeyValuePair{}
	for row.Next() {
		var Userid, adminId, Title, Description, Category, Audience, AgeRestrictions, TagInput, Privacy string
		err := row.Scan(&Userid, &adminId, &Title, &Description, &Category, &Audience, &AgeRestrictions, &TagInput, &Privacy)
		fmt.Println(err)
		fmt.Println(Title)
		//var nilArr []string
		var categoryArr []string
		for i := 0; i <= 10; i++ {
			categoryArr = append(categoryArr, "category_1")
		}

		sample := []dummy{}
		sample = append(sample, dummy{"data1"})
		sample = append(sample, dummy{"data1"})
		sample = append(sample, dummy{"data1"})
		sample = append(sample, dummy{"data1"})
		sample = append(sample, dummy{"data1"})
		sample = append(sample, dummy{"data1"})
		sample = append(sample, dummy{"data1"})

		propertyKeyValues = append(propertyKeyValues, propertyKeyValuePair{"Title", Title, sample})
		propertyKeyValues = append(propertyKeyValues, propertyKeyValuePair{"Description", Description, sample})
		propertyKeyValues = append(propertyKeyValues, propertyKeyValuePair{"Category", Category, sample})
		propertyKeyValues = append(propertyKeyValues, propertyKeyValuePair{"Audience", Audience, sample})
		propertyKeyValues = append(propertyKeyValues, propertyKeyValuePair{"AgeRestrictions", AgeRestrictions, sample})
		propertyKeyValues = append(propertyKeyValues, propertyKeyValuePair{"TagInput", TagInput, sample})
		propertyKeyValues = append(propertyKeyValues, propertyKeyValuePair{"Privacy", Privacy, sample})
	}

	// fmt.Println(propertyKeyValues)

	propertyKeyValueArr := propertyAndKeyvalue{propertyKeyValues}
	propertyKeyValueArr = propertyKeyValueArr

	parseTemplate, err := template.ParseFiles("./templates/videoReview.html")
	if err != nil {
		fmt.Println(err)
	}

	parseTemplate.Execute(w, propertyKeyValuePair{PropertyName: "property", PropertyValue: "val", propertyDefaultKey: []dummy{}})

}

type Property struct {
	Name        string
	Value       string
	StringArray []string
}

type Data struct {
	Properties []Property
}

var categoryList []string

func Test(w http.ResponseWriter, r *http.Request) {
	//Assuming profile i non-Admin
	userHandle := r.URL.Query().Get("id")
	config.UserRequest = userHandle
	var nonAdminUserId string
	var AdminUserId string

	//Assume came from nonadmin account
	AdminUserId = userHandle
	nonAdminUserId = config.USERID

	if config.LoginUserProfile == "Admin" { //override if not
		nonAdminUserId = userHandle
		AdminUserId = config.USERID
	}

	query := "select * from metadataFiledTbl where userId = '" + nonAdminUserId + "' AND AdminId = '" + AdminUserId + "'"
	fmt.Println(query)
	db := config.DB
	row, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	properties := []Property{}
	categoryList = append(categoryList, "category1")
	categoryList = append(categoryList, "category2")
	categoryList = append(categoryList, "category3")
	categoryList = append(categoryList, "category4")
	categoryList = append(categoryList, "category5")

	for row.Next() {
		var uid, Userid, adminId, Title, Description, Category, Audience, AgeRestrictions, TagInput, Privacy string
		err := row.Scan(&uid, &Userid, &adminId, &Title, &Description, &Category, &Audience, &AgeRestrictions, &TagInput, &Privacy)
		fmt.Println(err)
		if Title == "" || Description == "" || Category == "" || Audience == "" || AgeRestrictions == "" || Privacy == "" {
			continue
		} else {
			query = "select objectPathOnServer from objUserAndAdminAMpTbl where AdminId = '" + AdminUserId + "' And userId = '" + nonAdminUserId + "'"
			fmt.Println(query)
			var path string
			db.QueryRow(query).Scan(&path)
			property1 := []Property{
				{"Title", Title, []string{}},
				{"Description", Description, []string{}},
				{"Category", Category, categoryList},
				{"Audience", Audience, []string{"Yes, it's 'Made for Kids'", "No, it's not 'Made for Kids'"}},
				{"AgeRestrictions", AgeRestrictions, []string{"No, don't restrict my video to viewers over 18 only", "Yes, restrict my video to viewers over 18 Age"}},
				{"TagInput", TagInput, []string{}},
				{"Privacy", Privacy, []string{"Private", "Unlisted", "Public"}},
				{"Path", path, []string{}},
			}
			properties = append(properties, property1...)
			break
		}

	}

	data := Data{Properties: properties}
	fmt.Println(data)
	parseTemplate, err := template.ParseFiles("./templates/videoReview.html")
	if err != nil {
		fmt.Println(err)
	}

	parseTemplate.Execute(w, data)

}

func HandleUserDetails(w http.ResponseWriter, r *http.Request) {
	// if config.LoginUserProfile != "Admin" {
	Test(w, r)
	return
	// }
	userHandle := r.URL.Query().Get("id")
	fmt.Println(userHandle)
	var tmplName string
	categoryList = append(categoryList, "category1")
	categoryList = append(categoryList, "category2")
	categoryList = append(categoryList, "category3")
	categoryList = append(categoryList, "category4")
	categoryList = append(categoryList, "category5")

	if config.LoginUserProfile != "Admin" {
		//thne userId is the id of the admin
		// and my login id is config.userid
		query := "select * from metadataFiledTbl where userId = '" + config.USERID + "' AND AdminId = '" + userHandle
		fmt.Println(query)
		db := config.DB
		var uid, Userid, adminId, Title, Description, Category, Audience, AgeRestrictions, TagInput, Privacy string
		err := db.QueryRow(query).Scan(&uid, &Userid, &adminId, &Title, &Description, &Category, &Audience, &AgeRestrictions, &TagInput, &Privacy)
		//Assuming only active job/task will be there by the admin to me
		if err != nil {
			fmt.Println(err)
		}

		property1 := []Property{
			{"Title", Title, []string{}},
			{"Description", Description, []string{}},
			{"Category", Category, categoryList},
			{"Audience", Audience, []string{"It's 'Made for Kids", "It's not 'Made for Kids'"}},
			{"AgeRestrictions", AgeRestrictions, []string{"Yes, restrict my video to viewers over 18 A", "No, don't restrict my video to viewers over 18 only"}},
			{"TagInput", TagInput, []string{}},
			{"Privacy", Privacy, []string{"Option A", "Option B"}},
		}
		properties := []Property{}
		properties = append(properties, property1...)
		parseTemplate, err := template.ParseFiles("./templates/videoReview.html")
		if err != nil {
			fmt.Println(err)
		}
		parseTemplate.Execute(w, properties)

	} else {
		//go for the review page
		tmplName = "videoReview.html"
	}

	// var query string
	// query = "select objectPathOnServer from objUserAndAdminAMpTbl"
	// pf := config.LoginUserProfile
	// if pf != "Admin" {
	// 	//here userId will be of the admin who created your account .
	// 	query += " where AdminId = '" + userId + "' AND userId = '" + config.USERID + "'"
	// } else {
	// 	query += " where userId = '" + userId + "' AND AdminId = '" + config.USERID + "'"
	// }
	// fmt.Println(query)

	// db := config.DB

	// rows, err := db.Query(query)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// var paths []string
	// for rows.Next() {
	// 	var path string
	// 	if err := rows.Scan(&path); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	paths = append(paths, path)
	// }

	// tmplName = "co-userTask.html"

	tmplName = "videoReview.html"

	// if config.LoginUserProfile != "Admin" {
	// 	tmplName = "videoReview.html"
	// }
	tmplName = "./templates/" + tmplName

	parseTemplate, err := template.ParseFiles(tmplName)
	if err != nil {
		fmt.Println(err)
	}

	// type TemplateData struct {
	// 	AllPath []string
	// 	UserId  string
	// }

	// data := TemplateData{
	// 	AllPath: paths,
	// 	UserId:  userId,
	// }

	err = parseTemplate.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credential

	json.NewDecoder(r.Body).Decode(&creds)

	expPass := "123456789"
	if creds.Password != expPass {
		fmt.Println("Account not found")
		return
	}

	fmt.Println(creds.UserName)
	fmt.Println(creds.Password)

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)
	fmt.Println(expiresAt, time.Now().Add(2*time.Second))

	Sessions[sessionToken] = Session{
		UserId: creds.UserName,
		Expiry: expiresAt,
	}

	curUserCookie := &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}
	http.SetCookie(w, curUserCookie)
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) int {
	fmt.Println("Authentication called")
	m_cookie, err := r.Cookie("session_token")
	fmt.Println(err)
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("Coockie not found")
			http.Error(w, "Coockie not found", http.StatusUnauthorized)
			return http.StatusUnauthorized
		}
		fmt.Println("bad request")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return http.StatusBadRequest
	}
	sessionToken := m_cookie.Value
	userSession, ok := Sessions[sessionToken]

	if !ok {
		fmt.Println("Token not found in the map")
		http.Error(w, "Token not found in the map", http.StatusUnauthorized)
		return http.StatusUnauthorized
	}
	if userSession.isExpired() {
		delete(Sessions, sessionToken)
		fmt.Println("Token expired")
		http.Error(w, "Token expired", http.StatusUnauthorized)
		return http.StatusUnauthorized
	}
	fmt.Println("Is successful", userSession.UserId)
	return http.StatusOK
}

func (s Session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func Logout(w http.ResponseWriter, r *http.Request) {
	m_cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Cookie not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	sessionToken := m_cookie.Value
	delete(Sessions, sessionToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if AuthenticateUser(w, r) != http.StatusOK {
		return
	}
	m_cookie, _ := r.Cookie("session_token")
	oldSessionToken := m_cookie.Value
	userSession := Sessions[oldSessionToken]

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	Sessions[newSessionToken] = Session{
		UserId: userSession.UserId,
		Expiry: expiresAt,
	}

	delete(Sessions, oldSessionToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    newSessionToken,
		Expires:  time.Now().Add(120 * time.Second),
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func LoginData(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for the entire request
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	var userID string = r.FormValue("userId")
	var password string = r.FormValue("password")

	if userID == "" {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return
	}

	if len(password) <= 6 {
		http.Error(w, "Password doesn't meet the requirements.", http.StatusBadRequest)
		return
	}
	db := config.DB
	query := `SELECT userId, profileType FROM userTbl WHERE userId = @p1 AND password = @p2`
	var d_userId, d_profileType string
	err = db.QueryRow(query, userID, password).Scan(&d_userId, &d_profileType)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusBadRequest)
			return
		} else {
			// Handle other errors
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Printf("User found.\nUserID: %s\nProfileType: %s\n", d_userId, d_profileType)
	}
	http.Error(w, "Login successful", http.StatusAccepted)

	//In future it should go off for sure
	config.USERID = d_userId
	config.LoginUserProfile = d_profileType

	fmt.Println(d_userId)

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	Sessions[sessionToken] = Session{
		UserId: d_userId,
		Expiry: expiresAt,
	}

	fmt.Println("Setting up the cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SignupData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for the entire request
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	// Access form fields
	var userID string = r.FormValue("userId")
	var password string = r.FormValue("password")
	var profileType string = r.FormValue("profileType")
	if profileType == "Co-Admin" && config.USERID == "" {
		http.Error(w, "Please login first to create a co-user account", http.StatusBadRequest)
		return
	}

	if config.LoginUserProfile == "Co-Admin" {
		http.Error(w, "Server couldn't process the request, please login as a Admin User to create a co-user Account.", http.StatusBadRequest)
		return
	}
	fmt.Println("User ID:", userID)
	fmt.Println("Password:", password)
	fmt.Println("Profile", profileType)

	if userID == "" {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return
	}

	if len(password) <= 6 {
		http.Error(w, "Password doesn't meet the requirements.", http.StatusBadRequest)
		return
	}
	db := config.DB
	query := `SELECT count(userId) FROM userTbl WHERE userId = @p1`
	var count int
	err = db.QueryRow(query, userID, password).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		insertSql := `insert into userTbl (userId, password, profileType) VALUES (@P1, @P2, @p3)`
		db.Exec(insertSql, userID, password, profileType)

		if profileType == "Co-Admin" {
			insertSql = `insert into adminMapTbl(userId, AdminId) Values (@p1, @p2)`
			db.Exec(insertSql, userID, config.USERID)
		}
	} else {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	http.Error(w, "Account created successfully.", http.StatusAccepted)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	parseTemplate, err := template.ParseFiles("./templates/signup.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = parseTemplate.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	parseTemplate, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = parseTemplate.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {

	parsed, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	profile := config.LoginUserProfile
	err = parsed.Execute(w, profile)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upload Endpoint - Post request")
	query := "select userId from adminMapTbl where AdminId = '" + config.USERID + "'"
	fmt.Println(query)

	db := config.DB
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	var userIds []string
	for rows.Next() {
		var userId string
		if rows.Scan(&userId) != nil {
			fmt.Println(err)
		} else {
			userIds = append(userIds, userId)
		}
	}
	fmt.Println(userIds)
	parsed, err := template.ParseFiles("./templates/upload.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = parsed.Execute(w, userIds)
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
	PrivacyStatus           string `json:"privacyStatus"`
	SelfDeclaredMadeForKids bool   `json:"selfDeclaredMadeForKids"`
}

type mainBody struct {
	Id      string      `json:"id"`
	Snippet snippetBody `json:"snippet"`
	Status  statusBody  `json:"status"`
}

func updateVideoMetadataOnYoutubeServer(bodyArgs string, videoId string) {
	fmt.Println("Calling updateVideoMetadataOnYoutubeServer")
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

// This endpoint upload the video for the first time (callUpload(snippet string))) -> uploadVideoOnYoutubeServer
func uploadVideoOnYoutubeServer(snippet string) (error, string) {
	fmt.Println("Running uploadVideoOnYoutubeServer")

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

	// Add the video file
	file, err := os.Open(videoFilePath)
	if err != nil {
		fmt.Printf("Error opening video file: %v\n", err)
		return err, "nil"
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return err, "nil"
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		return err, "nil"
	} else {
		fmt.Println("File data copied.")
	}

	// Close the multipart writer
	writer.Close()

	// Create a POST request with the multipart body
	request, err := http.NewRequest("POST", apiEndpoint, body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err, "nil"
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
		// return
	}
	fmt.Println("full response:", response)

	if response.StatusCode/100 != 2 {
		return err, "nil"
	}
	defer response.Body.Close()

	// Read and print the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body", err)
		return err, "nil"
	}
	fmt.Println("Response:", string(responseBody))

	var youtubeResponse config.YoutubeResponse
	err = json.Unmarshal(responseBody, &youtubeResponse)
	if err != nil {
		return err, "nil"
	}
	fmt.Println(youtubeResponse)
	fmt.Println("Vid = ", youtubeResponse.Id)

	return nil, youtubeResponse.Id

}

type VideoRequest struct {
	FileLocation string `json:"fileLocation"`
}

func uploadObjectOnGCS(videoPath string) (error, int) {

	//extract the token from service account
	tokenKey, err := config.GenerateJWTToken()
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
		return fmt.Errorf("File name doesnot seems to be of a media type\n"), http.StatusBadRequest
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

	//If admin approved then send the video on utube server as well

	return err, resp.StatusCode

}

func sendError(w http.ResponseWriter, errorMessage string, statusCode int) {
	http.Error(w, errorMessage, statusCode)
	//w.Write([]byte(errorMessage))
}

// tested and working - request will be from the user.
func UploadVideoOnStorageServer(w http.ResponseWriter, r *http.Request) {
	// _, err := r.Cookie("session_token")
	// fmt.Println(err)
	// if r.Method != http.MethodPost || AuthenticateUser(w, r) != http.StatusOK {
	// 	fmt.Println("Not status ok")
	// 	return
	// }

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
	uploadObjectOnGCS(destnationFileLocation)

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

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	fmt.Println("> Running UpdateTable")

	file, handler, err := r.FormFile("userFile")
	if err != nil {
		fmt.Println("not creating the formfile")
		sendError(w, "", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := isValidFile(handler); err != nil {
		fmt.Println("File is Invalid", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	originalFileName := strings.TrimSuffix(handler.Filename, fileExt)

	objectId := originalFileName + "_" + time.Now().Format("20060102_150405") + fileExt
	destnationFileLocation := config.STATIC_FILE_PATH + "//" + objectId
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

	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")
	audience := r.FormValue("audience")
	ageRestrictions := r.FormValue("ageRestrictions")
	tagInput := r.FormValue("tagInput")
	privacy := r.FormValue("privacy")

	tags := strings.Split(tagInput, ",")

	// println(tagInput)
	// return

	var nonAdminUserId string
	var AdminUserId string
	//Assume loggin user is non admin

	nonAdminUserId = config.USERID
	AdminUserId = config.UserRequest

	if config.LoginUserProfile == "Admin" {
		nonAdminUserId = config.UserRequest
		AdminUserId = config.USERID
	}
	db := config.DB
	q := "update metadataFiledTbl set Title = (@p1), Description = (@p2), Category = (@p3), Audience = (@p4), AgeRestrictions = (@p5), TagInput = (@p6), Privacy = (@p7) where AdminId = (@p8) AND userId = (@p9)"
	_, err = db.Exec(q, title, description, category, audience, ageRestrictions, tagInput, privacy, AdminUserId, nonAdminUserId)
	fmt.Println(err)

	//Temproary comment

	// err, statusCode := uploadObjectOnGCS(destnationFileLocation)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// if statusCode%100 != 2 {
	// 	fmt.Println("Error Code : ", statusCode)
	// }

	// We have the object Id of the video, Lets approve to the utube server
	if config.LoginUserProfile != "Admin" {
		return
	}
	//else approve the video to uplaod on the youtube
	err, videoId := uploadVideoOnYoutubeServer("snippet")
	if err != nil {
		fmt.Println(err)
		return
	}
	if videoId == "nil" {
		fmt.Println("VideoId not found")
		return
	}
	fmt.Println("Uploaded video id: ", videoId)
	// q = "update metadataFiledTbl set youtubeVideoUploadId = (@p1) where AdminId = (@p8) AND userId = (@p9)"
	// _, err = db.Exec(q, videoUploadId, AdminUserId, nonAdminUserId)
	// fmt.Println(err)

	//now we have been uploaded the video on youtube Sever successfully, we now go for update the video
	//1. update the db table
	//2. send utube request to the utube server for update the video metadata
	snippetParam := snippetBody{
		Title:       title,
		CategoryId:  "22",
		Description: description,
		Tags:        tags,
	}
	isMadeForKids := (audience == "Yes, it's 'Made for Kids'")

	statusParam := statusBody{
		PrivacyStatus:           privacy,
		SelfDeclaredMadeForKids: isMadeForKids,
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
	updateVideoMetadataOnYoutubeServer(string(jsonData), videoId)

}

// Tested code for update the video metada after the successfull upload of the video on the utube.
func FetchAndUploadVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetch&upload")
	w.Header().Set("content-type", "application/json")
	type formValue struct {
		Title           string   `json:"title"`
		Description     string   `json:"description"`
		Category        string   `json:"category"`
		Audience        string   `json:"audience"`
		AgeRestrictions string   `json:"ageRestrictions"`
		TagInput        []string `json:"tagInput"`
		Privacy         string   `json:"privacy"`
		UserId          string   `json:"coUserId"`
	}
	var (
		d       formValue
		videoId = "zHA-YusgSSo"
	)
	json.NewDecoder(r.Body).Decode(&d)
	userId := config.USERID
	query := `insert into metadataFiledTbl( AdminId, userId,Title,Description,Category,Audience,AgeRestrictions,TagInput,Privacy) values (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`
	db := config.DB

	var taginputLists string
	sz := len(d.TagInput)

	for i := 0; i < sz; i++ {
		taginputLists += d.TagInput[i]
		if i < sz-1 {
			taginputLists += ","
		}
	}
	_, err := db.Exec(query, userId, d.UserId, d.Title, d.Description, d.Category, d.Audience, d.AgeRestrictions, taginputLists, d.Privacy)
	if err != nil {
		fmt.Println(err)
	}
	//return //why i did this we should check if video dont upload on utube or gcs

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

	//fmt.Println(string(jsonData))
	updateVideoMetadataOnYoutubeServer(string(jsonData), videoId)
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
	uploadVideoOnYoutubeServer(data)
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

func insertdata(objectid string) {

	// Query data from the table
	db := config.DB
	insertSql := `insert into objectTbl (userName, objectId, bucketId) VALUES (@p1, @p2, @p3)`
	_, err := db.Exec(insertSql, config.USERID, objectid, "ytc-media-storage")
	if err != nil {
		fmt.Println(err)
	}
}
