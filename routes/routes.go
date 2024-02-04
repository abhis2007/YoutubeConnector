package routes

import (
	"github.com/abhis2007/YOUTUECONNECTOR/controller"
	"github.com/gorilla/mux"
)

func RouterConfiguration(router *mux.Router) {
	router.HandleFunc("/", controller.Index).Methods("GET")
	router.HandleFunc("/signUp", controller.SignUp).Methods("GET")
	router.HandleFunc("/Login", controller.Login).Methods("GET")
	router.HandleFunc("/SignupData", controller.SignupData).Methods("POST")
	router.HandleFunc("/LoginData", controller.LoginData).Methods("POST")
	router.HandleFunc("/videos", controller.Videos).Methods("GET")
	router.HandleFunc("/upload", controller.UploadVideo).Methods("GET")
	router.HandleFunc("/fetchAndUploadVideo", controller.FetchAndUploadVideo).Methods("POST")
	router.HandleFunc("/UploadVideoOnStorageServer", controller.UploadVideoOnStorageServer).Methods("POST")

	// router.HandleFunc("/getVideoDetails", controller.GetVideoDetails).Methods("GET")
	// router.HandleFunc("/uploadVideo", controller.UploadVideo).Methods("POST")
	// router.HandleFunc("/updateVideo", controller.UpdateVideo).Methods("PUT")
	// router.HandleFunc("/addVideoThumbnail", controller.AddThumbnail).Methods("POST")
}
