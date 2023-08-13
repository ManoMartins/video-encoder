package services_test

import (
	"log"
	"testing"
	"time"

	"github.com/ManoMartins/video-encoder/application/repositories"
	"github.com/ManoMartins/video-encoder/application/services"
	"github.com/ManoMartins/video-encoder/domain"
	"github.com/ManoMartins/video-encoder/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {
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

	err = videoService.Finish()
	require.Nil(t, err)
}
