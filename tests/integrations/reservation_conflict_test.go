package integrations

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"booking/tests/testutil"
)

func TestReservations_Conflict409(t *testing.T) {
	ctx := context.Background()

	if err := testutil.ResetDB(ctx, db); err != nil {
		t.Fatal(err)
	}

	adminEmail := "admin@example.com"
	adminPass := "password123"
	//регистрация админа
	if status, body, err := testutil.DoJSON("POST", env.BaseURL+"/auth/register", "", map[string]any{
		"email":    adminEmail,
		"password": adminPass,
		"name":     "Admin",
	}); err != nil || status != 201 {
		t.Fatalf("register admin: status=%d err=%v body=%s", status, err, string(body))
	}

	if err := testutil.PromoteAdmin(ctx, db, adminEmail); err != nil {
		t.Fatal(err)
	}
	//логин под админом
	status, body, err := testutil.DoJSON("POST", env.BaseURL+"/auth/login", "", map[string]any{
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

	//создание комнаты
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

	//парсинг room.ID
	var room struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &room); err != nil {
		t.Fatalf("unmarshal room response: %v body=%s", err, string(body))
	}
	if room.ID == 0 {
		t.Fatalf("cannot parse room id from: %s", string(body))
	}

	//создание обычного пользователя
	userEmail := "u1@example.com"
	userPass := "password123"

	//регистрация обычного пользователя
	if status, body, err := testutil.DoJSON("POST", env.BaseURL+"/auth/register", "", map[string]any{
		"email":    userEmail,
		"password": userPass,
		"name":     "User",
	}); err != nil || status != 201 {
		t.Fatalf("register user: status=%d err=%v body=%s", status, err, string(body))
	}

	//логин обычного пользователя
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
	//формирование времени бронирования
	start := time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second)
	end := start.Add(1 * time.Hour)

	//создание первого бронирования
	status, body, err = testutil.DoJSON("POST", env.BaseURL+"/reservations", userToken, map[string]any{
		"room_id":    room.ID,
		"start_time": start.Format(time.RFC3339),
		"end_time":   end.Format(time.RFC3339),
	})
	if err != nil {
		t.Fatal(err)
	}
	if status != 201 {
		t.Fatalf("create reservation1: status=%d body=%s", status, string(body))
	}

	//создание конфликтующего бронирования
	status, body, err = testutil.DoJSON("POST", env.BaseURL+"/reservations", userToken, map[string]any{
		"room_id":    room.ID,
		"start_time": start.Add(30 * time.Minute).Format(time.RFC3339),
		"end_time":   end.Add(30 * time.Minute).Format(time.RFC3339),
	})
	if err != nil {
		t.Fatal(err)
	}
	if status != 409 {
		t.Fatalf("expected 409 conflict, got %d body=%s", status, string(body))
	}
}