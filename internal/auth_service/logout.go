package auth_service

import (
	"fmt"
	"github.com/Turalchik/authentication-service/internal/apperrors"
)

func (authService *AuthService) Logout(accessToken string, userID string) error {
	// заносим access токен в black-list
	// TODO
	// реализовать black-list с помощью redis
	fmt.Println(accessToken)

	// удаляем refresh токен из базы
	if err := authService.repo.DeleteSessionByUserID(userID); err != nil {
		return apperrors.ErrCantDeleteSession
	}

	return nil
}
