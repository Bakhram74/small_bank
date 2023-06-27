package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/Bakhram74/small_bank/db/mock"
	db "github.com/Bakhram74/small_bank/db/sqlc"
	"github.com/Bakhram74/small_bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
func TestGetAccountApi(t *testing.T) {
	account := randomAccount()
	testCase := []struct {
		name          string
		accountId     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			"ok",
			account.ID,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			"NotFound",
			account.ID,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			"InternalServer",
			account.ID,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			"BadRequest",
			0,
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountId)

			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
