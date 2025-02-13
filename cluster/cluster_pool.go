package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/resp/client"
)

type connectionFactory struct {
	node string
}

func (f connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	c, err := client.MakeClient(f.node)
	if err != nil {
		return nil, err
	}
	c.Start()
	return pool.NewPooledObject(c), nil
}

func (f connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("object is not a client")
	}
	c.Close()
	return nil
}

func (f connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	//不验证
	return true
}

func (f connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (f connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
