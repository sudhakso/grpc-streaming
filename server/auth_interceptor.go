package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager, accessibleRoles: accessibleRoles}
}

// Unary Interceptor
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	// Interceptor : unary
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Just print...
		log.Println("--> unary interceptor: ", info.FullMethod)
		// Check if Authroize
		if err := interceptor.Authorize(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// Stream Interceptor
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	// Interceptor :  stream
	return func(srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		// Just print...
		log.Println("--> stream interceptor: ", info.FullMethod)
		// Check if Authorized
		err := interceptor.Authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

// Authorize
func (interceptpr *AuthInterceptor) Authorize(context context.Context, method string) error {
	// Check if Method is resticted
	accessibleRoles, ok := interceptpr.accessibleRoles[method]
	if !ok {
		//everyone has access
		return nil
	}
	// If Restricted, get the JWT from context, and get the Users Role.
	md, ok := metadata.FromIncomingContext(context)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "Unauthorized, token not provided")
	}
	// Get the JWT from context
	// Check in accessible roles if the method is allowed by a User role.
	accessToken := values[0]
	claims, err := interceptpr.jwtManager.verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Unauthorized, invalid token ")
	}
	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}
	return status.Error(codes.PermissionDenied, "user doesn't have permissions to access the resource")
}
