package metrics

import (
	"fmt"
	"runtime"
	"schoolonline/postgree"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Структура для хранения метрик системы
type SystemMetrics struct {
	MemoryAlloc      uint64    `db:"memory_alloc"`       // Используемая память программой
	MemoryTotalAlloc uint64    `db:"memory_total_alloc"` // Всего выделено памяти программой
	MemorySys        uint64    `db:"memory_sys"`         // Память, выделенная системой
	NumGoroutines    int       `db:"num_goroutines"`     // Количество горутин
	CPULoad          float64   `db:"cpu_load"`           // Загрузка процессора в процентах
	Uptime           time.Time `db:"uptime"`
	NumCPU           int       `db:"num_cpu"`
	OS               string    `db:"os"`
	Arch             string    `db:"arch"`
	TotalMemory      uint64    `db:"total_memory"`     // Вся память системы
	MemoryUsagePct   float64   `db:"memory_usage_pct"` // Процент использования памяти
}

// Функция для сбора метрик
func getMetrics() (SystemMetrics, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Получение информации о системной памяти
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
		return SystemMetrics{}, fmt.Errorf("ошибка при получении системной памяти: %w", err)
	}

	// Получение загрузки CPU за последний интервал времени (1 секунда)
	cpuLoad, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		fmt.Println(err)
		return SystemMetrics{}, fmt.Errorf("ошибка при получении загрузки CPU: %w", err)
	}

	// Расчет процента использования памяти программой
	memoryUsagePct := (float64(memStats.Alloc) / float64(vmStat.Total)) * 100

	return SystemMetrics{
		MemoryAlloc:      memStats.Alloc,
		MemoryTotalAlloc: memStats.TotalAlloc,
		MemorySys:        memStats.Sys,
		NumGoroutines:    runtime.NumGoroutine(),
		CPULoad:          cpuLoad[0], // Используем первое значение, так как оно общее по всей системе
		Uptime:           time.Now(),
		NumCPU:           runtime.NumCPU(),
		OS:               runtime.GOOS,
		Arch:             runtime.GOARCH,
		TotalMemory:      vmStat.Total,
		MemoryUsagePct:   memoryUsagePct,
	}, nil
}

// Функция для красивой печати метрик
func (sm SystemMetrics) Print() {
	fmt.Println("System Metrics:")
	fmt.Println("-------------------------------")
	fmt.Printf("Memory Allocated:     %d bytes\n", sm.MemoryAlloc)
	fmt.Printf("Total Memory Alloc:   %d bytes\n", sm.MemoryTotalAlloc)
	fmt.Printf("System Memory:        %d bytes\n", sm.MemorySys)
	fmt.Printf("Total System Memory:  %d bytes\n", sm.TotalMemory)
	fmt.Printf("Memory Usage:         %.2f%%\n", sm.MemoryUsagePct)
	fmt.Printf("CPU Load:             %.2f%%\n", sm.CPULoad)
	fmt.Printf("Number of Goroutines: %d\n", sm.NumGoroutines)
	fmt.Printf("Number of CPUs:       %d\n", sm.NumCPU)
	fmt.Printf("Operating System:     %s\n", sm.OS)
	fmt.Printf("Architecture:         %s\n", sm.Arch)
	fmt.Printf("Uptime:               %s\n", sm.Uptime.Format("2006-01-02 15:04:05"))
	fmt.Println("-------------------------------")
}

func CreateMetricsTableIfNotExists() error {
	const query = `
		CREATE TABLE system_metrics (
			id SERIAL PRIMARY KEY,
			memory_alloc BIGINT,
			memory_total_alloc BIGINT,
			memory_sys BIGINT,
			total_memory BIGINT,
			memory_usage_pct REAL,
			num_goroutines INT,
			cpu_load REAL,
			num_cpu INT,
			os VARCHAR(50),
			arch VARCHAR(50),
			uptime TIMESTAMP
			);
			`

	_, err := postgree.MainDBX.Exec(query)
	return err
}

func RunMetrics() {
	for {
		<-time.After(60 * time.Second)
		m, err := getMetrics()
		if err != nil {
			fmt.Println(err)
			fmt.Println(err)
			continue
		}
		// m.Print()
		err = m.Rec()
		if err != nil {
			fmt.Println(err)
			fmt.Println(err)
			continue
		}
	}
}

func (sm *SystemMetrics) Rec() error {
	query := `
		INSERT INTO system_metrics 
		(memory_alloc, memory_total_alloc, memory_sys, num_goroutines, cpu_load, uptime, num_cpu, os, arch, total_memory, memory_usage_pct)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := postgree.MainDBX.Exec(query, sm.MemoryAlloc, sm.MemoryTotalAlloc, sm.MemorySys, sm.NumGoroutines,
		sm.CPULoad, sm.Uptime, sm.NumCPU, sm.OS, sm.Arch, sm.TotalMemory, sm.MemoryUsagePct)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ошибка при вставке данных: %w", err)
	}
	return nil
}
