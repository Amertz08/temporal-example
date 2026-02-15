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

	// PATCH /case/:id - update approved field
	approved := true
	patchReq1 := map[string]interface{}{
		"approved": approved,
	}

	patchBody1, _ := json.Marshal(patchReq1)
	req, _ := http.NewRequest(http.MethodPatch, baseURL+"/case/"+caseID, bytes.NewBuffer(patchBody1))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("PATCH request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("PATCH /case/%s (approved) response (%d):\n%s\n\n", caseID, resp.StatusCode, string(body))

	// PATCH /case/:id - update manufactured field
	manufactured := true
	patchReq2 := map[string]interface{}{
		"manufactured": manufactured,
	}

	patchBody2, _ := json.Marshal(patchReq2)
	req, _ = http.NewRequest(http.MethodPatch, baseURL+"/case/"+caseID, bytes.NewBuffer(patchBody2))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("PATCH request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("PATCH /case/%s (manufactured) response (%d):\n%s\n", caseID, resp.StatusCode, string(body))
}
