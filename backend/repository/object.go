package repository

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type object struct {
	width   int
	format  string
	quality int
}

const (
	key1 = "thumbnail"
	key2 = "image"
)

type ObjectRepository interface {
	FetchLocalThumbnail(baseFilepath string) (*[]byte, error)
	FetchCDNThumbnail(baseFilepath, publishDomain string) (*[]byte, error)
	FetchLocalImage(baseFilepath string) (*[]byte, error)
	FetchCDNImage(baseFilepath, publishDomain string) (*[]byte, error)
}

func NewObjectRepository() ObjectRepository {
	return &object{}
}

// baseFilepath = R2上のpath
func (i *object) createFilename(cacheKey, baseFilepath string) string {
	cacheDir := fmt.Sprintf("cache/%s", cacheKey)
	filename := strings.ReplaceAll(baseFilepath, "/", "_")
	thumbnailFilename := path.Join(cacheDir, filename)
	return thumbnailFilename
}

func (i *object) FetchLocalThumbnail(baseFilepath string) (*[]byte, error) {
	thumbnailFilename := i.createFilename(key1, baseFilepath)

	file, err := os.ReadFile(thumbnailFilename)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (i *object) FetchCDNThumbnail(baseFilepath, publishDomain string) (*[]byte, error) {
	thumbnailFilename := i.createFilename(key1, baseFilepath)
	publishURL := "https://" + path.Join(publishDomain, baseFilepath)
	imageFormat := "w=" + strconv.Itoa(i.width) + ",f=" + i.format + ",q=" + strconv.Itoa(i.quality)

	cdnURL := "https://www.maretol.xyz/cdn-cgi/image/#{option}/#{origin}"

	fetchURL := strings.Replace(cdnURL, "#{option}", imageFormat, 1)
	fetchURL = strings.Replace(fetchURL, "#{origin}", publishURL, 1)

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", fetchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status is not ok: %s", resp.Status)
	}
	file := make([]byte, 0)
	_, err = resp.Body.Read(file)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(thumbnailFilename, file, 0644)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (i *object) FetchLocalImage(baseFilepath string) (*[]byte, error) {
	imageFilename := i.createFilename(key2, baseFilepath)

	file, err := os.ReadFile(imageFilename)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (i *object) FetchCDNImage(baseFilepath, publishDomain string) (*[]byte, error) {
	imageFilename := i.createFilename(key2, baseFilepath)
	publishURL := "https://" + path.Join(publishDomain, baseFilepath)
	imageFormat := "w=1200,f=webp,q=80"

	cdnURL := "https://www.maretol.xyz/cdn-cgi/image/#{option}/#{origin}"

	fetchURL := strings.Replace(cdnURL, "#{option}", imageFormat, 1)
	fetchURL = strings.Replace(fetchURL, "#{origin}", publishURL, 1)

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", fetchURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("https status is not ok: %s", resp.Status)
	}
	file := make([]byte, 0)
	_, err = resp.Body.Read(file)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(imageFilename, file, 0644)
	if err != nil {
		return nil, err
	}
	return &file, nil
}
