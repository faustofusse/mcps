package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// func SayHello(ctx context.Context, req *mcp.ServerR)

func main() {
	// ctx := context.Background()

	implementation := &mcp.Implementation{ Name: "mcp-compose", Version: "v0.0.1" }

	server := mcp.NewServer(implementation, &mcp.ServerOptions{})

	type HelloParams struct {
		Name string `json:"name" jsonschema:"the name of the person to greet"`
	}
	// type HelloResult struct {
	// 	Text string `json:"name" jsonschema:"greet text"`
	// }
	mcp.AddTool(server, &mcp.Tool{
		Name: "greet",
		Description: "say hello",
	}, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[HelloParams]) (*mcp.CallToolResultFor[struct{}], error) {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Hello " + params.Name,
				},
			},
		}, nil
	})

	// handler := func(r *http.Request) *mcp.Server { return server }
	// http.ListenAndServe(":8080", mcp.NewStreamableHTTPHandler(handler, nil))

	ctx := context.Background()
	transport := mcp.NewStdioTransport()
	// transport := mcp.NewLoggingTransport(&mcp.StdioTransport{}, os.Stderr)
	err := server.Run(ctx, transport)
	if err != nil { panic(err) }
}
