package email

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	// "time"
)

// Función para procesar un archivo y extraer datos
func ProcessFile(filePath string, results chan<- EmailData, wg *sync.WaitGroup) error {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error abriendo el archivo %s: %v\n", filePath, err)
		return err
	}
	defer file.Close()

	var email EmailData
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Message-ID:") {
			email.MessageID = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		} else if strings.HasPrefix(line, "Date:") {
			email.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
			// dateStr := strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
			// parsedDate, parseErr := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", dateStr)
			// if parseErr != nil {
			// 	fmt.Printf("Error parsing date: %v\n", parseErr)
			// 	continue
			// }
			// email.Date = parsedDate
		} else if strings.HasPrefix(line, "From:") {
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		} else if strings.HasPrefix(line, "To:") {
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		} else if strings.HasPrefix(line, "Subject:") {
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error leyendo archivo %s: %v", filePath, err)
		return err
	}

	results <- email
	return nil
}
