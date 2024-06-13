package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"myapp/internal-xkcd/core/domain"
	"net/http"
	"os"
)

type GptAPI struct {
	Client    *http.Client
	APIKey    string
	APIURL    string
	CatalogID string
	FolderID  string
}

func NewGptAPI() *GptAPI {
	return &GptAPI{
		APIURL:    os.Getenv("YANDEX_GPT_API_URL"),
		Client:    &http.Client{},
		APIKey:    os.Getenv("YANDEX_GPT_API_KEY"),
		CatalogID: os.Getenv("YANDEX_GPT_CATALOG_ID"),
		FolderID:  os.Getenv("YANDEX_GPT_FOLDER_ID"),
	}
}

func (g *GptAPI) GetComicsDescription(comic domain.Comics) (string, error) {

	comic.Description = ""

	type prompt struct {
		ModelUri          string `json:"modelUri"`
		CompletionOptions struct {
			Stream      bool    `json:"stream"`
			Temperature float32 `json:"temperature"`
			MaxTokens   string  `json:"maxTokens"`
		} `json:"completionOptions"`
		Messages []struct {
			Role string `json:"role"`
			Text string `json:"text"`
		} `json:"messages"`
	}

	comicBytes, mErr := json.Marshal(comic)
	if mErr != nil {
		return "", mErr
	}
	comicDataString := string(comicBytes)

	//YandexGPT Pro	3	gpt://<идентификатор_каталога>/yandexgpt/latest	Асинхронный, синхронный
	//YandexGPT Lite	2	gpt://<идентификатор_каталога>/yandexgpt-lite/latest	Асинхронный, синхронный
	//YandexGPT Lite RC	3	gpt://<идентификатор_каталога>/yandexgpt-lite/rc	Асинхронный, синхронный
	//Краткий пересказ	2	gpt://<идентификатор_каталога>/summarization/latest	Асинхронный, синхронный
	//Модель, дообученная в Yandex DataSphere	3	ds://<идентификатор_дообученной_модели>

	newPrompt := prompt{
		ModelUri: fmt.Sprintf("gpt://%s/yandexgpt/latest", g.CatalogID),
		CompletionOptions: struct {
			Stream      bool    `json:"stream"`
			Temperature float32 `json:"temperature"`
			MaxTokens   string  `json:"maxTokens"`
		}{
			Stream:      false,
			Temperature: 0.7,
			MaxTokens:   "400",
		},
		Messages: []struct {
			Role string `json:"role"`
			Text string `json:"text"`
		}{
			{Role: "system", Text: "Ты — бот, который помогает пользователям понять шутки в комиксах, представленные в текстовом формате. Ваша задача — кратко описать происходящее в комиксе, а затем объяснить, в чём заключается шутка"},
			{Role: "user", Text: fmt.Sprintf("Можешь ли ты объяснить мне, что происходит в комиксе с этим описанием? %s", comicDataString)},
		},
	}

	jsonPrompt, mpErr := json.Marshal(newPrompt)
	if mpErr != nil {
		return "", mpErr
	}

	req, rErr := http.NewRequest("POST", g.APIURL, bytes.NewBuffer(jsonPrompt))
	if rErr != nil {
		return "", rErr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", g.APIKey))
	req.Header.Set("x-folder-id", g.FolderID)

	resp, rErr := g.Client.Do(req)
	if rErr != nil {
		return "", rErr
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Debug("unhandled error", "error", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return "", errors.New("bad response")
	}

	type Response struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Role string `json:"role"`
					Text string `json:"text"`
				} `json:"message"`
				Status string `json:"status"`
			} `json:"alternatives"`
			Usage struct {
				InputTextTokens  string `json:"inputTextTokens"`
				CompletionTokens string `json:"completionTokens"`
				TotalTokens      string `json:"totalTokens"`
			} `json:"usage"`
			ModelVersion string `json:"modelVersion"`
		} `json:"result"`
	}

	var response Response
	if dErr := json.NewDecoder(resp.Body).Decode(&response); dErr != nil {
		return "", dErr
	}

	if len(response.Result.Alternatives) != 1 {
		return "", errors.New("response empty error")

	}
	return response.Result.Alternatives[0].Message.Text, nil
}
