package blockip

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Структура для хранения данных о запросах по IP
type IPTracker struct {
	Count    int       // Количество некорректных запросов
	LastSeen time.Time // Время последнего запроса
}

// Глобальная карта для хранения данных о каждом IP
var ipData = make(map[string]*IPTracker)
var mu sync.Mutex // Мьютекс для синхронизации

// Подозрительные подстроки в URL, которые точно не относятся к Go-серверу
var suspiciousPatterns = []string{
	".php", ".asp", ".jsp", ".exe", ".sh", ".env", "wp-", "admin", "login",
}

// Middleware для отслеживания некорректных запросов
func RequestTrackerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path

		// Проверка, существует ли IP в карте
		mu.Lock()
		tracker, exists := ipData[ip]
		if !exists {
			tracker = &IPTracker{Count: 0, LastSeen: time.Now()}
			ipData[ip] = tracker
		}
		mu.Unlock()

		// Проверка на подозрительные запросы
		if isSuspiciousRequest(path) {
			fmt.Printf("--- запрос подозрителен  ---  IP %s/n", ip)
			mu.Lock()
			tracker.Count++
			tracker.LastSeen = time.Now()

			// Если больше 10 подозрительных запросов, блокируем IP
			if tracker.Count >= 10 {
				fmt.Printf("--- БЛОКИРОВКА ПОДОЗРИТЕЛЬНОГО IP  ---  IP %s/n", ip)
				// blockIP(ip)
				// fmt.Printf("IP %s заблокирован за 10 подозрительных запросов\n", ip)
			}
			mu.Unlock()

			// Возвращаем 403 Forbidden и не обрабатываем запрос дальше
			c.JSON(http.StatusForbidden, gin.H{"message": "Suspicious activity detected"})
			c.Abort()
			return
		}

		// Если запрос валидный, обнуляем счётчик некорректных запросов
		c.Next()

		// Если страница не найдена (404), увеличиваем счётчик
		if c.Writer.Status() == http.StatusNotFound {
			mu.Lock()
			tracker.Count++
			tracker.LastSeen = time.Now()

			// Если больше 10 некорректных запросов, блокируем IP
			if tracker.Count >= 10 {
				blockIP(ip)
				fmt.Printf("IP %s заблокирован за 10 некорректных запросов\n", ip)
			}
			mu.Unlock()
		} else {
			// Если запрос успешный, обнуляем счётчик
			mu.Lock()
			tracker.Count = 0
			mu.Unlock()
		}
	}
}

// Функция для блокировки IP
func blockIP(ip string) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`New-NetFirewallRule -DisplayName "Block IP %s" -Direction Inbound -RemoteAddress %s -Action Block`, ip, ip))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Ошибка при блокировке IP: %s\n", err)
	}
}

// Проверка, является ли запрос подозрительным
func isSuspiciousRequest(path string) bool {
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(strings.ToLower(path), pattern) {
			return true
		}
	}
	return false
}
