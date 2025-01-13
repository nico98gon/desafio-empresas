package zinc

import (
	"bytes"
	"desafio-empresas/internal/domain/email"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func IndexToZinc(emailData email.EmailData) error {
	url := os.Getenv("ZINC_URL")
	jsonData, _ := json.Marshal(emailData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	username := os.Getenv("ZINC_USERNAME")
	password := os.Getenv("ZINC_PASSWORD")
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error indexando a Zincsearch: %s", resp.Status)
	}
	return nil
}