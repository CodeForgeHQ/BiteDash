package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

const baseURL = "https://fakerestaurantapi.runasp.net"

type Restaurant struct {
	RestaurantID   int    `json:"restaurantID"`
	RestaurantName string `json:"restaurantName"`
	Description    string `json:"description"`
	Address        string `json:"address"`
	Type           string `json:"category"`
	ParkingLot     bool   `json:"parkingLot"`
}

type MenuItem struct {
	ItemID          int     `json:"itemID"`
	ItemName        string  `json:"itemName"`
	ItemDescription string  `json:"itemDescription"`
	ItemPrice       float64 `json:"itemPrice"`
	RestaurantName  string  `json:"restaurantName"`
	RestaurantID    int     `json:"restaurantID"`
	ImageURL        string  `json:"imageUrl"`
}

func FetchRestaurants(ctx context.Context) ([]Restaurant, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := baseURL + "/api/Restaurant"
	var lastErr error

	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if !isRetryableError(err) {
				return nil, err
			}
		} else {
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode == http.StatusOK {
				var restaurants []Restaurant
				if err := json.NewDecoder(resp.Body).Decode(&restaurants); err != nil {
					return nil, err
				}
				return restaurants, nil
			}

			lastErr = fmt.Errorf("api returned status: %d", resp.StatusCode)
			if !isRetryableStatus(resp.StatusCode) {
				return nil, lastErr
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 300 * time.Millisecond):
		}
	}

	return nil, fmt.Errorf("fetch restaurants failed after retries: %w", lastErr)
}

func isRetryableError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func isRetryableStatus(statusCode int) bool {
	return statusCode >= http.StatusInternalServerError
}

func FetchAllMenuItems(ctx context.Context) ([]MenuItem, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/Restaurant/items", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch menu items: unexpected status %d", resp.StatusCode)
	}

	var items []MenuItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("decode menu items: %w", err)
	}

	return items, nil
}
