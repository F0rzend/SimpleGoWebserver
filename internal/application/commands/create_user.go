package commands

import (
	"time"

	"github.com/F0rzend/SimpleGoWebserver/internal/domain"
)

type CreateUserCommand struct {
	Name     string
	Username string
	Email    string
}

type CreateUserCommandHandler struct {
	getId          func() uint64
	userRepository domain.UserRepository
}

func NewCreateUserCommand(userRepository domain.UserRepository) *CreateUserCommandHandler {
	if userRepository == nil {
		panic("userRepository is nil")
	}

	return &CreateUserCommandHandler{
		getId:          userIDGenerator(),
		userRepository: userRepository,
	}
}

func userIDGenerator() func() uint64 {
	var id uint64 = 0
	return func() uint64 {
		id++
		return id
	}
}

func (h *CreateUserCommandHandler) Handle(cmd CreateUserCommand) (uint64, error) {
	user, err := domain.NewUser(
		h.getId(),
		cmd.Name,
		cmd.Username,
		cmd.Email,
		0,
		0,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	if err := h.userRepository.Create(user); err != nil {
		return 0, err
	}

	return user.ID, nil
}
