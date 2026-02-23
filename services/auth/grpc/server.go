package grpc

import (
	"context"

	"github.com/aligh5331/godrop/services/auth/grpc/auth"
	"github.com/aligh5331/godrop/services/auth/internal/application/dto"
	"github.com/aligh5331/godrop/services/auth/internal/application/repository"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	authUc    repository.AuthUseCase
	sessionUc repository.SessionUseCase
}

func (a *AuthServer) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	input := dto.LoginInputDTO{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}
	meta := dto.SessionMetadataDTO{
		IP:          request.GetMetaData().GetIp(),
		ClientAgent: request.GetMetaData().GetUserAgent(),
	}
	loginDTO, err := a.authUc.Login(ctx, input, meta)
	if err != nil {
		return nil, err
	}

	out := &auth.LoginResponse{
		Tokens: &auth.TokenPair{
			AccessToken:  loginDTO.Tokens.AccessToken,
			RefreshToken: loginDTO.Tokens.RefreshToken,
		},
		User: &auth.User{
			Id:    loginDTO.User.ID,
			Name:  loginDTO.User.Name,
			Email: loginDTO.User.Email,
		},
	}
	return out, nil
}

func (a *AuthServer) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	input := dto.RegisterInputDTO{
		Name:     request.GetName(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}
	meta := dto.SessionMetadataDTO{
		IP:          request.GetMetaData().GetIp(),
		ClientAgent: request.GetMetaData().GetUserAgent(),
	}

	registerDTO, err := a.authUc.Register(ctx, input, meta)
	if err != nil {
		return nil, err
	}

	out := &auth.RegisterResponse{
		Tokens: &auth.TokenPair{
			AccessToken:  registerDTO.Tokens.AccessToken,
			RefreshToken: registerDTO.Tokens.RefreshToken,
		},
		User: &auth.User{
			Id:    registerDTO.User.ID,
			Name:  registerDTO.User.Name,
			Email: registerDTO.User.Email,
		},
	}

	return out, nil
}

func (a *AuthServer) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.LoginResponse, error) {
	input := dto.ChangePasswordDTO{
		UserID:  request.GetUserId(),
		OldPass: request.GetOldPassword(),
		NewPass: request.GetNewPassword(),
	}
	meta := dto.SessionMetadataDTO{
		IP:          request.GetMetaData().GetIp(),
		ClientAgent: request.GetMetaData().GetUserAgent(),
	}

	loginDTO, err := a.authUc.ChangePassword(ctx, input, meta)
	if err != nil {
		return nil, err
	}

	out := &auth.LoginResponse{
		Tokens: &auth.TokenPair{
			AccessToken:  loginDTO.Tokens.AccessToken,
			RefreshToken: loginDTO.Tokens.RefreshToken,
		},
		User: &auth.User{
			Id:    loginDTO.User.ID,
			Name:  loginDTO.User.Name,
			Email: loginDTO.User.Email,
		},
	}

	return out, nil
}

func (a *AuthServer) UpdateName(ctx context.Context, request *auth.UpdateNameRequest) (*auth.UpdateNameResponse, error) {
	input := dto.UpdateUserNameDTO{
		UserID: request.GetUserId(),
		Name:   request.GetNewName(),
	}

	if err := a.authUc.UpdateName(ctx, input); err != nil {
		return nil, err
	}
	return &auth.UpdateNameResponse{}, nil
}

func (a *AuthServer) UpdateEmail(ctx context.Context, request *auth.UpdateEmailRequest) (*auth.UpdateEmailResponse, error) {
	input := dto.UpdateUserEmailDTO{
		UserID: request.GetUserId(),
		Email:  request.GetNewEmail(),
		Pass:   request.GetPassword(),
	}

	if err := a.authUc.UpdateEmail(ctx, input); err != nil {
		return nil, err
	}
	return &auth.UpdateEmailResponse{}, nil
}

func (a *AuthServer) DeleteUser(ctx context.Context, request *auth.DeleteUserRequest) (*auth.DeleteUserResponse, error) {
	input := dto.DeleteUserDTO{
		UserID: request.GetUserId(),
		Pass:   request.GetPassword(),
	}

	if err := a.authUc.Delete(ctx, input); err != nil {
		return nil, err
	}
	return &auth.DeleteUserResponse{}, nil
}

func (a *AuthServer) RefreshSession(ctx context.Context, request *auth.RefreshSessionRequest) (*auth.TokenPair, error) {
	tokenPair, err := a.sessionUc.RefreshSession(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	out := &auth.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
	return out, nil
}

func (a *AuthServer) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	err := a.sessionUc.RevokeSessionWithAccessToken(ctx, request.GetAccessToken())
	if err != nil {
		return nil, err
	}
	return &auth.LogoutResponse{}, nil
}

func (a *AuthServer) RevokeSession(ctx context.Context, request *auth.RevokeSessionRequest) (*auth.RevokeSessionResponse, error) {
	err := a.sessionUc.RevokeSession(ctx, request.GetSessionId())
	if err != nil {
		return nil, err
	}
	return &auth.RevokeSessionResponse{}, nil
}

func (a *AuthServer) RevokeAllSessions(ctx context.Context, request *auth.RevokeAllSessionsRequest) (*auth.RevokeAllSessionsResponse, error) {
	err := a.sessionUc.RevokeAllUserSessions(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}
	return &auth.RevokeAllSessionsResponse{}, nil
}

func (a *AuthServer) CheckSession(ctx context.Context, request *auth.CheckSessionRequest) (*auth.CheckSessionResponse, error) {
	err := a.sessionUc.EnsureAccessTokenValid(ctx, request.GetAccessToken())
	if err != nil {
		return nil, err
	}
	return &auth.CheckSessionResponse{}, nil
}

func (a *AuthServer) GetUserSessions(ctx context.Context, request *auth.GetUserSessionsRequest) (*auth.GetUserSessionsResponse, error) {
	sessions, err := a.sessionUc.GetUserSessions(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	var out []*auth.Session
	for _, session := range sessions {
		meta := &auth.MetaData{
			Ip:        session.Metadata.IP,
			UserAgent: session.Metadata.ClientAgent,
		}
		out = append(out, &auth.Session{
			SessionId: session.SessionId,
			MetaData:  meta,
		})
	}
	return &auth.GetUserSessionsResponse{Sessions: out}, nil
}

func NewAuthServer(authUc repository.AuthUseCase, sessionUc repository.SessionUseCase) auth.AuthServiceServer {
	return &AuthServer{
		authUc:    authUc,
		sessionUc: sessionUc,
	}
}
