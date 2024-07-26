package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ImageURL struct {
	URL string `json:"url"`
}

type Content struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Payload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Question struct {
	Id int `json:"id"`
	Question string            `json:"question"`
	Options  map[string]string `json:"options"`
	Answer   string            `json:"answer"`
	Hint string `json:"hint"`
}

type Choice struct {
	Index   int              `json:"index"`
	Message AssistantMessage `json:"message"`
}

type Completion struct {
	Choices []Choice `json:"choices"`
}

type AssistantMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// System and user messages for the AI model

// GenerateQuestions generates questions based on images extracted from a PDF
func GenerateQuestions(imagesDir string, questionsDir string, pdfPath string) error {
	apiKey := os.Getenv("OPENAI_API")
	if apiKey == ""{
		return fmt.Errorf("OPENAI_API key not set in \".env\" file")
	}

	questionAmount := os.Getenv("QUESTION_AMOUNT")
	if questionAmount == ""{
		questionAmount = "50"
	}

	amount, err := strconv.Atoi(questionAmount)
	if err != nil{
		return fmt.Errorf("amount set as QUESTION_AMOUNT is not correct")
	}

	systemMessage := fmt.Sprintf(`Generate %d questions about the given image(s), each with 4 options and an answer and an explanation as to why the answer is correct. Output the result as a JSON array without spaces or line breaks. also dont reference the image(s), the questions should be self explanitory. Use the format: [{"id":question number starting from zero, "question":"question text","options":{"A":"option1","B":"option2","C":"option3","D":"option4"},"answer":"option","hint":"explanation"}]`, amount)
	userMessage := fmt.Sprintf(`Generate %d questions from the given image(s).`, amount)
	var images []string

	// Walk through the directory to find and read image files
	err = filepath.Walk(filepath.Join(imagesDir, strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".jpg" {
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			images = append(images, base64.StdEncoding.EncodeToString(bytes))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Construct payload for the AI API
	payload := Payload{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role: "system",
				Content: []Content{
					{Type: "text", Text: systemMessage},
				},
			},
			{
				Role: "user",
				Content: []Content{
					{Type: "text", Text: userMessage},
				},
			},
		},
	}

	// Add image data to the payload
	for _, image := range images {
		payload.Messages[1].Content = append(payload.Messages[1].Content, Content{Type: "image_url", ImageURL: ImageURL{URL: fmt.Sprintf("data:image/jpeg;base64,%s", image)}})
	}

	// Marshal the payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create and send the API request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Read and check the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %q, body: %s", res.Status, body)
	}

	// Unmarshal the response to get the generated questions
	var completion Completion
	err = json.Unmarshal(body, &completion)
	if err != nil {
		return err
	}

	var questions []Question
	err = json.Unmarshal([]byte(completion.Choices[0].Message.Content), &questions)
	if err != nil {
		return err
	}

	// Format the questions to JSON and write to file
	formattedJSON, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(questionsDir, strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath)) + ".json"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(formattedJSON)
	if err != nil {
		return err
	}

	return nil
}