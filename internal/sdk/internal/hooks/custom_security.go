package hooks

import (
	"errors"
	"log"
	"net/http"

	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
)

// Manages SDK custom security handling.
//
// Supports:
//   - Setting HTTP Authorization header using Service Token ID and Service Token.
//
// This can be extended in the future to support additional security mechanisms.
type CustomSecurityHook struct{}

func (i *CustomSecurityHook) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	logger := log.Default()
	securityObj, err := hookCtx.SDKConfiguration.Security(req.Context())

	if err != nil {
		return req, nil
	}

	security, ok := securityObj.(shared.Security)

	if !ok {
		return req, nil
	}

	serviceToken := security.GetServiceToken()
	serviceTokenID := security.GetServiceTokenID()

	if serviceToken == "" || serviceTokenID == "" {
		return nil, errors.New("missing Service Token and Service Token ID credentials")
	}

	logger.Printf("CustomSecurityHook: Setting HTTP Authorization header with Service Token ID: %s", serviceTokenID)
	req.Header.Set("Authorization", serviceTokenID+":"+serviceToken)

	return req, nil
}
