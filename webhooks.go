package primetrust

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BANKEX/go-primetrust/models"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateNewWebhook(webhook *models.Webhook) (*models.Webhook, error) {
	jsonData := new(bytes.Buffer)
	json.NewEncoder(jsonData).Encode(webhook)

	apiUrl := fmt.Sprintf("%s/webhook-configs", _apiPrefix)
	req, err := http.NewRequest("POST", apiUrl, jsonData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", _authHeader)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("%s: %s", res.Status, string(body)))
	}

	response := models.Webhook{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("unmarshal error")
	}

	return &response, nil
}

func UpdateWebhook(webhook *models.Webhook) (*models.Webhook, error) {
	jsonData := new(bytes.Buffer)
	json.NewEncoder(jsonData).Encode(webhook)

	apiUrl := fmt.Sprintf("%s/webhook-configs/%s", _apiPrefix, webhook.Data.ID)
	req, err := http.NewRequest("PATCH", apiUrl, jsonData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", _authHeader)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%s: %s", res.Status, string(body)))
	}

	response := models.Webhook{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("unmarshal error")
	}

	return &response, nil
}

func GetLastWebhook(accountId string) (*models.Webhook, error) {
	apiUrl := fmt.Sprintf("%s/webhook-configs?account.id=%s&one", _apiPrefix, accountId)
	req, err := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("Authorization", _authHeader)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	body, _ := ioutil.ReadAll(res.Body)

	response := models.Webhook{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("unmarshal error")
	}

	return &response, nil
}

func GetWebhook(webhookId string) (*models.Webhook, error) {
	apiUrl := fmt.Sprintf("%s/webhooks/%s", _apiPrefix, webhookId)
	req, err := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("Authorization", _authHeader)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	body, _ := ioutil.ReadAll(res.Body)

	response := models.Webhook{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("unmarshal error")
	}

	return &response, nil
}

func GetWebhookPayload(r *http.Request, secret string) (*models.WebhookPayload, error) {
	h := sha256.New()
	h.Write([]byte(secret))
	webhookHMAC := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if r.Header.Get("X-Prime-Trust-Webhook-Hmac") != webhookHMAC {
		log.Println("Wrong X-Prime-Trust-Webhook-Hmac:", r.Header.Get("X-Prime-Trust-Webhook-Hmac"))
		//return nil, errors.New("X-Prime-Trust-Webhook-Hmac header is absent or not valid")
	}

	body, _ := ioutil.ReadAll(r.Body)
	h2 := sha256.New()
	h2.Write(body)
	webhookHMAC2 := base64.StdEncoding.EncodeToString(h2.Sum(nil))

	log.Println("HMC1:", webhookHMAC)
	log.Println("HMC2:", webhookHMAC2)

	var webhookPayload models.WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&webhookPayload); err != nil {
		return nil, errors.New("error decoding webhook payload")
	}

	return &webhookPayload, nil
}
