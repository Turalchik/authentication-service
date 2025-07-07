package main

import (
	"errors"
	"fmt"
	"github.com/Turalchik/authentication-service/internal/apperrors"
)

func main() {
	fmt.Println(errors.Is(fmt.Errorf("can't get claims from access token: %w %w", apperrors.ErrCantCreateTokens, apperrors.ErrInvalidUserID), apperrors.ErrCantCreateTokens))
}
