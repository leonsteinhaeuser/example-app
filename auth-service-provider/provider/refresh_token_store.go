package provider

import (
	"context"
	"encoding/json"

	"github.com/leonsteinhaeuser/example-app/internal/keystore"
	"github.com/zitadel/oidc/v2/example/server/storage"
)

type RefreshTokenStore struct {
	kv keystore.KeyStore
}

func prefixRefreshToken(key string) string {
	return "rtk_" + key
}

func (s *RefreshTokenStore) Get(ctx context.Context, key string) (*storage.Token, error) {
	token := &storage.Token{}
	tk, err := s.kv.Get(ctx, prefixRefreshToken(key))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tk, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *RefreshTokenStore) Put(ctx context.Context, key string, token *storage.Token) error {
	tk, err := json.Marshal(token)
	if err != nil {
		return err
	}
	err = s.kv.Set(ctx, prefixRefreshToken(key), tk, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *RefreshTokenStore) Delete(ctx context.Context, key string) error {
	return s.kv.Delete(ctx, prefixRefreshToken(key))
}
