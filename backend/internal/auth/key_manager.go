package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"github.com/koo-arch/adjusta-backend/cache"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/jwtkey"
)

var jwtKeyTTL = 6 * 30 * 24 * time.Hour

type KeyManager struct {
	client *ent.Client
	cache *cache.Cache
}

func NewKeyManager(client *ent.Client, cache *cache.Cache) *KeyManager {
	return &KeyManager{
		client: client,
		cache: cache,
	}
}

func (km *KeyManager) GenerateRamdomKey(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (km *KeyManager) GenerateJWTKey(ctx context.Context, keyType string) (error) {
	key, err := km.GenerateRamdomKey(32)
	if err != nil {
		return err
	}

	// JWTキーの保存
	_, err = km.client.JWTKey.
		Create().
		SetType(keyType).
		SetKey(key).
		SetExpiresAt(time.Now().Add(jwtKeyTTL)).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (km *KeyManager) GetJWTKey(ctx context.Context, keyType string) ([]byte, error) {
	// キャッシュから取得
	if keyValue, found := km.GetJWTKeyFromCache(keyType); found {
		return []byte(keyValue), nil
	}

	// 有効期限内の最新のキーを取得
	jwtKey, err := km.client.JWTKey.
		Query().
		Where(
			jwtkey.Type(keyType),
			jwtkey.ExpiresAtGT(time.Now()),
		).
		Order(ent.Desc(jwtkey.FieldCreatedAt)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	// キャッシュに保存
	km.cache.JWTKeyCache.Set(keyType, jwtKey.Key, 5*time.Minute)

	byteKey := []byte(jwtKey.Key)

	return byteKey, nil
}

func (km *KeyManager) GetJWTKeyFromCache(keyType string) (string, bool) {
	keyValue, found := km.cache.JWTKeyCache.Get(keyType)
	if found {
		return keyValue.(string), true
	}

	return "", false
}

func (km *KeyManager) GetJWTKeys(ctx context.Context, keyType string) ([]*ent.JWTKey, error) {
	jwtKeys, err := km.client.JWTKey.
		Query().
		Where(
			jwtkey.Type(keyType)).
			Order(ent.Desc(jwtkey.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return jwtKeys, nil
}

func (km *KeyManager) DeleteJWTKeys(ctx context.Context, keyType string) (error) {
	// キーの削除猶予期間を設定
	expiresDuration := jwtKeyTTL

	switch keyType {
	case "access":
		expiresDuration += accessTokenTTL
	case "refresh":
		expiresDuration += refreshTokenTTL
	default:
		return nil
	}

	// keyTypeに応じて設定された期限を過ぎたキーを削除
	_, err := km.client.JWTKey.
		Delete().
		Where(
			jwtkey.Type(keyType),
			jwtkey.ExpiresAtLT(time.Now().Add(-expiresDuration)),
		).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (km *KeyManager) InitializeJWTKeys(ctx context.Context) (error) {
	// アクセスキーが作成されていない場合は新規作成
	_, err := km.GetJWTKey(ctx, "access")
	if err != nil {
		log.Println("access key not found - generating new key")
		err = km.GenerateJWTKey(ctx, "access")
		if err != nil {
			return err
		}
	} else {
		log.Println("access key already exists")
	}

	// リフレッシュキーが作成されていない場合は新規作成
	_, err = km.GetJWTKey(ctx, "refresh")
	if err != nil {
		log.Println("refresh key not found - generating new key")
		err = km.GenerateJWTKey(ctx, "refresh")
		if err != nil {
			return err
		}
	} else {
		log.Println("refresh key already exists")
	}
	return nil
}