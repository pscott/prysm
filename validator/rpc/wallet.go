package rpc

import (
	"context"
	"path/filepath"

	ptypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	pb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	v2 "github.com/prysmaticlabs/prysm/validator/accounts/v2"
	"github.com/prysmaticlabs/prysm/validator/flags"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WalletConfig --
func (s *Server) WalletConfig(ctx context.Context, _ *ptypes.Empty) (*pb.WalletResponse, error) {
	defaultWalletPath := filepath.Join(flags.DefaultValidatorDir(), flags.WalletDefaultDirName)
	err := v2.WalletExists(defaultWalletPath)
	if err != nil && errors.Is(err, v2.ErrNoWalletFound) {
		// If no wallet is found, we simply return an empty response.
		return nil, nil
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not check if wallet exists: %v", err)
	}
	return &pb.WalletResponse{
		WalletPath:       defaultWalletPath,
		KeymanagerConfig: nil, // Fill in by reading from disk.
	}, nil
}
