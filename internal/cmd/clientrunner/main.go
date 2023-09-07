package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/pkg/errors"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
)

func main() {
	accessToken := flag.String("access-token", "", "")
	serviceTokenID := flag.String("service-token-id", "", "")
	serviceToken := flag.String("service-token", "", "")
	flag.Parse()

	debugTpt := DebugRoundTripper(os.Stderr, http.DefaultTransport)
	var tpt http.RoundTripper
	if *accessToken != "" {
		tok := &oauth2.Token{AccessToken: *accessToken}
		tpt = &oauth2.Transport{Base: debugTpt, Source: oauth2.StaticTokenSource(tok)}
	} else if *serviceTokenID != "" && *serviceToken != "" {
		tpt = RoundtripperFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Set("Authorization", *serviceTokenID+":"+*serviceToken)
			return debugTpt.RoundTrip(r)
		})
	}
	cl := planetscale.NewClient(&http.Client{Transport: tpt}, nil)

	ctx := context.Background()

	res200, res403, res404, res500, err := cl.CreateDatabase(ctx, "ps-e2e-prod", planetscale.CreateDatabaseReq{
		Name: "hello_antoine",
	})
	if err != nil {
		slog.Error("failed to get current user", "err", err)
	} else {
		switch {
		case res200 != nil:
			slog.Info("response", "resp", res200)
		case res403 != nil:
			slog.Error("403 error")
		case res404 != nil:
			slog.Error("404 error")
		case res500 != nil:
			slog.Error("500 error")
		}
	}
}

func DebugRoundTripper(out io.Writer, tpt http.RoundTripper) http.RoundTripper {
	return RoundtripperFunc(func(r *http.Request) (*http.Response, error) {
		debugReq, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			return nil, errors.Wrap(err, "dumping request output")
		}
		debugReq = append(debugReq, '\n')
		_, err = out.Write(debugReq)
		if err != nil {
			return nil, errors.Wrap(err, "writing request output to stderr")
		}
		res, err := tpt.RoundTrip(r)
		if res == nil {
			return res, err
		}
		debugRes, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, errors.Wrap(err, "dumping response output")
		}
		debugRes = append(debugRes, '\n')
		_, err = out.Write(debugRes)
		if err != nil {
			return nil, errors.Wrap(err, "writing response output to stderr")
		}
		return res, err
	})
}

type RoundtripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundtripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func ptr[T any](v T) *T {
	return &v
}
