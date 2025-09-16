package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/auth"
)

// Client defines the client API
type Client interface {
	Project() v1.ProjectServiceClient
	ProjectMember() v1.ProjectMemberServiceClient
	Tenant() v1.TenantServiceClient
	TenantMember() v1.TenantMemberServiceClient
	Version() v1.VersionServiceClient
	Close() error
}

// GRPCClient is a Client implementation with grpc transport.
type GRPCClient struct {
	conn    *grpc.ClientConn
	log     *slog.Logger
	hmacKey string
}

type Config struct {
	Logger *slog.Logger

	Hostname string
	Port     uint

	CertFile string
	KeyFile  string
	CaFile   string
	Insecure bool

	HmacKey string

	// Namespace if set adds this namespace to namespaced requests such that it does not need to be passed all the time
	Namespace string
}

func (c *Config) validate() error {
	if c == nil {
		return fmt.Errorf("config must not be nil")
	}

	if c.Hostname == "" {
		return fmt.Errorf("hostname must be specified")
	}

	if c.KeyFile != "" || c.CertFile != "" {
		if c.KeyFile == "" || c.CertFile == "" {
			return fmt.Errorf("either both key and cert file must be specified or none of them")
		}
	}

	return nil
}

// NewClient creates a new client for the services for the given address, with the certificate and hmac.
func NewClient(config *Config) (Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	address := fmt.Sprintf("%s:%d", config.Hostname, config.Port)

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system credentials: %w", err)
	}

	if caFile := config.CaFile; caFile != "" {
		ca, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %w", err)
		}

		ok := certPool.AppendCertsFromPEM(ca)
		if !ok {
			return nil, fmt.Errorf("failed to append ca certs: %s", caFile)
		}
	}

	var (
		certificates []tls.Certificate
		opts         []grpc.DialOption
	)

	if config.CertFile != "" && config.KeyFile != "" {
		clientCertificate, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %w", err)
		}

		certificates = append(certificates, clientCertificate)

		creds := credentials.NewTLS(&tls.Config{
			ServerName:         config.Hostname,
			Certificates:       certificates,
			RootCAs:            certPool,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: config.Insecure, // nolint:gosec
		})

		opts = append(opts,
			// oauth.NewOauthAccess requires the configuration of transport
			// credentials.
			grpc.WithTransportCredentials(creds),
		)
	} else {
		opts = append(opts,
			grpc.WithTransportCredentials(grpcinsecure.NewCredentials()),
		)
	}

	client := GRPCClient{
		log: config.Logger,
	}

	if config.HmacKey != "" {
		client.hmacKey = config.HmacKey

		// Set up the credentials for the connection.
		perRPCHMACAuthenticator, err := auth.NewHMACAuther(config.HmacKey, auth.EditUser)
		if err != nil {
			return nil, fmt.Errorf("failed to create hmac-authenticator: %w", err)
		}

		opts = append(opts,
			// In addition to the following grpc.DialOption, callers may also use
			// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
			// itself.
			// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
			grpc.WithPerRPCCredentials(perRPCHMACAuthenticator))
	}

	if config.Namespace != "" {
		opts = append(opts, NamespaceInterceptor(config.Namespace))
	}

	// Set up a connection to the server.
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}

	client.conn = conn

	return client, nil
}

func NamespaceInterceptor(namespace string) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		switch r := req.(type) {
		case *v1.TenantMemberCreateRequest:
			if r.TenantMember.Namespace == "" {
				r.TenantMember.Namespace = namespace
			}
		case *v1.ProjectMemberCreateRequest:
			if r.ProjectMember.Namespace == "" {
				r.ProjectMember.Namespace = namespace
			}
		case *v1.TenantMemberFindRequest:
			if r.Namespace == "" {
				r.Namespace = namespace
			}
		case *v1.ProjectMemberFindRequest:
			if r.Namespace == "" {
				r.Namespace = namespace
			}
		case *v1.FindParticipatingProjectsRequest:
			if r.Namespace == "" {
				r.Namespace = namespace
			}
		case *v1.FindParticipatingTenantsRequest:
			if r.Namespace == "" {
				r.Namespace = namespace
			}
		case *v1.ListTenantMembersRequest:
			if r.Namespace == "" {
				r.Namespace = namespace
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	})
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

// Tenant is the root accessor for tenant related functions
func (c GRPCClient) TenantMember() v1.TenantMemberServiceClient {
	return v1.NewTenantMemberServiceClient(c.conn)
}

func (c GRPCClient) Version() v1.VersionServiceClient {
	return v1.NewVersionServiceClient(c.conn)
}
