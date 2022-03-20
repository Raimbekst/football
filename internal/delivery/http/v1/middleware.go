package v1

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userType            = "userType"
	adminCtx            = "adminId"
)

func (h *Handler) userIdentity(header string) (int, string, error) {
	var (
		userId   int
		id       string
		userType string
		err      error
	)
	if header != "" {
		headerParts := strings.Split(header, " ")

		id, userType, err = h.tokenManager.Parse(headerParts[1])

		if err != nil {
			return 0, "", fmt.Errorf("middleware.userIdentity: %w", err)
		}
		userId, err = strconv.Atoi(id)

		if err != nil {
			return 0, "", fmt.Errorf("middleware.userIdentity: %w", err)
		}

	}
	return userId, userType, nil

}
