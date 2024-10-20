package service

type Services struct {
	Users    *UsersService
	Links    *LinksService
	Sessions *SessionsService
	Email    *EmailService
}

func NewServices(users *UsersService, links *LinksService, sessions *SessionsService, email *EmailService) *Services {
	return &Services{
		Users:    users,
		Links:    links,
		Sessions: sessions,
		Email:    email,
	}
}
