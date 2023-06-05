package provider

import (
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
)

var (
	_ op.Client = (*ClientOP)(nil)
)

type Client struct {
	// ID is the client's unique identifier.
	ID string
	// Secret is the client's secret. It is used to authenticate the client when exchanging an authorization code for an access token.
	Secret string

	RedirectURIs       []string
	LogoutRedirectURIs []string
	ApplicationType    op.ApplicationType
	AuthMethod         oidc.AuthMethod
	ResponseTypes      []oidc.ResponseType
	GrantTypes         []oidc.GrantType
	AccessTokenType    op.AccessTokenType
	IDTokenLifetime    time.Duration

	Scopes []string
}

type ClientOP struct {
	Client *Client
}

func (c *ClientOP) GetID() string {
	return c.Client.ID
}

func (c *ClientOP) RedirectURIs() []string {
	return c.Client.RedirectURIs
}

func (c *ClientOP) PostLogoutRedirectURIs() []string {
	return c.Client.LogoutRedirectURIs
}

func (c *ClientOP) ApplicationType() op.ApplicationType {
	return c.Client.ApplicationType
}

func (c *ClientOP) AuthMethod() oidc.AuthMethod {
	return c.Client.AuthMethod
}

func (c *ClientOP) ResponseTypes() []oidc.ResponseType {
	return c.Client.ResponseTypes
}

func (c *ClientOP) GrantTypes() []oidc.GrantType {
	return c.Client.GrantTypes
}

func (c *ClientOP) LoginURL(state string) string {
	return ""
}

func (c *ClientOP) AccessTokenType() op.AccessTokenType {
	return c.Client.AccessTokenType
}

func (c *ClientOP) IDTokenLifetime() time.Duration {
	return c.Client.IDTokenLifetime
}

func (c *ClientOP) DevMode() bool {
	return false
}

func (c *ClientOP) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return nil
}

func (c *ClientOP) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return nil
}

func (c *ClientOP) IsScopeAllowed(scope string) bool {
	return true
}

func (c *ClientOP) IDTokenUserinfoClaimsAssertion() bool {
	return false
}

func (c *ClientOP) ClockSkew() time.Duration {
	return time.Minute
}
