package integrations

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"booking/tests/testutil"
)

//тест на гонку
func TestReservations_RaceOnlyOneWins(t *testing.T) {
	ctx := context.Background()

	if err := testutil.ResetDB(ctx, db); err != nil {
		t.Fatal(err)
	}

	adminEmail := "admin@example.com"
	adminPass := "password123"

	status, body, err := testutil.DoJSON("POST", env.BaseURL+"/auth/register", "", map[string]any{
		"email":    adminEmail,
		"password": adminPass,
		"name":     "Admin",
	})
	if err != nil || status != 201 {
		t.Fatalf("register admin: status=%d err=%v body=%s", status, err, string(body))
	}

	if err := testutil.PromoteAdmin(ctx, db, adminEmail); err != nil {
		t.Fatal(err)
	}

	status, body, err = testutil.DoJSON("POST", env.BaseURL+"/auth/login", "", map[string]any{
		"email":    adminEmail,
		"password": adminPass,
	})
	if err != nil {
		t.Fatal(err)
	}

	adminToken, err := testutil.MustTokenFromLoginResponse(status, body)
	if err != nil {
		t.Fatal(err)
	}

	status, body, err = testutil.DoJSON("POST", env.BaseURL+"/admin/rooms", adminToken, map[string]any{
		"name":        "Room A",
		"description": "test",
		"capacity":    4,
		"location":    "1st floor",
	})
	if err != nil {
		t.Fatal(err)
	}
	if status != 201 {
		t.Fatalf("create room: status=%d body=%s", status, string(body))
	}

	roomID := extractRoomID(t, body)

	userEmail := "user@example.com"
	userPass := "password123"

	status, body, err = testutil.DoJSON("POST", env.BaseURL+"/auth/register", "", map[string]any{
		"email":    userEmail,
		"password": userPass,
		"name":     "User",
	})
	if err != nil || status != 201 {
		t.Fatalf("register user: status=%d err=%v body=%s", status, err, string(body))
	}

	status, body, err = testutil.DoJSON("POST", env.BaseURL+"/auth/login", "", map[string]any{
		"email":    userEmail,
		"password": userPass,
	})
	if err != nil {
		t.Fatal(err)
	}

	userToken, err := testutil.MustTokenFromLoginResponse(status, body)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second)
	end := start.Add(1 * time.Hour)

	const n = 10 //будет запущено 10 конкурентных запросов
	var ok201 int32 //сколько запросов получили 201 Created
	var conflict409 int32 //сколько получили 409 Conflict
	var other int32 //сколько получили что-то неожиданное

	var wg sync.WaitGroup
	wg.Add(n)

	startCh := make(chan struct{})

	//запуск горутин в цикле
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			<-startCh
			//каждая горутина пытается создать одну и ту же бронь
			st, _, err := testutil.DoJSON("POST", env.BaseURL+"/reservations", userToken, map[string]any{
				"room_id":    roomID,
				"start_time": start.Format(time.RFC3339),
				"end_time":   end.Format(time.RFC3339),
			})

			if err != nil {
				atomic.AddInt32(&other, 1)
				return
			}

			switch st {
			case 201:
				atomic.AddInt32(&ok201, 1)
			case 409:
				atomic.AddInt32(&conflict409, 1)
			default:
				atomic.AddInt32(&other, 1)
			}
		}()
	}

	close(startCh)
	wg.Wait()

	if ok201 != 1 {
		t.Fatalf("expected exactly 1 success (201), got %d (409=%d other=%d)", ok201, conflict409, other)
	}
	if conflict409 != n-1 {
		t.Fatalf("expected %d conflicts (409), got %d (201=%d other=%d)", n-1, conflict409, ok201, other)
	}
	if other != 0 {
		t.Fatalf("expected other=0, got %d (201=%d 409=%d)", other, ok201, conflict409)
	}
}

func extractRoomID(t *testing.T, body []byte) int64 {
	t.Helper()

	var v struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatalf("unmarshal room body: %v body=%s", err, string(body))
	}
	if v.ID == 0 {
		t.Fatalf("cannot parse room id from: %s", string(body))
	}
	return v.ID
}