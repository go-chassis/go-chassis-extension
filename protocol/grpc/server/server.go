package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/util/iputil"
	"github.com/go-chassis/openlog"
)

//err define
var (
	ErrGRPCSvcDescMissing = errors.New("must use server.WithRPCServiceDesc to set desc")
	ErrGRPCSvcType        = errors.New("must set *grpc.ServiceDesc")
)

//const
const (
	Name = "grpc"
)

//Server is grpc server holder
type Server struct {
	s    *grpc.Server
	opts server.Options
}

//Request2Invocation convert grpc protocol to invocation
func Request2Invocation(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) *invocation.Invocation {
	md, _ := metadata.FromIncomingContext(ctx)
	sourceServices := md.Get(common.HeaderSourceName)
	var sourceService string
	if len(sourceServices) >= 1 {
		sourceService = sourceServices[0]
	}
	//TODO maybe need set headers
	m := make(map[string]string)
	inv := &invocation.Invocation{
		MicroServiceName:   runtime.ServiceName,
		SourceMicroService: sourceService,
		Args:               req,
		Protocol:           "grpc",
		SchemaID:           info.FullMethod,
		OperationID:        info.FullMethod,
		Ctx:                context.WithValue(ctx, common.ContextHeaderKey{}, m),
	}
	return inv
}

//New create grpc server
func New(opts server.Options) server.ProtocolServer {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handle grpc.UnaryHandler) (resp interface{}, err error) {
		c, err := handler.GetChain(common.Provider, opts.ChainName)
		if err != nil {
			openlog.Error(fmt.Sprintf("Handler chain init err [%s]", err.Error()))
			return nil, err
		}
		inv := Request2Invocation(ctx, req, info)
		var r *invocation.Response
		c.Next(inv, func(ir *invocation.Response) {
			ir.Result, ir.Err = handle(ctx, req)
			r = ir
		})
		return r.Result, r.Err
	}

	sops := make([]grpc.ServerOption, 0, 2)
	sops = append(sops, grpc.UnaryInterceptor(interceptor))

	if opts.BodyLimit != 0 {
		sops = append(sops, grpc.MaxRecvMsgSize(int(opts.BodyLimit)))
	}
	return &Server{
		opts: opts,
		s:    grpc.NewServer(sops...),
	}
}

//Register register grpc services
func (s *Server) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	opts := server.RegisterOptions{}
	for _, o := range options {
		o(&opts)
	}
	if opts.RPCSvcDesc == nil {
		return "", ErrGRPCSvcDescMissing
	}
	desc, ok := opts.RPCSvcDesc.(*grpc.ServiceDesc)
	if !ok {
		return "", ErrGRPCSvcType
	}
	s.s.RegisterService(desc, schema)
	return "", nil
}

//Start launch the server
func (s *Server) Start() error {
	listener, host, port, lisErr := iputil.StartListener(s.opts.Address, s.opts.TLSConfig)
	if lisErr != nil {
		openlog.Error("listening failed, reason:" + lisErr.Error())
		return lisErr
	}

	registry.InstanceEndpoints[Name] = net.JoinHostPort(host, port)

	go func() {
		if err := s.s.Serve(listener); err != nil {
			server.ErrRuntime <- err
		}
	}()
	return nil
}

//Stop gracfully shutdown grpc server
func (s *Server) Stop() error {
	s.s.GracefulStop()
	return nil
}

//String return server name
func (s *Server) String() string {
	return Name
}
func init() {
	server.InstallPlugin(Name, New)
}
