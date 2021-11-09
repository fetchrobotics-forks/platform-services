package main

// This client can be used to access a fully deployed service mesh example using both
// a LetsEcrypt certificate to secure the connection and a JWT token to authenticate
// on top of the secured connection

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fetchcore-forks/platform-services/internal/platform"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-stack/stack"
	"github.com/karlmutch/errors"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	experimentsrv "github.com/fetchcore-forks/platform-services/internal/gen/experimentsrv"
)

var (
	logger = platform.NewLogger("cli-experiment")

	serverAddr = flag.String("server-addr", "127.0.0.1:443", "The server address in the format of host:port")
	optToken   = flag.String("auth0-token", "", "The authentication token to be used for the user, (JWT Token)")
	caCertOpt  = flag.String("ca-cert", "", "The certificate for the certificate authority being used by the gRPC server")
	serverOpt  = flag.String("server-name", "", "The name of the server for which the certificate on the remote end will be verified for (IP SANS)")
)

func getToken() (token *oauth2.Token, err errors.Error) {
	return &oauth2.Token{
		AccessToken: *optToken,
	}, nil
}

func loadTLSCreds(caCert string, serverName string) (creds credentials.TransportCredentials, err errors.Error) {
	if len(caCert) == 0 {
		return nil, errors.New("A CA certificate must be supplied").With("stack", stack.Trace().TrimRuntime())
	}

	// Load certificate of the CA who signed server's certificate
	pemServerCA, errGo := ioutil.ReadFile(caCert)
	if errGo != nil {
		return nil, errors.Wrap(errGo).With("fn", caCert).With("stack", stack.Trace().TrimRuntime())
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, errors.New("failed to add server CA's certificate").With("fn", caCert).With("stack", stack.Trace().TrimRuntime())
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}
	if len(serverName) != 0 {
		config.ServerName = serverName
	}

	return credentials.NewTLS(config), nil
}

func main() {
	flag.Parse()

	// Obtain a JWT and prepare it for use with the gRPC connection and API calls
	token, err := getToken()
	if err != nil {
		logger.Fatal(fmt.Sprint(err.With("address", *serverAddr).With("stack", stack.Trace().TrimRuntime())))
	}
	perRPC := oauth.NewOauthAccess(token)

	cert, err := loadTLSCreds(*caCertOpt, *serverOpt)
	if err != nil {
		log.Fatal(err.Error())
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(cert),
	}

	logger.Info("", "server-addr", *serverAddr, "stack", stack.Trace().TrimRuntime())

	conn, errGo := grpc.Dial(*serverAddr, opts...)
	if errGo != nil {
		logger.Fatal(fmt.Sprint(errors.Wrap(errGo).With("address", *serverAddr).With("stack", stack.Trace().TrimRuntime())))
	}
	defer conn.Close()
	client := experimentsrv.NewExperimentsClient(conn)

	logger.Info("stack", stack.Trace().TrimRuntime())

	response, errGo := client.MeshCheck(context.Background(), &experimentsrv.CheckRequest{Live: true})
	if errGo != nil {
		logger.Error(fmt.Sprint(errors.Wrap(errGo).With("address", *serverAddr).With("stack", stack.Trace().TrimRuntime())))
		os.Exit(-1)
	}
	spew.Fdump(os.Stdout, response)
}
