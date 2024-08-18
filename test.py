import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.mime.application import MIMEApplication

def send_email(subject, body, to_email, from_email, password):
    # SMTP server configuration
    smtp_server = "smtp.outlook.com"
    smtp_port = 587
    
    # Create message container
    msg = MIMEMultipart()
    msg['From'] = from_email
    msg['To'] = to_email
    msg['Subject'] = subject
    
    # Attach the body with the msg instance
    msg.attach(MIMEText(body, 'plain'))
    
    # Setup the server
    server = smtplib.SMTP(smtp_server, smtp_port)
    server.starttls()  # Upgrade the connection to a secure encrypted SSL/TLS connection

    try:
        # Login to the server
        server.login(from_email, password)
        
        # Send email
        server.send_message(msg)
        print("Email sent successfully")
    except Exception as e:
        print(f"Failed to send email: {e}")
    finally:
        # Close the server connection
        server.quit()

# Example usage
if __name__ == "__main__":
    subject = "Test Email"
    body = "This is a test email sent from Python using Outlook SMTP server."
    to_email = "james61124@gmail.com"
    from_email = "james61124@gmail.com"
    password = "ABC123ABC456ABC78981461"  # Ensure you use an app password if you have two-factor authentication enabled

    send_email(subject, body, to_email, from_email, password)