package rpc

import (
	"context"

	ptypes "github.com/gogo/protobuf/types"
	pb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	v2 "github.com/prysmaticlabs/prysm/validator/accounts/v2"
)

// WalletConfig --
func (s *Server) WalletConfig(ctx context.Context, _ *ptypes.Empty) (*pb.WalletResponse, error) {
	wallet, err := v2.OpenWallet()
	//return &pb.WalletResponse{
	//	WalletPath:           "",
	//	KeymanagerConfig:     nil,
	//}, nil
	return nil, nil
}
