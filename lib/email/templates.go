package email

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
)

var (
	defaultTitle     = "سلام %s عزیز"
	defaultSignature = "ممنون از توجه شما"
)

var mailgen = hermes.Hermes{
	TextDirection: hermes.TDRightToLeft,
	Product: hermes.Product{ // TODO: Provide website's link and logo in production
		Name:      "ویش لیست",
		Copyright: "Copyright © 2020 Wishlist. All rights reserved",
	},
}

// GenEmailConfirmMail is used to generate an email confirmation mail
// containing user's name and confirmation code
func GenEmailConfirmMail(user string, confirmCode string) (string, error) {
	templ := hermes.Email{
		Body: hermes.Body{
			Title: fmt.Sprintf(defaultTitle, user),
			Intros: []string{
				"شما این ایمیل را به علت ثبت نام در سایت ویش لیست دریافت کردید.",
				fmt.Sprintf("کد فعال سازی حساب شما: %s", confirmCode),
			},
			Outros: []string{
				"در غیر اینصورت, اگر شما در سایت ثبت نام نکرده اید و این ایمیل به صورت اشتباه برای شما ارسال شده نیازی به انجام هیچ فرایندی نیست.",
			},
			Signature: defaultSignature,
		},
	}

	email, err := mailgen.GenerateHTML(templ)
	if err != nil {
		return "", err
	}

	return email, nil
}
