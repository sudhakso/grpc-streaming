package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

// goal: get access token and refresh token if expired
func NewAuthInterceptor(
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshTokenDuration time.Duration,
) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}
	err := interceptor.scheduleRefreshToken(refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return interceptor, nil
}

// return a unary interceptor context with access token
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		log.Printf("--> unary interceptor: %s", method)

		// check if method needs Authentication context
		if interceptor.authMethods[method] {
			// get Context with token
			c := interceptor.attachToken(ctx)
			return invoker(c, method, req, reply, cc, opts...)
		}
		// method doesnt need Authentication
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// returns a stream interceptor
func (interceptor *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("---> stream interceptor: %s", method)

		// check if method needs Authentication context
		if interceptor.authMethods[method] {
			c := interceptor.attachToken(ctx)
			return streamer(c, desc, cc, method, opts...)
		}
		// If method doesn't need auth context
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

func (interceptor *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := interceptor.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := interceptor.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}

func (interceptor *AuthInterceptor) refreshToken() error {
	accessToken, err := interceptor.authClient.Login()
	if err != nil {
		return err
	}

	// Access token
	interceptor.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)

	return nil
}
