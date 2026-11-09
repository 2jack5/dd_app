package callback_test

import (
	"dd_web/back/test/testkit"
	"net/http"
	"testing"
)

func TestPaymentCallbackDoesNotReturnFakeSuccess(t *testing.T) {
	response := testkit.RequestJSON(t, http.MethodPost, "/api/v1/callbacks/wechat/pay", map[string]string{"event": "test"})
	testkit.Require(t, response, http.StatusServiceUnavailable, "CHANNEL_NOT_CONFIGURED")
}
