package mail

import (
	"fmt"
	"slack-clone-api/domain/user"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

func MailHandler(usr user.User, authCode string) error {
	from := mail.NewEmail("Slack", viper.GetString("mail.sender"))
	subject := fmt.Sprintf("Slack confirmation code: %s", authCode)
	to := mail.NewEmail(fmt.Sprintf("%s %s", usr.FirstName, usr.LastName), usr.Email)
	plainTextContent := "Slack Clone Project-Personal Usage"
	htmlContent := fmt.Sprintf(`<html>
	<head>
	  <link rel="preconnect" href="https://fonts.googleapis.com" />
	  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
	  <link
		href="https://fonts.googleapis.com/css2?family=Kanit:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap"
		rel="stylesheet"
	  />
	</head>
	<body style="font-family: Kanit, sans-serif">
	  <img
		width="128"
		height="36"
		src="https://ci5.googleusercontent.com/proxy/EsLhfcFV4o_V7tZ9Ddhyfc6WO4T79EOgunxU5FuKWrZpZnWnW8xg00rFdW1A3zbS_PybqvczfifrQWCb4KC6Bu6A4Qq5OKVWxz-slA=s0-d-e1-ft#https://slack.com/x-a5219734518260/img/slack_logo_240.png"
		alt="Slack Logo"
	  />
	  <p style="font-size: 32px; font-weight: 500">Confirm Your Email Address</p>
	  <p style="font-size: 20px; font-weight: 300">
		Your confirmation code is below — enter it in your open browser<br />window
		and we'll help you get signed in.
	  </p>
	  <table border="0" cellspacing="0">
		<tr>
		  <td
			bgcolor="#0d6bd6"
			style="
			  background-color: #edede9;
			  padding-top: 20px;
			  padding-right: 120px;
			  padding-bottom: 20px;
			  padding-left: 120px;
			  min-width: 50px;
			  border-radius: 4px;
			"
		  >
			<p
			  style="
				font-family: Kanit, sans-serif;
				font-size: 18px;
				text-align: center;
				text-decoration: none;
				letter-spacing: 0.02em;
				color: #3c6e71;
			  "
			>
			  %s
			</p>
		  </td>
		</tr>
	  </table>
	  <p style="font-weight: 300">
		If you didn’t request this email, there’s nothing to worry about — you<br />
		can safely ignore it.
	  </p>
	</body>
  </html>`, authCode)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(viper.GetString("mail.key"))
	rs, err := client.Send(message)

	fmt.Println(rs.StatusCode)
	fmt.Println(rs.Body)
	fmt.Println(rs.Headers)

	return err
}
