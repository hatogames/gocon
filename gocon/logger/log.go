package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendDiscord(msg string) error {
	webhookURL := "https://discord.com/api/webhooks/1424467607154720953/0PESG1ICvshDMriZgFL7NH1iFBQ8KTlfcc554lo77Sh0l_3UAv-Bn8medb6vd-MzEiSs"

	payload := map[string]string{
		"content": msg, // Nachricht
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("Discord Webhook returned status %d", resp.StatusCode)
	}

	return nil
}
