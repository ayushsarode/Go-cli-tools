package main

import (
	"net/http"
	"os"
	"encoding/json"
	"fmt"
	"io"
"os/exec"
)


type DictionaryResponse struct {
	Word     string `json:"word"`
	Meanings []struct {
		Definitions []struct {
			Definition string `json:"definition"`
		} `json:"definitions"`
	} `json:"meanings"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}


func getDefination(word string) {
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		body, _ := io.ReadAll(resp.Body)
		// Try parsing the error message from the API response
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			fmt.Println("Error parsing response:", err)
			return
		}
		fmt.Printf("Error: %s\n", errorResponse.Message)
		return
	}


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var dictionaryData []DictionaryResponse

	if err := json.Unmarshal(body, &dictionaryData); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if len(dictionaryData) > 0 {
		fmt.Printf("Definition of %s:\n", dictionaryData[0].Word)
		var definitionText string
		for _, meaning := range dictionaryData[0].Meanings {
			for _, definition := range meaning.Definitions {
				fmt.Println("-", definition.Definition)
				definitionText += definition.Definition + ". "
			}
		}

		speakText(dictionaryData[0].Word + ". " + definitionText)
	} else {
		fmt.Println("Word not found.")
	}
	
}



func speakText(text string) {
	cmd := exec.Command("flite", "-voice", "rms", "-t", text)
	err := cmd.Run()

	if err != nil {
		fmt.Println("Error speaking: ", err)
	}
}

func main(){
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run amin.go word")
	}

	word := os.Args[1]
	getDefination(word)
}