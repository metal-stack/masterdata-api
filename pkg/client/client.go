package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/metal-stack/masterdata-api/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
)

// Client defines the client API
type Client interface {
	Project() v1.ProjectServiceClient
	ProjectMember() v1.ProjectMemberServiceClient
	Tenant() v1.TenantServiceClient
	Close() error
}

// GRPCClient is a Client implementation with grpc transport.
type GRPCClient struct {
	conn    *grpc.ClientConn
	log     *slog.Logger
	hmacKey string
}

// NewClient creates a new client for the services for the given address, with the certificate and hmac.
func NewClient(ctx context.Context, hostname string, port int, certFile string, keyFile string, caFile string, hmacKey string, insecure bool, logger *slog.Logger) (Client, error) {

	address := fmt.Sprintf("%s:%d", hostname, port)

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system credentials: %w", err)
	}

	if caFile != "" {
		ca, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %w", err)
		}

		ok := certPool.AppendCertsFromPEM(ca)
		if !ok {
			return nil, fmt.Errorf("failed to append ca certs: %s", caFile)
		}
	}

	clientCertificate, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load client key pair: %w", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:         hostname,
		Certificates:       []tls.Certificate{clientCertificate},
		RootCAs:            certPool,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: insecure, // nolint:gosec
	})

	if hmacKey == "" {
		return nil, errors.New("no hmac-key specified")
	}

	client := GRPCClient{
		log:     logger,
		hmacKey: hmacKey,
	}

	// Set up the credentials for the connection.
	perRPCHMACAuthenticator, err := auth.NewHMACAuther(hmacKey, auth.EditUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create hmac-authenticator: %w", err)
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
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
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

// ProjectMember is the root accessor for project member related functions
func (c GRPCClient) ProjectMember() v1.ProjectMemberServiceClient {
	return v1.NewProjectMemberServiceClient(c.conn)
}

// Tenant is the root accessor for tenant related functions
func (c GRPCClient) Tenant() v1.TenantServiceClient {
	return v1.NewTenantServiceClient(c.conn)
}
