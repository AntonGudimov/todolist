package domain

import "errors"

var (
	ErrUserAlreadyExists          = errors.New("email already exists")
	ErrUserNotFound               = errors.New("user not found")
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrEmailVerificationNotFound  = errors.New("email verification not found")
	ErrEmailVerificationExpired   = errors.New("email verification expired")
	ErrEmailVerificationUsed      = errors.New("email verification used")
	ErrRefreshTokenNotFound       = errors.New("refresh token not found")
	ErrRefreshTokenInvalid        = errors.New("refresh token invalid")
	ErrPasswordResetTokenNotFound = errors.New("password reset token not found")
	ErrPasswordResetTokenExpired  = errors.New("password reset token expired")
	ErrPasswordResetTokenRevoked  = errors.New("password reset token revoked")
	ErrPasswordResetTokenUsed     = errors.New("password reset token already used")
)

var (
	ErrPrivateKeyNotFound             = errors.New("private key not found")
	ErrPublicKeyNotFound              = errors.New("public key not found")
	ErrKeyMetadataNotFound            = errors.New("key metadata not found")
	ErrInvalidPem                     = errors.New("invalid PEM")
	ErrKeyNotPKCS8PFormat             = errors.New("key is not PKCS8")
	ErrKeyMetadataInvalidFormat       = errors.New("invalid format for key metadata")
	ErrPrivateKeyNotRSA               = errors.New("cannot parse private key to RSA")
	ErrPrivateKeyNotECDSA             = errors.New("cannot parse private key to ECDSA")
	ErrPrivateKeyUnsupportedAlgorithm = errors.New("cannot parse private key")
	ErrPrivateAndPublicKeysDontMatch  = errors.New("private and public keys don't match")
	ErrActiveKeyIDFileNotValid        = errors.New("active key ID file not valid")
	ErrSigningKeyNotFound             = errors.New("signing key not found")
	ErrNoSigningMethodFound           = errors.New("no signing method found")
)
