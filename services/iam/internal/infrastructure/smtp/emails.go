package smtp

const VerifyEmail string = "Please verify your email: %s"

const ForgotPasswordEmail string = `
Hi, We received a request to reset your password.

Click the link below to choose a new password:
%s

This link expires in %d minutes and can only be used once.

If you didn't request a password reset, you can safely ignore this email.

Thanks,
TodoApp
`
