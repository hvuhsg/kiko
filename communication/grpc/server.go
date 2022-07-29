package grpc

import (
	context "context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/hvuhsg/kiko/execution"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	UnimplementedKikoServiceServer
	engine     execution.IEngine
	grpcServer grpc.Server
	host       string
	port       int
}

func NewServer(engine *execution.IEngine, host string, port int) *server {
	grpcServer := grpc.NewServer()
	s := new(server)
	s.engine = *engine
	s.grpcServer = *grpcServer
	s.host = host
	s.port = port

	RegisterKikoServiceServer(grpcServer, s)

	return s
}

func (s *server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) AddNode(ctx context.Context, in *emptypb.Empty) (*AddNodeResponse, error) {
	nodeUuid := s.engine.AddNode()
	return &AddNodeResponse{Node: nodeUuid.String()}, nil
}

func (s *server) DeleteNode(ctx context.Context, in *DeleteNodeRequest) (*emptypb.Empty, error) {
	nodeUuidString := in.GetNode()

	nodeUuid, err := uuid.Parse(nodeUuidString)
	if err != nil {
		return nil, err
	}

	err = s.engine.RemoveNode(nodeUuid)
	return nil, err
}

func (s *server) CreateConnection(ctx context.Context, in *CreateConnectionRequest) (*emptypb.Empty, error) {
	fromString := in.GetFrom()
	toString := in.GetTo()

	fromUuid, err := uuid.Parse(fromString)
	if err != nil {
		return nil, err
	}

	toUuid, err := uuid.Parse(toString)
	if err != nil {
		return nil, err
	}

	err = s.engine.AddConnection(fromUuid, toUuid, uint(in.GetWeight()))

	return nil, err
}

func (s *server) DeleteConnection(ctx context.Context, in *DeleteConnectionRequest) (*emptypb.Empty, error) {
	fromString := in.GetFrom()
	toString := in.GetTo()

	fromUuid, err := uuid.Parse(fromString)
	if err != nil {
		return nil, err
	}

	toUuid, err := uuid.Parse(toString)
	if err != nil {
		return nil, err
	}

	err = s.engine.RemoveConnection(fromUuid, toUuid)

	return nil, err
}

func (s *server) GetKRecommendations(ctx context.Context, in *GetKRecommendationsRequest) (*GetKRecommendationsReponse, error) {
	nodeUuidString := in.GetNode()
	nodeUuid, err := uuid.Parse(nodeUuidString)
	if err != nil {
		return nil, err
	}

	recommendations, err := s.engine.Recommend(nodeUuid, int(in.GetK()))
	if err != nil {
		return nil, err
	}

	recommendationsString := make([]string, len(recommendations))

	for i, recommendation := range recommendations {
		recommendationsString[i] = recommendation.String()
	}

	return &GetKRecommendationsReponse{Recommendations: recommendationsString}, nil
}

func (s *server) GetNode(ctx context.Context, in *GetNodeRequest) (*GetNodeResponse, error) {
	nodeUuidString := in.GetNode()
	nodeUuid, err := uuid.Parse(nodeUuidString)
	if err != nil {
		return nil, err
	}

	connections, err := s.engine.GetNodeConnections(nodeUuid)
	if err != nil {
		return nil, err
	}

	connectionsMap := make(map[string]int64, len(connections))

	for connection, weight := range connections {
		connectionsMap[connection.String()] = int64(weight)
	}

	return &GetNodeResponse{Connections: connectionsMap}, nil
}
