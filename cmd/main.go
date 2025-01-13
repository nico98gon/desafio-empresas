package main

import (
	"desafio-empresas/internal/domain/email"
	"desafio-empresas/internal/infrastructure/db"
	"desafio-empresas/internal/infrastructure/zinc"
	"desafio-empresas/internal/utils"

	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)


func withSemaphore(sem chan struct{}, f func()) {
	sem <- struct{}{}  // Adquirir el semáforo
	defer func() { <-sem }() // Liberar el semáforo
	f()
}

func main() {
	start := time.Now()

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

	// repo := email.NewEmailRepository(dbConn)

	dir := os.Getenv("EMAILS_DIR")
	if dir == "" {
		fmt.Println("EMAILS_DIR no está configurado")
		return
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("El directorio %s no existe\n", dir)
		return
	}

	const maxGoroutines = 100
	sem := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup
	results := make(chan email.EmailData, 100)

	// Goroutine para recolectar y guardar en MySQL y Zincsearch
	go func() {
		defer close(results)

		for emailData := range results {
			// if err := repo.Save(emailData); err != nil {
			// 	fmt.Printf("Error guardando email en MySQL: %v\n", err)
			// 	continue
			// }

			if err := zinc.IndexToZinc(emailData); err != nil {
				fmt.Printf("Error indexando email en Zincsearch: %v\n", err)
			}
		}
	}()

	// Recorrer los archivos del directorio
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			wg.Add(1)
			withSemaphore(sem, func() {
				defer wg.Done() // Asegurar que el contador de goroutines se decremente

				err := email.ProcessFile(path, results, &wg)
				if err != nil {
					fmt.Printf("Error procesando archivo %s: %v\n", path, err)
				}
			})
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error recorriendo el directorio: %v\n", err)
		return
	}

	wg.Wait()

	fmt.Printf("Tiempo de ejecución: %v\n", time.Since(start))
}
