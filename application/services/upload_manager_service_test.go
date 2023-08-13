package services_test

import (
	"log"
	"os"
	"testing"

	"github.com/ManoMartins/video-encoder/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()
	videoService := services.NewVideoService()

	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("path")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUploadService := services.NewVideoUploadService()

	videoUploadService.OutputBucket = "video-encoder-bucket"
	videoUploadService.VideoPath = os.Getenv("STORAGE_PATH") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUploadService.ProcessUpload(50, doneUpload)

	result := <-doneUpload

	require.Equal(t, result, "upload completed")
}
