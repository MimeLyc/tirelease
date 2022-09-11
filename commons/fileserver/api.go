package fileserver

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadFile(fromFilePath, toFilePath string) (string, error) {
	uploadUrl := fmt.Sprintf("%s/%s", FsUrl, UploadPath)
	fileBytes, multiFormWriter, err := constructFormData(fromFilePath, toFilePath)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", uploadUrl, fileBytes)
	if err != nil {
		return "", err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", multiFormWriter.FormDataContentType())

	client := &http.Client{}
	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status uploading file to fileserver: %s", res.Status)
		return "", err
	}

	downloadPath := fmt.Sprintf(DownloadPath, toFilePath)
	return fmt.Sprintf("%s/%s", FsUrl, downloadPath), nil
}

func constructFormData(fromFilePath, toFilePath string) (*bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	multiWriter := multipart.NewWriter(&b)
	defer multiWriter.Close()

	file, err := os.Open(fromFilePath)
	defer file.Close()
	if err != nil {
		return nil, nil, err
	}
	fileReader := io.Reader(file)

	var fw io.Writer
	if x, ok := fileReader.(io.Closer); ok {
		defer x.Close()
	}

	if fw, err = multiWriter.CreateFormFile(toFilePath, fromFilePath); err != nil {
		return nil, nil, err
	}

	if _, err = io.Copy(fw, fileReader); err != nil {
		return nil, nil, err
	}

	return &b, multiWriter, nil
}
