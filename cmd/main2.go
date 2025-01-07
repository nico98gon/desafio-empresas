package main

import (
	"desafio-empresas/internal/domain/email"
	"desafio-empresas/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main2() {
	start := time.Now()

	if err := utils.LoadEnvFile(".env"); err != nil {
		fmt.Printf("Error cargando archivo .env: %v\n", err)
		return
	}

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

	// Controlar concurrencia con semáforo
	const maxGoroutines = 100 // Máximo de goroutines activas
	sem := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup
	results := make(chan email.EmailData, 100) // Canal con buffer

	// Goroutine para recolectar e imprimir en terminal
	go func() {
		defer close(results)
		for email := range results {
			fmt.Printf("Email procesado: %+v\n", email)
		}
	}()

	// Procesar archivos en el directorio
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Adquirir permiso del semáforo
			sem <- struct{}{}
			wg.Add(1)

			go func(filePath string) {
				defer wg.Done()
				defer func() { <-sem }() // Liberar el semáforo

				// Procesar archivo
				email.ProcessFile(filePath, results, &wg)
			}(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Mostrar tiempo de ejecución
	fmt.Printf("Tiempo de ejecución: %v\n", time.Since(start))
}
