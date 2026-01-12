package services

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// EmailService handles email sending
type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("SMTP_FROM_NAME")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 587
	}

	dialer := gomail.NewDialer(host, port, user, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}

	from := fmt.Sprintf("%s <%s>", fromName, fromEmail)

	return &EmailService{
		dialer: dialer,
		from:   from,
	}
}

// SendEmail sends a basic email
func (s *EmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}

// SendBookingConfirmation sends a booking confirmation email
func (s *EmailService) SendBookingConfirmation(to, clientName, eventType, eventDate, language string) error {
	var subject, body string
	if language == "en" {
		subject = "Your Booking Confirmation - Angel Event"
		body = s.getBookingConfirmationTemplateEn(clientName, eventType, eventDate)
	} else {
		subject = "Confirmation de votre réservation - Angel Event"
		body = s.getBookingConfirmationTemplateFr(clientName, eventType, eventDate)
	}
	return s.SendEmail(to, subject, body)
}

// SendContactFormEmail sends contact form submission to admin
func (s *EmailService) SendContactFormEmail(name, email, phone, message string) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	subject := fmt.Sprintf("Nouvelle demande de contact - %s", name)
	body := s.getContactFormTemplate(name, email, phone, message)
	return s.SendEmail(adminEmail, subject, body)
}

// SendAdminBookingNotification sends notification to admin about new booking
func (s *EmailService) SendAdminBookingNotification(clientName, clientEmail, clientPhone, eventType, eventDate, location string, guestCount int, budget float64, message string) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	subject := fmt.Sprintf("🎉 Nouvelle Réservation - %s", eventType)
	body := s.getAdminBookingTemplate(clientName, clientEmail, clientPhone, eventType, eventDate, location, guestCount, budget, message)
	return s.SendEmail(adminEmail, subject, body)
}

// SendNewsletterEmail sends a newsletter email
func (s *EmailService) SendNewsletterEmail(to, subject, content string) error {
	body := s.getNewsletterTemplate(content)
	return s.SendEmail(to, subject, body)
}

// Email Templates

func (s *EmailService) getBookingConfirmationTemplateFr(name, eventType, eventDate string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: 'Inter', Arial, sans-serif; line-height: 1.6; color: #1A1A1A; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { text-align: center; padding: 30px 0; background: linear-gradient(135deg, #FFFEF7 0%%, #F7E7CE 100%%); }
		.logo { font-family: 'Great Vibes', cursive; font-size: 36px; color: #D4AF37; }
		.content { padding: 30px; background: white; }
		.footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
		.button { display: inline-block; padding: 12px 30px; background: #D4AF37; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<div class="logo">Angel Event</div>
			<p style="color: #666; margin: 10px 0;">Créer l'instant parfait</p>
		</div>
		<div class="content">
			<h2 style="color: #D4AF37;">Merci pour votre réservation !</h2>
			<p>Bonjour %s,</p>
			<p>Nous avons bien reçu votre demande de réservation pour votre <strong>%s</strong> prévu le <strong>%s</strong>.</p>
			<p>Notre équipe va examiner votre demande et vous contactera sous peu pour confirmer tous les détails et personnaliser votre événement.</p>
			<p>Nous sommes impatients de créer avec vous un moment inoubliable !</p>
			<p style="margin-top: 30px;">Cordialement,<br><strong>L'équipe Angel Event</strong></p>
		</div>
		<div class="footer">
			<p>Angel Event - L'art de sublimer vos moments précieux</p>
			<p>Pour toute question, contactez-nous à contact@angelevent.com</p>
		</div>
	</div>
</body>
</html>
	`, name, eventType, eventDate)
}

func (s *EmailService) getBookingConfirmationTemplateEn(name, eventType, eventDate string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: 'Inter', Arial, sans-serif; line-height: 1.6; color: #1A1A1A; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { text-align: center; padding: 30px 0; background: linear-gradient(135deg, #FFFEF7 0%%, #F7E7CE 100%%); }
		.logo { font-family: 'Great Vibes', cursive; font-size: 36px; color: #D4AF37; }
		.content { padding: 30px; background: white; }
		.footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
		.button { display: inline-block; padding: 12px 30px; background: #D4AF37; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<div class="logo">Angel Event</div>
			<p style="color: #666; margin: 10px 0;">Creating the perfect moment</p>
		</div>
		<div class="content">
			<h2 style="color: #D4AF37;">Thank you for your booking!</h2>
			<p>Hello %s,</p>
			<p>We have received your booking request for your <strong>%s</strong> scheduled for <strong>%s</strong>.</p>
			<p>Our team will review your request and contact you shortly to confirm all details and customize your event.</p>
			<p>We look forward to creating an unforgettable moment with you!</p>
			<p style="margin-top: 30px;">Best regards,<br><strong>The Angel Event Team</strong></p>
		</div>
		<div class="footer">
			<p>Angel Event - The art of sublimating your precious moments</p>
			<p>For any questions, contact us at contact@angelevent.com</p>
		</div>
	</div>
</body>
</html>
	`, name, eventType, eventDate)
}

func (s *EmailService) getContactFormTemplate(name, email, phone, message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; background: #f9f9f9; }
		.content { background: white; padding: 30px; border-radius: 8px; }
		.field { margin-bottom: 15px; }
		.label { font-weight: bold; color: #D4AF37; }
	</style>
</head>
<body>
	<div class="container">
		<div class="content">
			<h2 style="color: #D4AF37;">Nouvelle demande de contact</h2>
			<div class="field">
				<div class="label">Nom:</div>
				<div>%s</div>
			</div>
			<div class="field">
				<div class="label">Email:</div>
				<div>%s</div>
			</div>
			<div class="field">
				<div class="label">Téléphone:</div>
				<div>%s</div>
			</div>
			<div class="field">
				<div class="label">Message:</div>
				<div>%s</div>
			</div>
		</div>
	</div>
</body>
</html>
	`, name, email, phone, message)
}

func (s *EmailService) getNewsletterTemplate(content string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: 'Inter', Arial, sans-serif; line-height: 1.6; color: #1A1A1A; }
		.container { max-width: 600px; margin: 0 auto; }
		.header { text-align: center; padding: 40px 20px; background: linear-gradient(135deg, #FFFEF7 0%%, #F7E7CE 100%%); }
		.logo { font-family: 'Great Vibes', cursive; font-size: 42px; color: #D4AF37; }
		.content { padding: 40px 20px; background: white; }
		.footer { text-align: center; padding: 20px; color: #666; font-size: 12px; background: #f9f9f9; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<div class="logo">Angel Event</div>
		</div>
		<div class="content">
			%s
		</div>
		<div class="footer">
			<p>Angel Event - Des émotions mises en scène</p>
			<p><a href="#" style="color: #D4AF37;">Se désabonner</a></p>
		</div>
	</div>
</body>
</html>
	`, content)
}

func (s *EmailService) getAdminBookingTemplate(clientName, clientEmail, clientPhone, eventType, eventDate, location string, guestCount int, budget float64, message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; background: #f9f9f9; }
		.header { background: linear-gradient(135deg, #D4AF37 0%%, #C5A028 100%%); color: white; padding: 30px; border-radius: 8px 8px 0 0; text-align: center; }
		.content { background: white; padding: 30px; border-radius: 0 0 8px 8px; }
		.field { margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px solid #eee; }
		.field:last-child { border-bottom: none; }
		.label { font-weight: bold; color: #D4AF37; margin-bottom: 5px; }
		.value { color: #333; }
		.highlight { background: #FFFEF7; padding: 15px; border-left: 4px solid #D4AF37; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1 style="margin: 0;">🎉 Nouvelle Réservation</h1>
			<p style="margin: 10px 0 0 0; opacity: 0.9;">Une nouvelle demande vient d'être reçue</p>
		</div>
		<div class="content">
			<div class="field">
				<div class="label">Type d'événement:</div>
				<div class="value"><strong>%s</strong></div>
			</div>
			<div class="field">
				<div class="label">Date de l'événement:</div>
				<div class="value">%s</div>
			</div>
			<div class="field">
				<div class="label">Lieu:</div>
				<div class="value">%s</div>
			</div>
			<div class="field">
				<div class="label">Nombre d'invités:</div>
				<div class="value">%d personnes</div>
			</div>
			<div class="field">
				<div class="label">Budget:</div>
				<div class="value"><strong>%.2f $</strong></div>
			</div>
			<div class="highlight">
				<h3 style="margin-top: 0; color: #D4AF37;">Informations Client</h3>
				<div class="field" style="border: none; margin: 10px 0; padding: 5px 0;">
					<div class="label">Nom:</div>
					<div class="value">%s</div>
				</div>
				<div class="field" style="border: none; margin: 10px 0; padding: 5px 0;">
					<div class="label">Email:</div>
					<div class="value"><a href="mailto:%s">%s</a></div>
				</div>
				<div class="field" style="border: none; margin: 10px 0; padding: 5px 0;">
					<div class="label">Téléphone:</div>
					<div class="value">%s</div>
				</div>
			</div>
			<div class="field">
				<div class="label">Message du client:</div>
				<div class="value">%s</div>
			</div>
			<div style="text-align: center; margin-top: 30px; padding: 20px; background: #FFFEF7; border-radius: 8px;">
				<p style="margin: 0; color: #666;">🔔 Connectez-vous à l'administration pour gérer cette réservation</p>
			</div>
		</div>
	</div>
</body>
</html>
	`, eventType, eventDate, location, guestCount, budget, clientName, clientEmail, clientEmail, clientPhone, message)
}
