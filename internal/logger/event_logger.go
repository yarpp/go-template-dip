package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type EventLogger struct {
	eventsChan chan string
	done       chan struct{}
	file       *os.File
	wg         sync.WaitGroup
	once       sync.Once
}

func NewEventLogger(filePath string) *EventLogger {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("event-logger: failed to open log file %s: %v", filePath, err)
		return &EventLogger{
			eventsChan: make(chan string, 100),
			done:       make(chan struct{}),
		}
	}
	return &EventLogger{
		eventsChan: make(chan string, 100),
		done:       make(chan struct{}),
		file:       file,
	}
}

// LogEvent отправляет событие в канал (неблокирующе)
func (el *EventLogger) LogEvent(event string) {
	if el == nil {
		return
	}
	select {
	case el.eventsChan <- event:
	default:
		log.Printf("event-logger: channel full, dropping event: %s", event)
	}
}

// Start запускает worker в отдельной горутине
func (el *EventLogger) Start() {
	if el == nil {
		return
	}
	el.wg.Add(1)
	go el.worker()
}

func (el *EventLogger) worker() {
	defer el.wg.Done()
	defer func() {
		if el.file != nil {
			_ = el.file.Close()
		}
		close(el.done)
	}()

	for event := range el.eventsChan {
		// искусственная задержка 1 секунда по требованиям задания
		time.Sleep(1 * time.Second)
		line := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), event)
		if el.file != nil {
			if _, err := el.file.WriteString(line); err != nil {
				log.Printf("event-logger: write error: %v", err)
			}
			_ = el.file.Sync()
		}
		log.Printf("event-logger: %s", line)
	}
}

// Stop корректно завершает работу worker'а
func (el *EventLogger) Stop() {
	if el == nil {
		return
	}
	el.once.Do(func() {
		close(el.eventsChan)
		el.wg.Wait()
		log.Println("event-logger: stopped")
	})
}
