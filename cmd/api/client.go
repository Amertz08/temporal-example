package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	baseURL := "http://localhost:8080"

	// POST /case
	newCase := Case{
		Name:      "adam",
		Address:   "1600 Pennselvania Ave.",
		Email:     "hello@example.com",
		VinNumber: "abcdef",
	}

	postBody, _ := json.Marshal(newCase)
	resp, err := http.Post(baseURL+"/case", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("POST /case response (%d):\n%s\n\n", resp.StatusCode, string(body))

	var caseResp CaseResponse
	json.Unmarshal(body, &caseResp)
	caseID := caseResp.Id

	// GET /case/:id
	resp, err = http.Get(baseURL + "/case/" + caseID)
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("GET /case/%s response (%d):\n%s\n\n", caseID, resp.StatusCode, string(body))

	// PUT /case/:id
	updateReq := UpdateRequest{
		Approved: true,
	}

	putBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest(http.MethodPut, baseURL+"/case/"+caseID, bytes.NewBuffer(putBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("PUT request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("PUT /case/%s response (%d):\n%s\n", caseID, resp.StatusCode, string(body))
}
