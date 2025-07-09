package types

import (
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"slices"
	"strconv"
)

type SignInForm struct {
	Email    string
	Password string
	Errors   map[string]FieldError
}

func NewSignInForm(request *http.Request) (SignInForm, error) {
	errors := make(map[string]FieldError)

	if request == nil {
		return SignInForm{Errors: errors}, nil
	}

	if err := request.ParseForm(); err != nil {
		return SignInForm{}, err
	}

	return SignInForm{
		Email:    request.FormValue("email"),
		Password: request.FormValue("password"),
		Errors:   errors,
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
	errors := make(map[string]FieldError)

	if request == nil {
		return SignUpFrom{Errors: errors}, nil
	}

	if err := request.ParseForm(); err != nil {
		return SignUpFrom{}, err
	}

	return SignUpFrom{
		Email:          request.FormValue("email"),
		Password:       request.FormValue("password"),
		RepeatPassword: request.FormValue("repeatPassword"),
		Errors:         errors,
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
	Id       string
	Name     string
	Url      string
	Status   int
	Interval string
	Notify   []string
	Errors   map[string]FieldError
}

func NewWebAppForm(request *http.Request) (WebAppForm, error) {
	errors := make(map[string]FieldError)
	notify := []string{}

	if request == nil {
		return WebAppForm{Status: 200, Notify: notify, Errors: errors}, nil
	}

	if err := request.ParseForm(); err != nil {
		return WebAppForm{}, err
	}

	status, err := strconv.Atoi(request.FormValue("status"))

	if err != nil {
		status = -1
	}

	notify = append(notify, request.Form["notify[]"]...)

	return WebAppForm{
		Name:     request.FormValue("name"),
		Url:      request.FormValue("url"),
		Status:   status,
		Interval: request.FormValue("interval"),
		Notify:   notify,
		Errors:   errors,
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

		if err != nil || url.Scheme == "" || url.Host == "" {
			f.Errors["url"] = NewFieldError("url", "Invalid URL")
		}
	}

	if f.Status <= 0 || f.Status >= 1000 {
		f.Errors["status"] = NewFieldError("status", "Invalid HTTP status")
	}

	if f.Interval == "" {
		f.Errors["interval"] = NewFieldError("interval", "Interval is required")
	} else {
		exists := slices.Contains(Intervals, f.Interval)

		if !exists {
			f.Errors["interval"] = NewFieldError("interval", "Invalid interval")
		}
	}

	if len(f.Notify) > 0 {
		for i, email := range f.Notify {
			field := fmt.Sprintf("email%d", i)

			if email == "" {
				f.Errors[field] = NewFieldError(field, "Address e-mail is required")
				continue
			}

			if _, err := mail.ParseAddress(email); err != nil {
				f.Errors[field] = NewFieldError(field, "Invalid address e-mail")
			}
		}
	}
}
