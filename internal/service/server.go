package service

import (
	"git.miem.hse.ru/1206/app"
	"git.miem.hse.ru/1206/app/errs"
	"git.miem.hse.ru/1206/app/middleware"
	"git.miem.hse.ru/1206/material-library/internal/domain"
	"git.miem.hse.ru/1206/proto/pb"
	"google.golang.org/grpc"
)

type LibraryServer struct {
	pb.UnimplementedLibraryServer
	uc *domain.Usecase
}

func NewLibraryServer(cfg *app.GRPCConfig, uc *domain.Usecase, serverOptions ...grpc.UnaryServerInterceptor) (*app.GRPCServer, error) {
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.ChainUnaryInterceptor(
			middleware.InterceptorLogger(),
			middleware.CredentialInterceptor(),
		),
	)
	opts = append(opts,
		grpc.ChainUnaryInterceptor(
			serverOptions...,
		),
	)

	server, err := app.NewGRPCServer(cfg, opts...)
	if err != nil {
		return nil, errs.NewInternal(err)
	}

	pb.RegisterLibraryServer(server.Ser, &LibraryServer{uc: uc})

	return server, nil
}
