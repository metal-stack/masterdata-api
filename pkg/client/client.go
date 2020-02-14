package client

import (
	"github.com/metal-stack/masterdata-api/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"go.uber.org/zap"
)

// Client defines the client API
type Client interface {
	Project() v1.ProjectServiceClient
	Tenant() v1.TenantServiceClient
	Close() error
}

// GRPCClient is a Client implementation with grpc transport.
type GRPCClient struct {
	conn    *grpc.ClientConn
	log     *zap.Logger
	hmacKey string
}

// NewClient creates a new client for the services for the given address, with the certificate and hmac.
func NewClient(address string, certFile string, hmacKey string, logger *zap.Logger) (Client, error) {

	if hmacKey == "" {
		logger.Sugar().Fatal("no hmac-key specified")
	}

	client := GRPCClient{
		log:     logger,
		hmacKey: hmacKey,
	}

	// Set up the credentials for the connection.
	perRPCHMACAuthenticator, err := auth.NewHMACAuther(logger, hmacKey, auth.EditUser)
	if err != nil {
		logger.Sugar().Fatalf("failed to create hmac-authenticator: %v", err)
	}
	// TODO serverNameOverride should only be there for tests...?
	creds, err := credentials.NewClientTLSFromFile(certFile, "metal-stack.io")
	if err != nil {
		logger.Sugar().Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPCHMACAuthenticator),
		// oauth.NewOauthAccess requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(creds),

		// grpc.WithInsecure(),
		grpc.WithBlock(),
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logger.Sugar().Errorf("did not connect: %v", err)
		return nil, err
	}
	client.conn = conn

	return client, nil
}

// Close the underlying connection
func (c GRPCClient) Close() error {
	return c.conn.Close()
}

// Project is the root accessor for project related functions
func (c GRPCClient) Project() v1.ProjectServiceClient {
	return v1.NewProjectServiceClient(c.conn)
}

// Tenant is the root accessor for tenant related functions
func (c GRPCClient) Tenant() v1.TenantServiceClient {
	return v1.NewTenantServiceClient(c.conn)
}
