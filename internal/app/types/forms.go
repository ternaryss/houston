package types

import (
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
)

type SignInForm struct {
	Email    string
	Password string
	Errors   map[string]FieldError
}

func NewSignInForm(request *http.Request) (SignInForm, error) {
	if request == nil {
		return SignInForm{}, nil
	}

	if err := request.ParseForm(); err != nil {
		return SignInForm{}, err
	}

	return SignInForm{
		Email:    request.FormValue("email"),
		Password: request.FormValue("password"),
		Errors:   make(map[string]FieldError),
	}, nil
}

func (f SignInForm) Validate() {
	if f.Email == "" || f.Password == "" {
		f.Errors["password"] = NewFieldError("password", "Invalid address e-mail or password.")
		return
	}

	if _, err := mail.ParseAddress(f.Email); err != nil {
		f.Errors["password"] = NewFieldError("password", "Invalid address e-mail or password.")
		return
	}
}

type SignUpFrom struct {
	Email          string
	Password       string
	RepeatPassword string
	Errors         map[string]FieldError
}

func NewSignUpForm(request *http.Request) (SignUpFrom, error) {
	if request == nil {
		return SignUpFrom{}, nil
	}

	if err := request.ParseForm(); err != nil {
		return SignUpFrom{}, err
	}

	return SignUpFrom{
		Email:          request.FormValue("email"),
		Password:       request.FormValue("password"),
		RepeatPassword: request.FormValue("repeatPassword"),
		Errors:         make(map[string]FieldError),
	}, nil
}

func (f SignUpFrom) Validate() {
	if f.Email == "" {
		f.Errors["email"] = NewFieldError("email", "Address e-mail is required.")
	} else {
		if _, err := mail.ParseAddress(f.Email); err != nil {
			f.Errors["email"] = NewFieldError("email", "Invalid address e-mail.")
		}
	}

	if f.Password == "" {
		f.Errors["password"] = NewFieldError("password", "Password is required.")
	} else {
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(f.Password)
		hasDigit := regexp.MustCompile(`\d`).MatchString(f.Password)
		hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(f.Password)

		if len(f.Password) < 8 || !hasUpper || !hasDigit || !hasSpecial {
			f.Errors["password"] = NewFieldError("password", "Invalid password (8 characters, one capital, one digit, one special character).")
		} else {
			if f.Password != f.RepeatPassword {
				f.Errors["password"] = NewFieldError("password", "Passwords not match.")
			}
		}
	}
}

type WebAppForm struct {
	Id     string
	Name   string
	Url    string
	Errors map[string]FieldError
}

func NewWebAppForm(request *http.Request) (WebAppForm, error) {
	if request == nil {
		return WebAppForm{}, nil
	}

	if err := request.ParseForm(); err != nil {
		return WebAppForm{}, err
	}

	return WebAppForm{
		Name:   request.FormValue("name"),
		Url:    request.FormValue("url"),
		Errors: make(map[string]FieldError),
	}, nil
}

func (f WebAppForm) Validate() {
	if f.Name == "" {
		f.Errors["name"] = NewFieldError("name", "Name is required")
	}

	if f.Url == "" {
		f.Errors["url"] = NewFieldError("url", "URL is required")
	} else {
		url, err := url.ParseRequestURI(f.Url)

		if err != nil {
			f.Errors["url"] = NewFieldError("url", "Invalid URL")
			return
		}

		if url.Scheme == "" || url.Host == "" {
			f.Errors["url"] = NewFieldError("url", "Invalid URL")
			return
		}
	}
}
