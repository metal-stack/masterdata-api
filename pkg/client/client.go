package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/metal-stack/masterdata-api/pkg/auth"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"

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
func NewClient(ctx context.Context, hostname string, port int, certFile string, keyFile string, caFile string, hmacKey string, logger *zap.Logger) (Client, error) {

	address := fmt.Sprintf("%s:%d", hostname, port)

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system credentials: %w", err)
	}

	if caFile != "" {
		ca, err := ioutil.ReadFile(caFile)
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
		ServerName:   hostname,
		Certificates: []tls.Certificate{clientCertificate},
		RootCAs:      certPool,
	})

	if hmacKey == "" {
		return nil, errors.New("no hmac-key specified")
	}

	client := GRPCClient{
		log:     logger,
		hmacKey: hmacKey,
	}

	// Set up the credentials for the connection.
	perRPCHMACAuthenticator, err := auth.NewHMACAuther(logger, hmacKey, auth.EditUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create hmac-authenticator")
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
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

// Tenant is the root accessor for tenant related functions
func (c GRPCClient) Tenant() v1.TenantServiceClient {
	return v1.NewTenantServiceClient(c.conn)
}
