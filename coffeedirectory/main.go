package main

import (
	"context"
	"net"
	"os"

	"cloud.google.com/go/datastore"
	pb "github.com/ahmetalpbalkan/coffeelog/coffeelog"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	projectID = "ahmetb-starter" // TODO make configurable
)

var (
	userDirectoryBackend string
	log                  *logrus.Entry
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log = logrus.WithField("service", "coffeedirectory")

	if env := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); env == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}
	userDirectoryBackend = os.Getenv("USER_DIRECTORY_HOST")
	if userDirectoryBackend == "" {
		log.Fatal("USER_DIRECTORY_HOST not set")
	}

	ds, err := datastore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.WithField("error", err).Fatal("failed to create client")
	}
	defer ds.Close()

	addr := "0.0.0.0:8002" // TODO make configurable
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	svc := &service{ds}
	pb.RegisterRoasterDirectoryServer(grpcServer, svc)
	pb.RegisterActivityDirectoryServer(grpcServer, svc)
	log.WithField("addr", addr).WithField("service", "coffeedirectory").Info("starting to listen on grpc")
	log.Fatal(grpcServer.Serve(lis))
}
