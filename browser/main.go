package main

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	implementation := &mcp.Implementation{ Name: "mcp-browser", Version: "v0.0.1" }

	server := mcp.NewServer(implementation, &mcp.ServerOptions{})

	mcp.AddTool(server, &ScreenshotTool, Screenshot)

	// handler := func(r *http.Request) *mcp.Server { return server }
	// http.ListenAndServe(":8080", mcp.NewStreamableHTTPHandler(handler, nil))

	err := server.Run(ctx, mcp.NewStdioTransport())
	if err != nil { panic(err) }
}

type ScreenshotParams struct {
	Url string `json:"url" jsonschema:"url of the webpage"`
}

var ScreenshotTool = mcp.Tool{
	Name: "screenshot",
	Description: "take screenshot of webpage",
}

type ScreenshotResult struct {}

func Screenshot(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ScreenshotParams]) (*mcp.CallToolResultFor[ScreenshotResult], error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var image []byte

	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(params.Arguments.Url),
		chromedp.FullScreenshot(&image, 90),
	})
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResultFor[ScreenshotResult]{
		Content: []mcp.Content{
			&mcp.ImageContent{
				Data: image,
				MIMEType: "image/jpeg",
			},
		},
	}, nil
}
