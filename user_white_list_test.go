package onappgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testID          = 1
	testUserID      = 1
	testIP          = "0.0.0.0/0"
	testDescription = "Test description for IP"
	testTime        = time.Date(2020, 4, 20, 0, 0, 0, 0, time.UTC)
	testTimeString  = testTime.Format(time.RFC3339)
)

func TestUserWhiteList_Create(t *testing.T) {
	setup()
	defer teardown()

	want := &UserWhiteList{
		CreatedAt:   "2020-03-20T12:57:58.000+02:00",
		Description: "test",
		ID:          2,
		IP:          "1.1.1.1/0",
		UpdatedAt:   "2020-03-20T12:57:58.000+02:00",
		UserID:      testUserID,
	}

	createRequest := &UserWhiteListCreateRequest{
		Description: "test",
		IP:          "1.1.1.1/0",
	}

	createResponseJSON := `{"user_white_list":{"id":2,"user_id":1,"ip":"1.1.1.1/0","description":"test","created_at":"2020-03-20T12:57:58.000+02:00","updated_at":"2020-03-20T12:57:58.000+02:00"}}`

	path := fmt.Sprintf("/users/%d/user_white_lists.json", testUserID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		res := make(map[string]interface{})

		expected := map[string]interface{}{
			"description": "test",
			"ip":          "1.1.1.1/0",
		}

		res["user_white_list"] = expected

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(v, res) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, res)
		}

		fmt.Fprint(w, createResponseJSON)
	})

	got, _, err := client.UserWhiteLists.Create(ctx, testUserID, createRequest)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestUserWhiteList_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &UserWhiteList{
		CreatedAt:   "2020-03-20T12:57:58.000+02:00",
		Description: "",
		ID:          testID,
		IP:          testIP,
		UpdatedAt:   "2020-03-20T12:57:58.000+02:00",
		UserID:      testUserID,
	}

	getResponseJSON := `
	{
		"user_white_list": {
			"created_at": "2020-03-20T12:57:58.000+02:00",
			"description": "",
			"id": 1,
			"ip": "0.0.0.0/0",
			"updated_at": "2020-03-20T12:57:58.000+02:00",
			"user_id": 1
		}
	}`

	path := fmt.Sprintf("/users/%d/user_white_lists/%d.json", testUserID, testID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, getResponseJSON)
	})

	got, _, err := client.UserWhiteLists.Get(ctx, testUserID, testID)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestUserWhiteList_Delete(t *testing.T) {
	setup()
	defer teardown()

	path := fmt.Sprintf("/users/%d/user_white_lists/%d.json", testUserID, testID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.UserWhiteLists.Delete(ctx, testUserID, testID, nil)
	require.NoError(t, err)
}

func TestUserWhiteList_List(t *testing.T) {
	setup()
	defer teardown()

	wantUserWhiteLists := []UserWhiteList{
		{
			CreatedAt:   testTimeString,
			Description: testDescription,
			ID:          testID,
			IP:          testIP,
			UpdatedAt:   testTimeString,
			UserID:      testUserID,
		},
	}

	getResponseJSON := `[
		{
			"user_white_list": {
				"created_at": "` + testTimeString + `",
				"description": "` + testDescription + `",
				"id": ` + strconv.Itoa(testID) + `,
				"ip": "` + testIP + `",
				"updated_at": "` + testTimeString + `",
				"user_id": ` + strconv.Itoa(testUserID) + `
			}
		}
	]`

	path := fmt.Sprintf("/users/%d/user_white_lists.json", testUserID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, getResponseJSON)
	})

	got, _, err := client.UserWhiteLists.List(ctx, testUserID, nil)

	require.NoError(t, err)
	require.Equal(t, wantUserWhiteLists, got)
}
