package main

import (
	"github.com/Kong/go-pdk/server"
	"keycloak-guard/cmd/plugin"
	"log"
)

func main() {
	log.Println("Starting plugin server...")
	//socketDir := "/tmp"
	//os.MkdirAll(socketDir, 0755)
	//socketPath := socketDir + "/keycloak-guard.socket"
	//log.Printf("Creating socket at %s", socketPath)
	//os.Remove(socketPath)
	err := server.StartServer(plugin.New, "0.1", 1000)
	if err != nil {
		log.Fatal("Failed to start plugin:", err)
	}
	log.Println("Plugin server started and listening on socket")
}
