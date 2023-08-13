package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/ManoMartins/video-encoder/application/repositories"
	"github.com/ManoMartins/video-encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)
	r, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("STORAGE_PATH") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		return err
	}

	_, err = f.Write(body)

	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("video %s downloaded", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("STORAGE_PATH")+"/"+v.Video.ID, os.ModePerm)

	if err != nil {
		return err
	}

	source := os.Getenv("STORAGE_PATH") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("STORAGE_PATH") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}

	cmdArgs = append(cmdArgs, os.Getenv("STORAGE_PATH")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("STORAGE_PATH")+"/"+v.Video.ID+".mp4")
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(os.Getenv("STORAGE_PATH") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		log.Println("error removing mp4", err)
		return err
	}

	err = os.Remove(os.Getenv("STORAGE_PATH") + "/" + v.Video.ID + ".frag")

	if err != nil {
		log.Println("error removing frag", err)
		return err
	}

	err = os.RemoveAll(os.Getenv("STORAGE_PATH") + "/" + v.Video.ID)

	if err != nil {
		log.Println("error removing all file", err)
		return err
	}

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s", string(out))
	}
}
