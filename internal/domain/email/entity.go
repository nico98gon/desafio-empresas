package email

// Estructura para almacenar datos extraídos de un archivo
type EmailData struct {
	MessageID string
	Date      string
	// Date      time.Time
	From      string
	To        string
	Subject   string
}