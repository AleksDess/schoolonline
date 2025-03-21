package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"schoolonline/createtable"
	"schoolonline/crypto"
	"schoolonline/launch"
	"schoolonline/metrics"
	"schoolonline/postgree"
	"schoolonline/routes"
)

func main() {

	fmt.Println("Hello")
	fmt.Println(generateSecretKey(32))

	fmt.Println("pass:", crypto.GeneratePassword())

	// GoTo
	GoToInit()
	routes.GoToPathHTML()
	routes.GoToRunRoutes()
	createtable.GoToCreateTable()

	// os.Exit(1)

	// запуск БД
	postgree.RunDB()
	defer postgree.MainDB.Close()
	defer postgree.MainDBX.Close()

	createtable.CreateTable()
	createtable.CreateManager()

	postgree.Set_idle_in_transaction_session_timeout()

	if launch.Launch == "server" {
		go metrics.RunMetrics()
	}

	routes.RunRoutes()
}

// Функция для генерации случайного секретного ключа
func generateSecretKey(length int) string {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(key)
}
