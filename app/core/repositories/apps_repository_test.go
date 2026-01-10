package repositories

import (
	"context"
	"testing"

	"github.com/brunojet/my-store-go/app/core/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestAppsRepository_CRUDAndPatch(t *testing.T) {
	ctx := context.Background()

	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := gormDB.AutoMigrate(&models.App{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repo := NewAppsRepository(gormDB)

	created, err := repo.Create(ctx, &models.App{Nome: "A", Descricao: "D", CodigoParceiroExterno: "C"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected ID to be set")
	}

	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil || got.Nome != "A" {
		t.Fatalf("unexpected get result: %+v", got)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}

	newNome := "B"
	patched, err := repo.Patch(ctx, created.ID, map[string]any{"nome": newNome})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if patched == nil || patched.Nome != "B" {
		t.Fatalf("unexpected patched result: %+v", patched)
	}

	got2, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get after patch: %v", err)
	}
	if got2 == nil || got2.Nome != "B" {
		t.Fatalf("unexpected get after patch: %+v", got2)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	missing, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get missing: %v", err)
	}
	if missing != nil {
		t.Fatalf("expected nil after delete, got %+v", missing)
	}
}
