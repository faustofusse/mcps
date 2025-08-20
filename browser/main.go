package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	implementation := &mcp.Implementation{ Name: "mcp-browser", Version: "v0.0.1" }

	server := mcp.NewServer(implementation, &mcp.ServerOptions{})

	mcp.AddTool(server, &ScreenshotTool, Screenshot)

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
	// ctx, cancel := chromedp.NewExecAllocator(ctx, chromedp.WindowSize(1920, 1080))
	// defer cancel()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second * 15)
	defer cancel()

	var image []byte

	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(params.Arguments.Url),
		WaitForRequests(time.Millisecond * 1000),
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

// espera a que haya un cierto tiempo (gap) sin requests
func WaitForRequests(gap time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		timer := time.NewTimer(gap)

		// escucho eventos de red
		chromedp.ListenTarget(ctx, func(event any) {
			switch parsed := event.(type) {
			// cuando termine de cargar una request
			case *network.EventLoadingFinished:
				log.Printf("termino de cargar la request %v", parsed.RequestID)
				// reinicio el timer
				timer.Reset(gap)
			}
		})

		select {
		// si se me cerro el context (capaz timeouteó) me voy con error
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return ctx.Err()
		// si termino el timer, me voy contento
		case <-timer.C:
			return nil
		}
	})
}
