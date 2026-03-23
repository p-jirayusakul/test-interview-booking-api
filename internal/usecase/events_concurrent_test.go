package usecase

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/books"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/events"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/postgres"
)

func TestConcurrentBooking(t *testing.T) {

	ctx := context.Background()

	eventIDStr := "019d1a7c-67f5-7376-a505-61718ec959f5"
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		t.Fatalf("failed to parse event ID: %v", err)
	}

	dbConn, err := postgres.NewConnection(postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "1234",
		DBName:   "events_booking_db",
	})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	var wg sync.WaitGroup
	totalRequests := 100

	timeLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatalf("failed to load timezone: %v", err)
	}

	eventsRepo := events.NewEventsRepository(dbConn)
	booksRepo := books.NewBooksRepository(dbConn)
	txManager := postgres.NewTransactionManager(dbConn)
	eventsUseCase := NewEventsUseCase(eventsRepo, booksRepo, txManager, timeLocation)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			userID := uuid.New()

			_, err := eventsUseCase.BookEvent(ctx, eventID, userID)
			if err != nil {
				// optional: log error
				fmt.Println(err)
			}
		}(i)
	}

	wg.Wait()

	// assert
	event, _ := eventsRepo.GetEvenById(ctx, eventID)

	if event.BookedCount > 50 {
		t.Fatalf("overbooking detected")
	}

	if event.WaitlistCount > 5 {
		t.Fatalf("waitlist overflow")
	}
}
