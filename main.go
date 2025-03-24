package main

import (
	"fmt"
	discovery "github.com/Thomas-PEYROT/discovery-client"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Enregistrement du microservice
	discovery.RegisterMicroservice("service-a")

	// Création du routeur Gin
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello microservice A !")
	})

	// Démarrage du serveur dans une goroutine
	port := discovery.ServiceInformations.Port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur lors du démarrage du serveur : %v", err)
		}
	}()

	// Capture des signaux système pour exécuter le code avant la fermeture
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Attente du signal d'arrêt
	<-quit
	log.Println("Arrêt du microservice...")

	// Code de nettoyage avant l'arrêt
	discovery.UnregisterMicroservice() // Exemple d'une éventuelle fonction de nettoyage
	log.Println("Microservice arrêté proprement.")
}
