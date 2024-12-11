package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

type RequestBody struct {
	OriginalTarget string `json:"originalTarget"`
	ArtifactName   string `json:"artifactName"`
	CallbackIri    string `json:"callbackIri"`
}

type StoredRequest struct {
	ID     int
	Body   string
	Header http.Header
}

var forwardingEnabled = true
var mu sync.Mutex

var DOCKER_LOCATION = "http://yggdrasil:8080"

// var ORIGINAL_CALLBACK = "http://localhost:8083/notifications"
var ORIGINAL_CALLBACK = "http://host.docker.internal:8083/notifications"

var ALTERED_CALLBACK = "http://websub_proxy:3000/altered-callback"

var storedRequests []StoredRequest
var requestIDCounter = 0

// unlockRoomHandler executes the specified POST request
func unlockRoomHandler(c *fiber.Ctx) error {
	// Define the request URL
	url := "http://host.docker.internal:8080/workspaces/103/artifacts/r3/unlockRoom"

	// Create a new POST request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create request")
	}

	// Add required headers
	req.Header.Set("X-Agent-WebID", "http://localhost:8080/agents/valentin")

	// Execute the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute request: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to execute request")
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}

	// Return the response body to the client
	return c.Status(resp.StatusCode).SendString(string(responseBody))
}

func main() {
	// Initialize the template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Configure CORS to allow the X-Agent-WebID header
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this to specific origins if needed
		AllowHeaders: "Origin, Content-Type, Accept, X-Agent-WebID",
	}))

	// Route to toggle forwarding on or off
	app.Post("/toggle", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()

		forwardingEnabled = !forwardingEnabled
		status := "disabled"
		if forwardingEnabled {
			status = "enabled"
		}

		log.Printf("Request forwarding %s\n", status)
		return c.SendString(fmt.Sprintf("Request forwarding %s", status))
	})

	// Route to handle the focus requests and modify the callback URI
	app.Post("/focus", func(c *fiber.Ctx) error {
		log.Printf("Received POST request with body: %s\n", c.Body())

		mu.Lock()
		if !forwardingEnabled {
			mu.Unlock()
			log.Println("Forwarding is disabled, responding without forwarding.")
			return c.SendString("Forwarding is currently disabled")
		}
		mu.Unlock()

		var body RequestBody
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			log.Printf("Failed to parse JSON body: %v\n", err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
		}

		originalTarget := body.OriginalTarget
		if strings.HasPrefix(originalTarget, "http://localhost:8080") {
			originalTarget = strings.Replace(originalTarget, "http://localhost:8080", DOCKER_LOCATION, 1)
			log.Printf("Modified OriginalTarget to: %s\n", originalTarget)
		}

		modifiedBody := map[string]string{
			"artifactName": body.ArtifactName,
			"callbackIri":  ALTERED_CALLBACK,
		}

		modifiedBodyJson, err := json.Marshal(modifiedBody)
		if err != nil {
			log.Printf("Failed to marshal modified body to JSON: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create JSON")
		}

		log.Printf("Modified request body: %s\n", modifiedBodyJson)

		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("POST", originalTarget, bytes.NewBuffer(modifiedBodyJson))
		if err != nil {
			log.Printf("Failed to create HTTP request: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create request")
		}

		req.Header = http.Header{}
		for k, v := range c.GetReqHeaders() {
			req.Header.Set(k, v[0])
		}

		log.Printf("Forwarding headers: %v\n", req.Header)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to forward request: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to forward request")
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response from original target: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response")
		}

		log.Printf("Received response with status: %d and body: %s\n", resp.StatusCode, respBody)

		for k, v := range resp.Header {
			c.Set(k, v[0])
		}
		return c.Status(resp.StatusCode).Send(respBody)
	})

	// Route to handle the redirected messages and store them
	app.Post("/altered-callback", func(c *fiber.Ctx) error {
		log.Printf("Received message for altered callback with body: %s\n", c.Body())

		body := string(c.Body())

		if strings.Contains(body, "@context") {
			log.Println("JSON-LD message received, storing for later action")
			mu.Lock()
			requestIDCounter++
			storedRequests = append(storedRequests, StoredRequest{
				ID:     requestIDCounter,
				Body:   body,
				Header: c.GetReqHeaders(),
			})
			mu.Unlock()

			return c.SendString("JSON-LD message stored and waiting for action")
		} else {
			log.Println("Forwarding non-JSON-LD message to original callback")

			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("POST", ORIGINAL_CALLBACK, bytes.NewBuffer(c.Body()))
			if err != nil {
				log.Printf("Failed to create HTTP request: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to create request")
			}

			req.Header = http.Header{}
			for k, v := range c.GetReqHeaders() {
				req.Header.Set(k, v[0])
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Failed to forward message to original callback: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to forward message")
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response from original callback: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response")
			}

			log.Printf("Received response from original callback with status: %d and body: %s\n", resp.StatusCode, respBody)

			for k, v := range resp.Header {
				c.Set(k, v[0])
			}
			return c.Status(resp.StatusCode).Send(respBody)
		}
	})

	// Route to display stored requests with "forward" options
	app.Get("/requests", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()

		// Render the stored requests using the template
		return c.Render("requests", fiber.Map{
			"Requests": storedRequests,
		})
	})

	app.Get("/actions", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()

		// Render the stored actions using the template
		return c.Render("actions", fiber.Map{})
	})

	// Route to forward a stored request to the original callback
	app.Post("/forward", func(c *fiber.Ctx) error {
		payload := struct {
			ID string `json:"id"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		id := payload.ID
		log.Printf("Forwarding request with ID: %s\n", id)

		var reqToForward *StoredRequest
		mu.Lock()
		for i, req := range storedRequests {
			if strconv.Itoa(req.ID) == id {
				reqToForward = &storedRequests[i]
				storedRequests = append(storedRequests[:i], storedRequests[i+1:]...)
				break
			}
		}
		mu.Unlock()

		if reqToForward == nil {
			return c.Status(fiber.StatusNotFound).SendString("Request not found")
		}

		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("POST", ORIGINAL_CALLBACK, bytes.NewBuffer([]byte(reqToForward.Body)))
		if err != nil {
			log.Printf("Failed to create HTTP request: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create request")
		}

		req.Header = reqToForward.Header

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to forward message to original callback: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to forward message")
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response from original callback: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response")
		}

		log.Printf("Received response from original callback with status: %d and body: %s\n", resp.StatusCode, respBody)

		for k, v := range resp.Header {
			c.Set(k, v[0])
		}
		return c.Status(resp.StatusCode).Send(respBody)
	})

	// Register the unlock room handler
	app.Post("/unlock-room", unlockRoomHandler)

	log.Println("Server listening on port 3000")
	log.Fatal(app.Listen(":3000"))
}
