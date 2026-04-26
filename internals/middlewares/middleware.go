package middleware

type AuthServiceInterface interface {
	GetUserIDFromSession(sid string) (int64, error)
}
