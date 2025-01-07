package main

import (
	"desafio-empresas/internal/domain/email"
	"desafio-empresas/internal/infrastructure/db"
	"desafio-empresas/internal/utils"

	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	if err := utils.LoadEnvFile(".env"); err != nil {
		fmt.Printf("Error cargando archivo .env: %v\n", err)
		return
	}

	dbConn, err := db.InitMySQL()
	if err != nil {
		fmt.Printf("Error inicializando base de datos: %v\n", err)
		return
	}
	defer dbConn.Close()

	repo := email.NewEmailRepository(dbConn)

	dir := os.Getenv("EMAILS_DIR")
	if dir == "" {
		fmt.Println("EMAILS_DIR no está configurado")
		return
	}

	// Verificar si el directorio existe
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("El directorio %s no existe\n", dir)
		return
	}

	var wg sync.WaitGroup
	results := make(chan email.EmailData, 100) // Canal con buffer

	// Goroutine para recolectar e imprimir en terminal
	go func() {
		for email := range results {
			fmt.Printf("Processed Email: %+v\n", email)
		}
	}()

	// Goroutine para recolectar y guardar en mysql
	go func() {
		for emailData := range results {
			if err := repo.Save(emailData); err != nil {
				fmt.Printf("Error guardando email: %v\n", err)
			}
		}
	}()

	// Recorrer los archivos del directorio
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Verificar que sea un archivo regular
		if !info.IsDir() {
			wg.Add(1)
			go email.ProcessFile(path, results, &wg)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()
	close(results)
}