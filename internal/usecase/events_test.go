package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockEventsRepository struct {
	CreateEventFunc         func(ctx context.Context, payload *domain.CreateEvent) error
	GetEvenByIdFunc         func(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	TXGetEventForUpdateFunc func(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.Event, error)
	TXExistsBookingFunc     func(ctx context.Context, tx *gorm.DB, eventId, userId uuid.UUID) (bool, error)
	TXUpdateEventFunc       func(ctx context.Context, tx *gorm.DB, payload *domain.Event) error
	SearchEventsFunc        func(ctx context.Context, payload *domain.EventFilter) ([]*domain.Event, error)
	CountEventsFunc         func(ctx context.Context, search string) (int64, error)
}

func (m *mockEventsRepository) CreateEvent(ctx context.Context, payload *domain.CreateEvent) error {
	return m.CreateEventFunc(ctx, payload)
}

func (m *mockEventsRepository) GetEvenById(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return m.GetEvenByIdFunc(ctx, id)
}

func (m *mockEventsRepository) TXGetEventForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.Event, error) {
	return m.TXGetEventForUpdateFunc(ctx, tx, id)
}

func (m *mockEventsRepository) TXExistsBooking(ctx context.Context, tx *gorm.DB, eventId, userId uuid.UUID) (bool, error) {
	return m.TXExistsBookingFunc(ctx, tx, eventId, userId)
}

func (m *mockEventsRepository) TXUpdateEvent(ctx context.Context, tx *gorm.DB, payload *domain.Event) error {
	return m.TXUpdateEventFunc(ctx, tx, payload)
}

func (m *mockEventsRepository) SearchEvents(ctx context.Context, payload *domain.EventFilter) ([]*domain.Event, error) {
	return m.SearchEventsFunc(ctx, payload)
}

func (m *mockEventsRepository) CountEvents(ctx context.Context, search string) (int64, error) {
	return m.CountEventsFunc(ctx, search)
}

type mockBooksRepository struct {
	TXCreateBookingFunc func(ctx context.Context, tx *gorm.DB, payload domain.Booking) error
}

func (m *mockBooksRepository) TXCreateBooking(ctx context.Context, tx *gorm.DB, payload domain.Booking) error {
	return m.TXCreateBookingFunc(ctx, tx, payload)
}

type mockTransactionManager struct {
	WithTransactionFunc func(ctx context.Context, fn func(tx *gorm.DB) error) error
}

func (m *mockTransactionManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if m.WithTransactionFunc != nil {
		return m.WithTransactionFunc(ctx, fn)
	}
	return fn(nil)
}

func TestEventsUseCase_CreateEvent(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name        string
		payload     *domain.CreateEvent
		mockRepoErr error
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			payload: &domain.CreateEvent{
				Name:          "Go Workshop",
				MaxSeats:      100,
				WaitlistLimit: 10,
				Price:         999,
				StartTime:     now.Add(time.Hour),
				EndTime:       now.Add(2 * time.Hour),
			},
			mockRepoErr: nil,
			wantErr:     false,
		},
		{
			name: "validation failed - empty name",
			payload: &domain.CreateEvent{
				Name:          "",
				MaxSeats:      100,
				WaitlistLimit: 10,
				Price:         999,
				StartTime:     now.Add(time.Hour),
				EndTime:       now.Add(2 * time.Hour),
			},
			wantErr:     true,
			errContains: "event name cannot be empty",
		},
		{
			name: "validation failed - max seats less than 1",
			payload: &domain.CreateEvent{
				Name:          "Go Workshop",
				MaxSeats:      0,
				WaitlistLimit: 10,
				Price:         999,
				StartTime:     now.Add(time.Hour),
				EndTime:       now.Add(2 * time.Hour),
			},
			wantErr:     true,
			errContains: "max seats cannot be less than 1",
		},
		{
			name: "validation failed - start after end",
			payload: &domain.CreateEvent{
				Name:          "Go Workshop",
				MaxSeats:      100,
				WaitlistLimit: 10,
				Price:         999,
				StartTime:     now.Add(2 * time.Hour),
				EndTime:       now.Add(time.Hour),
			},
			wantErr:     true,
			errContains: "start time cannot be after end time",
		},
		{
			name: "repository failed",
			payload: &domain.CreateEvent{
				Name:          "Go Workshop",
				MaxSeats:      100,
				WaitlistLimit: 10,
				Price:         999,
				StartTime:     now.Add(time.Hour),
				EndTime:       now.Add(2 * time.Hour),
			},
			mockRepoErr: errors.New("db error"),
			wantErr:     true,
			errContains: "failed to create event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventsRepo := &mockEventsRepository{
				CreateEventFunc: func(ctx context.Context, payload *domain.CreateEvent) error {
					return tt.mockRepoErr
				},
			}
			booksRepo := &mockBooksRepository{}
			txManager := &mockTransactionManager{}

			uc := NewEventsUseCase(eventsRepo, booksRepo, txManager, time.UTC)

			err := uc.CreateEvent(ctx, tt.payload)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEventsUseCase_GetEvent(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	eventID := uuid.New()

	tests := []struct {
		name        string
		id          uuid.UUID
		mockResult  *domain.Event
		mockErr     error
		wantErr     bool
		errContains string
		wantEvent   *domain.Event
	}{
		{
			name:        "invalid input - nil uuid",
			id:          uuid.Nil,
			wantErr:     true,
			errContains: "event id cannot be nil",
		},
		{
			name:        "not found",
			id:          eventID,
			mockErr:     gorm.ErrRecordNotFound,
			wantErr:     true,
			errContains: "event not found",
		},
		{
			name:        "repository error",
			id:          eventID,
			mockErr:     errors.New("db error"),
			wantErr:     true,
			errContains: "failed to get event",
		},
		{
			name: "success",
			id:   eventID,
			mockResult: &domain.Event{
				ID:            eventID,
				Name:          "Go Workshop",
				MaxSeats:      100,
				WaitlistLimit: 10,
				BookedCount:   2,
				WaitlistCount: 0,
				Price:         999,
				StartTime:     now,
				EndTime:       now.Add(time.Hour),
				CreatedAt:     now,
				UpdatedAt:     nil,
			},
			wantEvent: &domain.Event{
				ID:            eventID,
				Name:          "Go Workshop",
				MaxSeats:      100,
				WaitlistLimit: 10,
				BookedCount:   2,
				WaitlistCount: 0,
				Price:         999,
				StartTime:     now,
				EndTime:       now.Add(time.Hour),
				CreatedAt:     now,
				UpdatedAt:     nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventsRepo := &mockEventsRepository{
				GetEvenByIdFunc: func(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
					return tt.mockResult, tt.mockErr
				},
			}

			uc := NewEventsUseCase(eventsRepo, &mockBooksRepository{}, &mockTransactionManager{}, time.UTC)

			got, err := uc.GetEvent(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.wantEvent, got)
		})
	}
}

func TestEventsUseCase_SearchEvents(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name         string
		payload      *domain.EventFilter
		mockItems    []*domain.Event
		mockTotal    int64
		mockItemsErr error
		mockTotalErr error
		wantErr      bool
		errContains  string
		wantPage     int
		wantLimit    int
		wantOffset   int
		wantTotal    int64
		wantHasNext  bool
		wantHasPrev  bool
		wantLen      int
	}{
		{
			name: "success - default page and limit",
			payload: &domain.EventFilter{
				Name:  "go",
				Page:  0,
				Limit: 0,
			},
			mockItems: []*domain.Event{
				{ID: uuid.New(), Name: "Go Workshop", StartTime: now, EndTime: now.Add(time.Hour)},
			},
			mockTotal:   1,
			wantErr:     false,
			wantPage:    1,
			wantLimit:   10,
			wantOffset:  0,
			wantTotal:   1,
			wantHasNext: false,
			wantHasPrev: false,
			wantLen:     1,
		},
		{
			name: "success - limit capped at 100",
			payload: &domain.EventFilter{
				Name:  "go",
				Page:  2,
				Limit: 200,
			},
			mockItems:   []*domain.Event{},
			mockTotal:   150,
			wantErr:     false,
			wantPage:    2,
			wantLimit:   100,
			wantOffset:  100,
			wantTotal:   150,
			wantHasNext: false,
			wantHasPrev: true,
			wantLen:     0,
		},
		{
			name: "repo search error",
			payload: &domain.EventFilter{
				Name:  "go",
				Page:  1,
				Limit: 10,
			},
			mockItemsErr: errors.New("search failed"),
			mockTotal:    1,
			wantErr:      true,
			errContains:  "failed to search event",
		},
		{
			name: "repo count error",
			payload: &domain.EventFilter{
				Name:  "go",
				Page:  1,
				Limit: 10,
			},
			mockItems:    []*domain.Event{},
			mockTotalErr: errors.New("count failed"),
			wantErr:      true,
			errContains:  "failed to count events",
		},
		{
			name: "success - no results",
			payload: &domain.EventFilter{
				Name:  "empty",
				Page:  1,
				Limit: 10,
			},
			mockItems:   []*domain.Event{},
			mockTotal:   0,
			wantErr:     false,
			wantPage:    1,
			wantLimit:   10,
			wantOffset:  0,
			wantTotal:   0,
			wantHasNext: false,
			wantHasPrev: false,
			wantLen:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventsRepo := &mockEventsRepository{
				SearchEventsFunc: func(ctx context.Context, payload *domain.EventFilter) ([]*domain.Event, error) {
					return tt.mockItems, tt.mockItemsErr
				},
				CountEventsFunc: func(ctx context.Context, search string) (int64, error) {
					return tt.mockTotal, tt.mockTotalErr
				},
			}

			uc := NewEventsUseCase(eventsRepo, &mockBooksRepository{}, &mockTransactionManager{}, time.UTC)

			got, err := uc.SearchEvents(ctx, tt.payload)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got.Items, tt.wantLen)

			require.Equal(t, tt.wantPage, got.Pagination.Page)
			require.Equal(t, tt.wantLimit, got.Pagination.PageSize)
			require.Equal(t, tt.wantTotal, got.Pagination.Total)
			require.Equal(t, tt.wantHasNext, got.Pagination.HasNext)
			require.Equal(t, tt.wantHasPrev, got.Pagination.HasPrevious)

			require.Equal(t, tt.wantOffset, tt.payload.Offset)
		})
	}
}

func TestEventsUseCase_BookEvent(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	eventID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name        string
		event       *domain.Event
		exists      bool
		txGetErr    error
		existsErr   error
		createErr   error
		updateErr   error
		wantErr     bool
		errContains string
		wantStatus  string
		wantBooked  int
		wantWait    int
	}{
		{
			name: "success confirmed",
			event: &domain.Event{
				ID:            eventID,
				Name:          "Go Workshop",
				MaxSeats:      50,
				WaitlistLimit: 5,
				BookedCount:   10,
				WaitlistCount: 1,
				StartTime:     now.Add(-time.Hour),
				EndTime:       now.Add(time.Hour),
			},
			exists:     false,
			wantErr:    false,
			wantStatus: string(domain.BookingStatusConfirmed),
			wantBooked: 11,
			wantWait:   1,
		},
		{
			name: "success waitlist",
			event: &domain.Event{
				ID:            eventID,
				Name:          "Go Workshop",
				MaxSeats:      10,
				WaitlistLimit: 5,
				BookedCount:   10,
				WaitlistCount: 1,
				StartTime:     now.Add(-time.Hour),
				EndTime:       now.Add(time.Hour),
			},
			exists:     false,
			wantErr:    false,
			wantStatus: string(domain.BookingStatusWaitlist),
			wantBooked: 10,
			wantWait:   2,
		},
		{
			name:        "event not found",
			txGetErr:    gorm.ErrRecordNotFound,
			wantErr:     true,
			errContains: "event not found",
		},
		{
			name: "booking window closed - before start",
			event: &domain.Event{
				ID:            eventID,
				MaxSeats:      50,
				WaitlistLimit: 5,
				BookedCount:   0,
				WaitlistCount: 0,
				StartTime:     now.Add(time.Hour),
				EndTime:       now.Add(2 * time.Hour),
			},
			wantErr:     true,
			errContains: "booking window closed",
		},
		{
			name: "already booked",
			event: &domain.Event{
				ID:            eventID,
				MaxSeats:      50,
				WaitlistLimit: 5,
				BookedCount:   10,
				WaitlistCount: 1,
				StartTime:     now.Add(-time.Hour),
				EndTime:       now.Add(time.Hour),
			},
			exists:      true,
			wantErr:     true,
			errContains: "you have already booked this event",
		},
		{
			name: "event full",
			event: &domain.Event{
				ID:            eventID,
				MaxSeats:      1,
				WaitlistLimit: 1,
				BookedCount:   1,
				WaitlistCount: 1,
				StartTime:     now.Add(-time.Hour),
				EndTime:       now.Add(time.Hour),
			},
			exists:      false,
			wantErr:     true,
			errContains: "event is full",
		},
		{
			name: "create booking duplicated key",
			event: &domain.Event{
				ID:            eventID,
				MaxSeats:      50,
				WaitlistLimit: 5,
				BookedCount:   10,
				WaitlistCount: 1,
				StartTime:     now.Add(-time.Hour),
				EndTime:       now.Add(time.Hour),
			},
			exists:      false,
			createErr:   gorm.ErrDuplicatedKey,
			wantErr:     true,
			errContains: "you have already booked this event",
		},
		{
			name: "update event failed",
			event: &domain.Event{
				ID:            eventID,
				MaxSeats:      50,
				WaitlistLimit: 5,
				BookedCount:   10,
				WaitlistCount: 1,
				StartTime:     now.Add(-time.Hour),
				EndTime:       now.Add(time.Hour),
			},
			exists:      false,
			updateErr:   errors.New("db error"),
			wantErr:     true,
			errContains: "failed to update event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventsRepo := &mockEventsRepository{
				TXGetEventForUpdateFunc: func(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.Event, error) {
					return tt.event, tt.txGetErr
				},
				TXExistsBookingFunc: func(ctx context.Context, tx *gorm.DB, eventId, userId uuid.UUID) (bool, error) {
					return tt.exists, tt.existsErr
				},
				TXUpdateEventFunc: func(ctx context.Context, tx *gorm.DB, payload *domain.Event) error {
					return tt.updateErr
				},
			}

			booksRepo := &mockBooksRepository{
				TXCreateBookingFunc: func(ctx context.Context, tx *gorm.DB, payload domain.Booking) error {
					return tt.createErr
				},
			}

			txManager := &mockTransactionManager{
				WithTransactionFunc: func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(&gorm.DB{})
				},
			}

			uc := NewEventsUseCase(eventsRepo, booksRepo, txManager, time.UTC)

			got, err := uc.BookEvent(ctx, eventID, userID)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, got)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantStatus, got)

			// assert state mutation on event object
			require.Equal(t, tt.wantBooked, tt.event.BookedCount)
			require.Equal(t, tt.wantWait, tt.event.WaitlistCount)
		})
	}
}
