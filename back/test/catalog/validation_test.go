package catalog_test

import (
	"dd_web/back/test/testkit"
	"net/http"
	"testing"
)

func TestGamesRejectPageSizeAboveLimit(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/games?page_size=101", nil)
	testkit.Require(t, response, http.StatusBadRequest, "INVALID_ARGUMENT")
}

func TestGamesRejectNegativeCursor(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/games?cursor=-1", nil)
	testkit.Require(t, response, http.StatusBadRequest, "INVALID_ARGUMENT")
}

func TestRegionsRejectNonNumericGameID(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/games/not-id/regions", nil)
	testkit.Require(t, response, http.StatusBadRequest, "INVALID_ARGUMENT")
}

func TestRanksRejectNonPositiveGameID(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/games/0/ranks", nil)
	testkit.Require(t, response, http.StatusBadRequest, "INVALID_ARGUMENT")
}

func TestServiceCategoriesRequireGameID(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/service-categories", nil)
	testkit.Require(t, response, http.StatusBadRequest, "INVALID_ARGUMENT")
}

func TestServiceCategoriesRejectInvalidGameID(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/service-categories?game_id=invalid", nil)
	testkit.Require(t, response, http.StatusBadRequest, "INVALID_ARGUMENT")
}
