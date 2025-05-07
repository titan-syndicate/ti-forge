package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium-plugin-sdk/pkg/logger"
	"github.com/titan-syndicate/titanium-plugin-sdk/pkg/pluginapi"
	"google.golang.org/grpc"

	"github.com/titan-syndicate/ti-scaffold/cmd/scaffold"
)

// GRPCPlugin is the gRPC implementation of the plugin
type GRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
}

// GRPCServer implements the gRPC server for the plugin
func (p *GRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pluginapi.RegisterPluginServer(s, &pluginServer{})
	// logger.Log.Info("Registering gRPC server")
	return nil
}

// GRPCClient implements the gRPC client for the plugin
func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	logger.Log.Info("Creating gRPC client")
	return &pluginClient{client: pluginapi.NewPluginClient(c)}, nil
}

// pluginServer implements the gRPC server interface
type pluginServer struct {
	pluginapi.UnimplementedPluginServer
}

// Name implements the Name RPC method
func (s *pluginServer) Name(ctx context.Context, req *pluginapi.Empty) (*pluginapi.NameResponse, error) {
	logger.Log.Debug("Name called")
	return &pluginapi.NameResponse{Name: "ti-scaffold"}, nil
}

// Version implements the Version RPC method
func (s *pluginServer) Version(ctx context.Context, req *pluginapi.Empty) (*pluginapi.VersionResponse, error) {
	logger.Log.Debug("Version called")
	return &pluginapi.VersionResponse{Version: "1.0.0"}, nil
}

// Execute implements the Execute RPC method
func (s *pluginServer) Execute(ctx context.Context, req *pluginapi.ExecuteRequest) (*pluginapi.ExecuteResponse, error) {
	logger.Log.Infow("Execute called v2", "args", req.Args)

	// Set up os.Args for Cobra
	os.Args = append([]string{"ti-scaffold"}, req.Args...)

	// Execute the Cobra command
	logger.Log.Info("About to execute scaffold command...")
	if err := scaffold.Execute(); err != nil {
		logger.Log.Errorw("Error executing plugin", "error", err)
		// Ensure logs are flushed before returning error
		logger.Sync()
		return &pluginapi.ExecuteResponse{
			Result: fmt.Sprintf("Error executing plugin: %v", err),
		}, nil
	}

	logger.Log.Info("Scaffold command completed successfully")
	// Ensure logs are flushed before returning success
	logger.Sync()
	return &pluginapi.ExecuteResponse{
		Result: "Plugin executed successfully",
	}, nil
}

// pluginClient implements the PluginInterface for the gRPC client
type pluginClient struct {
	client pluginapi.PluginClient
}

// Name implements the PluginInterface
func (c *pluginClient) Name() string {
	logger.Log.Debug("Client Name called")
	resp, err := c.client.Name(context.Background(), &pluginapi.Empty{})
	if err != nil {
		logger.Log.Errorw("Name error", "error", err)
		return ""
	}
	return resp.Name
}

// Version implements the PluginInterface
func (c *pluginClient) Version() string {
	logger.Log.Debug("Client Version called")
	resp, err := c.client.Version(context.Background(), &pluginapi.Empty{})
	if err != nil {
		logger.Log.Errorw("Version error", "error", err)
		return ""
	}
	return resp.Version
}

// Execute implements the PluginInterface
func (c *pluginClient) Execute(args []string) (string, error) {
	logger.Log.Infow("Client Execute called", "args", args)
	resp, err := c.client.Execute(context.Background(), &pluginapi.ExecuteRequest{
		Args: args,
	})
	if err != nil {
		logger.Log.Errorw("Execute error", "error", err)
		return "", err
	}
	return resp.Result, nil
}

func main() {
	// Initialize logger with default level
	if err := logger.Init(logger.Config{
		Level:      "info",
		PluginName: "ti-scaffold",
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Check if we're running in dev mode
	if len(os.Args) > 1 && os.Args[1] == "--dev" {
		logger.Log.Info("Running in development mode...")
		// Remove the --dev flag from args
		os.Args = append(os.Args[:1], os.Args[2:]...)
		if err := scaffold.Execute(); err != nil {
			logger.Log.Errorw("Error executing command", "error", err)
			os.Exit(1)
		}
		return
	}

	// Create the plugin map
	pluginMap := map[string]plugin.Plugin{
		"plugin": &GRPCPlugin{},
	}

	// Serve the plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "TITANIUM_PLUGIN",
			MagicCookieValue: "titanium",
		},
		Plugins:    pluginMap,
		GRPCServer: plugin.DefaultGRPCServer,
	})

	// This line will never be reached as plugin.Serve blocks
	logger.Log.Info("Starting plugin...")
	os.Exit(0)
}
