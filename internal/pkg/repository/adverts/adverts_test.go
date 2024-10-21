package adverts

// import (
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/pkg/domain"
// )

// func TestAdvertsList_Add(t *testing.T) {
// 	t.Parallel()

// 	list := NewAdvertsList()
// 	advert := &domain.Advert{Title: "Test Advert"}

// 	list.Add(advert)

// 	if len(list.adverts) != 1 {
// 		t.Fatalf("expected 1 advert, got %d", len(list.adverts))
// 	}

// 	if list.adverts[0].ID != 1 {
// 		t.Fatalf("expected ID 1, got %d", list.adverts[0].ID)
// 	}
// }

// func TestAdvertsList_Update(t *testing.T) {
// 	t.Parallel()

// 	list := NewAdvertsList()
// 	advert := &domain.Advert{ID: 1, Title: "Test Advert"}

// 	list.Add(advert)

// 	updatedAdvert := &domain.Advert{ID: 1, Title: "Updated Advert"}

// 	err := list.Update(updatedAdvert)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if list.adverts[0].Title != "Updated Advert" {
// 		t.Fatalf("expected 'Updated Advert', got %s", list.adverts[0].Title)
// 	}
// }

// func TestAdvertsList_DeleteAdvert(t *testing.T) {
// 	t.Parallel()

// 	list := NewAdvertsList()
// 	advert := &Advert{ID: 1, Title: "Test Advert"}

// 	list.Add(advert)

// 	err := list.DeleteAdvert(1)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if len(list.adverts) != 0 {
// 		t.Fatalf("expected 0 adverts, got %d", len(list.adverts))
// 	}
// }

// func TestAdvertsList_GetAdverts(t *testing.T) {
// 	t.Parallel()

// 	list := NewAdvertsList()
// 	advert := &Advert{ID: 1, Title: "Test Advert"}

// 	list.Add(advert)

// 	adverts, err := list.GetAdverts()
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if len(adverts) != 1 {
// 		t.Fatalf("expected 1 advert, got %d", len(adverts))
// 	}

// 	if adverts[0].ID != 1 {
// 		t.Fatalf("expected ID 1, got %d", adverts[0].ID)
// 	}
// }

// func TestAdvertsList_GetAdvertByID(t *testing.T) {
// 	t.Parallel()

// 	list := NewAdvertsList()
// 	advert := &Advert{ID: 1, Title: "Test Advert"}

// 	list.Add(advert)

// 	result, err := list.GetAdvertByID(1)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if result.ID != 1 {
// 		t.Fatalf("expected ID 1, got %d", result.ID)
// 	}

// 	_, err = list.GetAdvertByID(2)
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}
// }
