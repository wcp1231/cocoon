package agent

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

func (r *RpcClient) ClientPostStart(ctx context.Context, appname, session string) {
	req := &rpc.PostStartReq{
		Appname: appname,
		Session: session,
	}
	resp := &rpc.PostStartResp{}
	err := r.client.Call(ctx, "ClientPostStart", req, resp)
	if err != nil {
		// TODO
	}
	if resp.Error != nil {
		// TODO
	}
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

func (r *RpcClient) ConnClose(ctx context.Context, req *rpc.ConnCloseReq) {
	resp := &rpc.ConnCloseResp{}
	err := r.client.Call(ctx, "ConnClose", req, resp)
	if err != nil {
		//TODO
		// log.Fatalf("failed to call: %v", err)
	}
}

func (r *RpcClient) RequestOutbound(ctx context.Context, req *rpc.OutboundReq) (*client.Call, error) {
	resp := &rpc.OutboundResp{}
	return r.client.Go(ctx, "RequestOutbound", req, resp, nil)
}

func (r *RpcClient) RecordRequestResponse(ctx context.Context, req *rpc.RecordReq) {
	resp := &rpc.RecordResp{}
	r.client.Go(ctx, "RecordRequestResponse", req, resp, nil)
	// TODO
}
