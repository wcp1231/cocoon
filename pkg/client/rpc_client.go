package client

import (
	"cocoon/pkg/model/rpc"
	"context"
	"github.com/smallnest/rpcx/client"
)

type RpcClient struct {
	client client.XClient
}

func NewRpcClient(remote string) *RpcClient {
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+remote, "")
	c := client.NewXClient(rpc.COCOON_SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)

	return &RpcClient{
		client: c,
	}
}

func (r *RpcClient) GenerateUploadRequest() {

}

func (r *RpcClient) Upload(ctx context.Context, req *rpc.UploadReq) {
	resp := &rpc.UploadResp{}
	call, err := r.client.Go(ctx, "Upload", req, resp, nil)
	if err != nil {
		//TODO
		// log.Fatalf("failed to call: %v", err)
	}

	replyCall := <-call.Done
	if replyCall.Error != nil {
		// TODO
		// log.Fatalf("failed to call: %v", replyCall.Error)
	} else {
		// TODO
	}
}
