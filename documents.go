package primetrust

import (
	"fmt"
	"mime/multipart"
	"io/ioutil"
	"bytes"
	"net/http"
	"log"
	"encoding/json"
	"github.com/moul/http2curl"
	"os"
)

func UploadDocument(path string, accountId string, contactId string, description string, extension string, label string, mimeType string) (*map[string]interface{}, error) {
	apiUrl := fmt.Sprintf("%s/uploaded-documents", _apiPrefix)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	data := map[string]interface{}{
		"account-id":  accountId,
		"contact-id":  contactId,
		"description": description,
		"extension":   extension,
		"label":       label,
		"mime_type":   mimeType,
	}

		for key, val := range data {
		_ = writer.WriteField(key, val.(string))
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", _authHeader)

	client := &http.Client{}
	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bodyResp, _ := ioutil.ReadAll(res.Body)
	log.Println(res.StatusCode)
	if res.StatusCode != http.StatusOK {
		log.Println(bodyResp)
		return nil, err
	}

	var resData map[string]interface{}

	log.Printf("result  ", bodyResp)
	if err := json.Unmarshal(bodyResp, &resData); err != nil {
		return nil, err
	}
	return &resData, nil
}
