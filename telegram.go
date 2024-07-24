package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type TelegramBot struct {
	BotToken string
	ChatID   int
}

type SendMessageRequest struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func NewTelegramBot() (*TelegramBot, error) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatIDStr == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID is not set in environment variables")
	}

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TELEGRAM_CHAT_ID: %v", err)
	}

	return &TelegramBot{BotToken: botToken, ChatID: chatID}, nil
}

func (tb *TelegramBot) SendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tb.BotToken)

	payload := SendMessageRequest{
		ChatID: tb.ChatID,
		Text:   message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int `json:"message_id"`
		} `json:"result"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if !result.OK {
		return fmt.Errorf("failed to send message")
	}

	return nil
}
