package middleware

import (
	"net/http"

	"github.com/AssetPortal/assets-api/pkg/adapters/auth"
	"github.com/AssetPortal/assets-api/pkg/adapters/tokens"
	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/ggicci/httpin"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type PolkadotAuth struct {
	tokensRepo tokens.Repository
	authClient auth.Client
	enabled    bool
}

// NewPolkadotAuth initializes PolkadotAuth with a TokenRepository.
func NewPolkadotAuth(tokensRepo tokens.Repository, authClient auth.Client, enabled bool) *PolkadotAuth {
	return &PolkadotAuth{
		tokensRepo: tokensRepo,
		authClient: authClient,
		enabled:    enabled,
	}
}

// Middleware for Polkadot authentication.
func (p *PolkadotAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !p.enabled {
			next.ServeHTTP(w, r)
		} else {
			headers := r.Context().Value(httpin.Input).(*model.AuthHeaders)
			address := r.Header.Get("X-Address")
			signature := r.Header.Get("X-Signature")
			message := r.Header.Get("X-Message")

			if address == "" || signature == "" || message == "" {
				render.JSON(w, r, model.NewResponseError("Missing authentication headers"))
				return
			}
			// Verify the token
			dbToken, err := p.tokensRepo.GetToken(r.Context(), headers.Message)
			if err != nil {
				logrus.Errorf("Error retrieving token: %s", err)
				render.JSON(w, r, model.NewResponseError("Cannot verify the token"))
				return
			}
			if dbToken == nil {
				render.JSON(w, r, model.NewResponseError("Message was not generated with /nonce"))
				return
			}
			if !dbToken.IsValid() {
				render.JSON(w, r, model.NewResponseError("Invalid or expired token"))
				return
			}

			// Verify the Polkadot signature
			auth, err := p.authClient.VerifySignature(r.Context(), headers.Message, headers.Address, headers.Signature)
			if err != nil {
				logrus.Errorf("Error verifying the signature: %s", err)
				render.JSON(w, r, model.NewResponseError("Cannot verify signature now"))
				return
			}
			if !auth.OK {
				render.JSON(w, r, model.NewResponseError("Invalid authentication: "+auth.Message))
				return
			}

			// Mark the token as used
			if err := p.tokensRepo.MarkTokenAsUsed(r.Context(), headers.Message); err != nil {
				logrus.Errorf("Error marking token as used: %s", err)
				render.JSON(w, r, model.NewResponseError("Failed to mark token as used"))
				return
			}

			next.ServeHTTP(w, r)
		}
	})
}
