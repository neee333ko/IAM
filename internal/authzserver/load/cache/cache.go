package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/neee333ko/IAM/internal/authzserver/store"
	pb "github.com/neee333ko/api/proto/v1"
	"github.com/neee333ko/log"
	"github.com/ory/ladon"
)

type Cache struct {
	secretMutex *sync.RWMutex
	policyMutex *sync.RWMutex
	secrets     *ristretto.Cache
	policies    *ristretto.Cache
	client      store.Factory
}

var (
	cache             *Cache
	once              sync.Once
	ErrSecretNotFound = errors.New("secret not found")
	ErrPolicyNotFound = errors.New("policy not found")
)

func GetCacheInsOr() *Cache {
	once.Do(func() {
		cache = new(Cache)
		cache.secretMutex = new(sync.RWMutex)
		cache.policyMutex = new(sync.RWMutex)
		cache.client = store.Client()

		config := &ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		}

		var err error
		cache.secrets, err = ristretto.NewCache(config)
		if err != nil {
			log.Fatalf("cache initialize failed: %s", err.Error())
		}

		cache.policies, err = ristretto.NewCache(config)
		if err != nil {
			log.Fatalf("cache initialize failed: %s", err.Error())
		}
	})

	return cache
}

func (cache *Cache) GetSecret(id string) (*pb.SecretInfo, error) {
	cache.secretMutex.RLock()
	defer cache.secretMutex.RUnlock()

	value, found := cache.secrets.Get(id)
	if !found {
		return nil, ErrSecretNotFound
	}

	key, _ := value.(*pb.SecretInfo)

	return key, nil
}

func (cache *Cache) GetPolicy(username string) ([]*ladon.DefaultPolicy, error) {
	cache.policyMutex.RLock()
	defer cache.policyMutex.RUnlock()

	value, found := cache.policies.Get(username)
	if !found {
		return nil, ErrPolicyNotFound
	}

	policies, _ := value.([]*ladon.DefaultPolicy)

	return policies, nil
}

func (cache *Cache) Reload() error {
	log.Info("secrets and policies reloading...")

	cache.secretMutex.Lock()
	cache.policyMutex.Lock()
	defer cache.secretMutex.Unlock()
	defer cache.policyMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	policies, err := cache.client.NewPolicyStore().List(ctx)
	if err != nil {
		return err
	}

	cache.policies.Clear()

	for key, value := range *policies {
		cache.policies.Set(key, value, 1)
	}

	secrets, err := cache.client.NewSecretStore().List(ctx)
	if err != nil {
		return err
	}

	cache.secrets.Clear()

	for key, value := range *secrets {
		cache.secrets.Set(key, value, 1)
	}

	log.Info("secrets and policies reloaded.")

	return nil
}
