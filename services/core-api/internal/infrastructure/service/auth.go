package service

import (
	"context"
	"time"

	"github.com/aligh5331/godrop/services/auth/grpc/auth"
	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
)

type AuthService struct {
	client auth.AuthServiceClient
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*dto.TokenPairs, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.LoginRequest{
		Email:    email,
		Password: password,
		MetaData: nil,
	}

	res, err := s.client.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &dto.TokenPairs{
		User: &dto.User{
			ID:    res.User.Id,
			Email: res.User.Email,
			Name:  res.User.Name,
		},
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	}
	return out, nil
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*dto.TokenPairs, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
		MetaData: nil,
	}
	res, err := s.client.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &dto.TokenPairs{
		User: &dto.User{
			ID:    res.User.Id,
			Email: res.User.Email,
			Name:  res.User.Name,
		},
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	}
	return out, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, newPass, OldPass string) (*dto.TokenPairs, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.ChangePasswordRequest{
		UserId:      userID,
		OldPassword: OldPass,
		NewPassword: newPass,
		MetaData:    nil,
	}
	res, err := s.client.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &dto.TokenPairs{
		User: &dto.User{
			ID:    res.User.Id,
			Email: res.User.Email,
			Name:  res.User.Name,
		},
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	}
	return out, nil
}

func (s *AuthService) UpdateName(ctx context.Context, userID, newName string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.UpdateNameRequest{
		UserId:  userID,
		NewName: newName,
	}
	_, err := s.client.UpdateName(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) UpdateEmail(ctx context.Context, userID, newEmail, pass string) (*dto.TokenPairs, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.UpdateEmailRequest{
		UserId:   userID,
		NewEmail: newEmail,
		Password: pass,
	}
	_, err := s.client.UpdateEmail(ctx, req)
	if err != nil {
		return nil, err
	}
	return s.Login(ctx, newEmail, pass)
}

func (s *AuthService) DeleteUser(ctx context.Context, userID, pass string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.DeleteUserRequest{
		UserId:   userID,
		Password: pass,
	}
	_, err := s.client.DeleteUser(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairs, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.RefreshSessionRequest{
		RefreshToken: refreshToken,
	}
	res, err := s.client.RefreshSession(ctx, req)
	if err != nil {
		return nil, err
	}

	out := &dto.TokenPairs{
		User:         nil,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
	return out, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.LogoutRequest{
		AccessToken: accessToken,
	}
	_, err := s.client.Logout(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) RevokeSession(ctx context.Context, sessionID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req := &auth.RevokeSessionRequest{
		SessionId: sessionID,
	}
	_, err := s.client.RevokeSession(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) RevokeAllSessions(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &auth.RevokeAllSessionsRequest{
		UserId: userID,
	}
	_, err := s.client.RevokeAllSessions(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) CheckSessionIsValid(ctx context.Context, accessToken string) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &auth.CheckSessionRequest{
		AccessToken: accessToken,
	}
	_, err := s.client.CheckSession(ctx, req)
	if err != nil {
		return false
	}
	return true
}

func (s *AuthService) GetUserSessions(ctx context.Context, userID string) ([]*dto.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &auth.GetUserSessionsRequest{
		UserId: userID,
	}
	res, err := s.client.GetUserSessions(ctx, req)
	if err != nil {
		return nil, err
	}
	var sessions []*dto.Session
	for _, session := range res.Sessions {
		sessions = append(sessions, &dto.Session{
			ID:        session.SessionId,
			IP:        session.MetaData.Ip,
			UserAgent: session.MetaData.UserAgent,
		})
	}
	return sessions, nil
}

func (s *AuthService) GetUser(context.Context, string) (*dto.User, error) {
	//TODO: Need to implement in the auth service. it does not exist in grpc right now
	panic("implement me")
}
