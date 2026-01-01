package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// Configurations for messengers
type Config struct {
	SlackWebhookURL   string
	TelegramBotToken  string
	TelegramChatID    string
	DiscordWebhookURL string
}

var cfg Config

func main() {
	rootCmd := &cobra.Command{
		Use:   "messenger",
		Short: "CLI Messenger Center to send messages to Slack, Telegram, and Discord",
	}

	rootCmd.AddCommand(discordCmd())
	rootCmd.AddCommand(slackCmd())
	rootCmd.AddCommand(telegramCmd())

	// Set configurations directly in code
	cfg = Config{
		SlackWebhookURL:   "https://hooks.slack.com/services/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",
		TelegramBotToken:  "8344022264:AAEWEMvPCTYC0Oz51O5NlJgNyJvAWVU9fDE",
		TelegramChatID:    "725269091",
		DiscordWebhookURL: "https://discord.com/api/webhooks/1454602976101404813/oF12i3ixHY2N59wWdBfwGXlhyLb0sR-RdEspsJu-XKExoxzwWGIrHEnvgYjHB5mWo5SV",
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}

// Command for Discord
func discordCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "discord [message]",
		Short: "Send a message to Discord",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := args[0]
			if cfg.DiscordWebhookURL == "" {
				fmt.Println("Error: DISCORD_WEBHOOK_URL is not set.")
				return
			}
			err := sendDiscordMessage(cfg.DiscordWebhookURL, message)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Message sent to Discord!")
			}
		},
	}
}

// Command for Slack
func slackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "slack [message]",
		Short: "Send a message to Slack",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := args[0]
			if cfg.SlackWebhookURL == "" {
				fmt.Println("Error: SLACK_WEBHOOK_URL is not set.")
				return
			}
			err := sendSlackMessage(cfg.SlackWebhookURL, message)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Message sent to Slack!")
			}
		},
	}
}

// Command for Telegram
func telegramCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "telegram [message]",
		Short: "Send a message to Telegram",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := args[0]
			if cfg.TelegramBotToken == "" || cfg.TelegramChatID == "" {
				fmt.Println("Error: TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID is not set.")
				return
			}
			err := sendTelegramMessage(cfg.TelegramBotToken, cfg.TelegramChatID, message)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Message sent to Telegram!")
			}
		},
	}
}

// Sending message to Discord
func sendDiscordMessage(webhookURL, message string) error {
	payload := map[string]string{"content": message}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message to Discord, status code: %d", resp.StatusCode)
	}

	return nil
}

// Sending message to Slack
func sendSlackMessage(webhookURL, message string) error {
	payload := map[string]string{"text": message}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Slack, status code: %d", resp.StatusCode)
	}

	return nil
}

// Sending message to Telegram
func sendTelegramMessage(botToken, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Telegram, status code: %d", resp.StatusCode)
	}

	return nil
}
