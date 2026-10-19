package profilemanager

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"core/internal/domain"
)

type UserStore interface {
	Create(ctx context.Context, input domain.CreateProfileInput) (domain.Profile, error)
	GetByID(ctx context.Context, id int64) (domain.Profile, error)
	Update(ctx context.Context, id int64, input domain.UpdateProfileInput) (domain.Profile, error)
}

type ProfileManager struct {
	users UserStore
}

func NewProfileManager(users UserStore) *ProfileManager {
	return &ProfileManager{users: users}
}

func (m *ProfileManager) Create(ctx context.Context, input domain.CreateProfileInput) (domain.Profile, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Bio = strings.TrimSpace(input.Bio)

	if err := validateEmail(input.Email); err != nil {
		return domain.Profile{}, err
	}
	if err := validateDisplayName(input.DisplayName); err != nil {
		return domain.Profile{}, err
	}
	if err := validateBio(input.Bio); err != nil {
		return domain.Profile{}, err
	}

	return m.users.Create(ctx, input)
}

func (m *ProfileManager) Get(ctx context.Context, id int64) (domain.Profile, error) {
	if id <= 0 {
		return domain.Profile{}, fmt.Errorf("%w: identyfikator musi być większy od zera", domain.ErrInvalidInput)
	}

	return m.users.GetByID(ctx, id)
}

func (m *ProfileManager) Update(ctx context.Context, id int64, input domain.UpdateProfileInput) (domain.Profile, error) {
	if id <= 0 {
		return domain.Profile{}, fmt.Errorf("%w: identyfikator musi być większy od zera", domain.ErrInvalidInput)
	}
	if input.DisplayName == nil && input.Bio == nil {
		return domain.Profile{}, fmt.Errorf("%w: podaj przynajmniej jedno pole do aktualizacji", domain.ErrInvalidInput)
	}

	if input.DisplayName != nil {
		displayName := strings.TrimSpace(*input.DisplayName)
		if err := validateDisplayName(displayName); err != nil {
			return domain.Profile{}, err
		}
		input.DisplayName = &displayName
	}

	if input.Bio != nil {
		bio := strings.TrimSpace(*input.Bio)
		if err := validateBio(bio); err != nil {
			return domain.Profile{}, err
		}
		input.Bio = &bio
	}

	return m.users.Update(ctx, id, input)
}

func validateEmail(email string) error {
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || !strings.Contains(email, "@") {
		return fmt.Errorf("%w: podaj poprawny adres e-mail", domain.ErrInvalidInput)
	}
	return nil
}

func validateDisplayName(displayName string) error {
	length := utf8.RuneCountInString(displayName)
	if length < 2 || length > 80 {
		return fmt.Errorf("%w: nazwa wyświetlana musi mieć od 2 do 80 znaków", domain.ErrInvalidInput)
	}
	return nil
}

func validateBio(bio string) error {
	if utf8.RuneCountInString(bio) > 500 {
		return fmt.Errorf("%w: opis może mieć maksymalnie 500 znaków", domain.ErrInvalidInput)
	}
	return nil
}
