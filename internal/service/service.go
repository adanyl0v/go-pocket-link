package service

type Services struct {
	Auth     *AuthService
	Users    *UsersService
	Sessions *SessionsService
}
