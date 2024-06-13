package service

type SessionService interface {
	GenerateToken() (string, error)
	ValidateToken(token string) (bool, error)
}

type sessionService struct {
}

func NewSessionService() SessionService {
	return &sessionService{}
}

func (s *sessionService) GenerateToken() (string, error) {
	return "", nil
}

func (s *sessionService) ValidateToken(token string) (bool, error) {
	return true, nil
}
