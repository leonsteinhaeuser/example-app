package provider

import (
	"context"
	"encoding/json"

	"github.com/leonsteinhaeuser/example-app/internal/keystore"
	"github.com/zitadel/oidc/v2/example/server/storage"
)

type TokenStore struct {
	kv keystore.KeyStore
}

func prefixUserToken(key string) string {
	return "utk_" + key
}

func (s *TokenStore) Get(ctx context.Context, key string) (*storage.Token, error) {
	token := &storage.Token{}
	tk, err := s.kv.Get(ctx, prefixUserToken(key))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tk, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *TokenStore) Put(ctx context.Context, key string, token *storage.Token) error {
	tk, err := json.Marshal(token)
	if err != nil {
		return err
	}
	err = s.kv.Set(ctx, prefixUserToken(key), tk, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *TokenStore) Delete(ctx context.Context, key string) error {
	return s.kv.Delete(ctx, prefixUserToken(key))
}
