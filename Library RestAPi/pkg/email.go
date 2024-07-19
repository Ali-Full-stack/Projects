package pkg

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail(email, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "juraboevizzatillo5@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Registration Status")
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "juraboevizzatillo5@gmail.com", "uqdw nfqi geyh omoy")

	if err := d.DialAndSend(m); err != nil {
		log.Println("failed to send an email:", err)
		return err
	}

	return nil
}

func SendClientCode(code int, name string) string {
	body := fmt.Sprintf(`
    <html>
    <body>
		<h>Confirmation  Code </h>
        <p>Hello Dear <strong>%v</strong>,</p>
        <p>I Hope you are doing Well ..</p>
		<p> Your one time CODE : %v
        <p><strong>PLEASE DO NOT SHARE WITH ANYONE</strong></p>
        <p>Thanks and have a nice day </p>
        <p>from <strong>Ali Team</strong>.........</p>
    </body>
</html>
    `, name, code)
	return body
}

func SendClientResponse(id, name string)string{
	body := fmt.Sprintf(`
    <html>
    <body>
		<h>Congrats, Registration Succesfull </h>
        <p>Hello Dear,  <strong>%v</strong>,</p>
        <p>I Hope you are doing Well ..</p>
		<p> Your ID :<strong> %v</strong>
        <p><strong>PLEASE DO NOT SHARE WITH ANYONE</strong></p>
        <p>Thanks and have a nice day </p>
        <p>from <strong>Ali Team</strong>.........</p>
    </body>
</html>
    `, name, id)
	return body
}
