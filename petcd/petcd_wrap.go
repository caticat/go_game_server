package petcd

import (
	"context"

	"go.etcd.io/etcd/clientv3"
)

func getRevision() (int64, error) {
	cli := getClient()
	if cli == nil {
		return 0, ErrorNilClient
	}

	ctx, cancel := context.WithTimeout(context.Background(), getConfigOperationTimeout())
	defer cancel()
	resp, err := cli.Get(ctx, "", clientv3.WithPrefix(), clientv3.WithLimit(1))
	if err != nil {
		return 0, err
	}

	return resp.Header.Revision, nil
}
