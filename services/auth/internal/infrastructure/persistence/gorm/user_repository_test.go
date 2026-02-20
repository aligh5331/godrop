package gorm

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aligh5331/godrop/services/auth/internal/domain"
	"github.com/aligh5331/godrop/services/auth/internal/domain/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	sqlDB, _ := db.DB()

	// Important for concurrency: allow multiple connections
	sqlDB.SetMaxOpenConns(10)

	t.Cleanup(func() {
		sqlDB.Close()
	})

	if err := db.AutoMigrate(&gUser{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}
func newTestUser(id, email string) *entity.User {
	now := time.Now()
	u, _ := entity.NewUser(
		id,
		"Ali",
		email,
		entity.HashedPassword("hashed"),
		now,
		now,
	)
	return u
}

func TestUserRepository_Create_And_FindById(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	user := newTestUser("id-1", "a@test.com")

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindById(ctx, "id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email() != "a@test.com" {
		t.Errorf("expected email a@test.com, got %s", found.Email())
	}
}

func TestUserRepository_FindById_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	_, err := repo.FindById(ctx, "missing")
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	user := newTestUser("id-2", "b@test.com")
	_ = repo.Create(ctx, user)

	found, err := repo.FindByEmail(ctx, "b@test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Id() != "id-2" {
		t.Errorf("expected id id-2, got %s", found.Id())
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	_, err := repo.FindByEmail(ctx, "missing@test.com")
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound")
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	u1 := newTestUser("id-1", "dup@test.com")
	u2 := newTestUser("id-2", "dup@test.com")

	if err := repo.Create(ctx, u1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := repo.Create(ctx, u2)
	if err == nil {
		t.Fatalf("expected unique constraint error")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	user := newTestUser("id-3", "c@test.com")
	_ = repo.Create(ctx, user)

	user.ChangeName("NewName", time.Now())
	user.ChangeEmail("new@test.com", time.Now())

	if err := repo.Update(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindById(ctx, "id-3")

	if found.Name() != "NewName" {
		t.Errorf("name not updated")
	}

	if found.Email() != "new@test.com" {
		t.Errorf("email not updated")
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	user := newTestUser("missing", "x@test.com")

	err := repo.Update(ctx, user)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	user := newTestUser("id-del", "del@test.com")
	_ = repo.Create(ctx, user)

	if err := repo.Delete(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := repo.FindById(ctx, "id-del")
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound after delete")
	}
}

func TestUserRepository_Create_Concurrent(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserRepository{db: db}
	ctx := context.Background()

	const workers = 20

	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			u := newTestUser(
				fmt.Sprintf("id-%d", i),
				fmt.Sprintf("u%d@test.com", i),
			)

			errCh <- func() error {
				tch := make(chan struct{})
				go func() {
					select {
					case <-time.After(5 * time.Second):
						tch <- struct{}{}
					}
				}()
				for {
					select {
					case <-tch:
						return errors.New("timeout")
					default:
						err := repo.Create(ctx, u)
						if err == nil {
							return nil
						}
						if strings.Contains(err.Error(), "database table is locked") {
							time.Sleep(time.Duration(5+rand.Intn(5)) * time.Millisecond)
							continue
						}
						return err
					}
				}
			}()
		}(i)
	}

	wg.Wait()
	close(errCh)

	// Assert no errors
	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Assert all rows exist
	var count int64
	if err := db.Model(&gUser{}).Count(&count).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}

	if count != workers {
		t.Fatalf("expected %d users, got %d", workers, count)
	}
}
