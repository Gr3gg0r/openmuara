package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/server"
)

func newSecurityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security",
		Short:   "Security helpers for OpenMuara",
		Long:    "Generate password hashes, self-signed certificates, and audit the security posture.",
		Example: "  muara security hash-password --password mypassword\n  muara security audit",
	}
	cmd.AddCommand(
		newSecurityHashPasswordCommand(),
		newSecurityGenCertCommand(),
		newSecurityAuditCommand(),
	)
	return cmd
}

func newSecurityHashPasswordCommand() *cobra.Command {
	var password string
	cmd := &cobra.Command{
		Use:     "hash-password",
		Short:   "Generate a bcrypt hash of a password",
		Example: "  muara security hash-password --password mypassword",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if password == "" {
				return fmt.Errorf("--password is required")
			}
			hash, err := server.HashPassword(password)
			if err != nil {
				return fmt.Errorf("hash password: %w", err)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), hash)
			return nil
		},
	}
	cmd.Flags().StringVar(&password, "password", "", "password to hash")
	return cmd
}

func newSecurityGenCertCommand() *cobra.Command {
	var host, certOut, keyOut string
	cmd := &cobra.Command{
		Use:     "gen-cert",
		Short:   "Generate a self-signed TLS certificate for local testing",
		Example: "  muara security gen-cert --host localhost --cert-out cert.pem --key-out key.pem",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if host == "" {
				host = "localhost"
			}
			if certOut == "" {
				certOut = "cert.pem"
			}
			if keyOut == "" {
				keyOut = "key.pem"
			}

			priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				return fmt.Errorf("generate key: %w", err)
			}

			tmpl := x509.Certificate{
				SerialNumber: big.NewInt(1),
				Subject:      pkix.Name{Organization: []string{"OpenMuara"}},
				NotBefore:    time.Now(),
				NotAfter:     time.Now().Add(365 * 24 * time.Hour),
				KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				DNSNames:     []string{host},
			}

			certDER, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
			if err != nil {
				return fmt.Errorf("create certificate: %w", err)
			}

			// #nosec G304 -- certOut is a user-provided CLI output path for self-signed cert generation
			certFile, err := os.Create(certOut)
			if err != nil {
				return fmt.Errorf("create cert file: %w", err)
			}
			defer func() { _ = certFile.Close() }()
			if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
				return fmt.Errorf("write cert: %w", err)
			}

			keyBytes, err := x509.MarshalECPrivateKey(priv)
			if err != nil {
				return fmt.Errorf("marshal key: %w", err)
			}
			// #nosec G304 -- keyOut is a user-provided CLI output path for self-signed key generation
			keyFile, err := os.OpenFile(keyOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
			if err != nil {
				return fmt.Errorf("create key file: %w", err)
			}
			defer func() { _ = keyFile.Close() }()
			if err := pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}); err != nil {
				return fmt.Errorf("write key: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "generated %s and %s for host %q\n", certOut, keyOut, host)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "localhost", "host name for the certificate")
	cmd.Flags().StringVar(&certOut, "cert-out", "cert.pem", "output path for the certificate")
	cmd.Flags().StringVar(&keyOut, "key-out", "key.pem", "output path for the key")
	return cmd
}

func newSecurityAuditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "audit",
		Short:   "Print the security posture of the current config",
		Example: "  muara security audit",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(rootConfigPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			var issues []string
			if cfg.Server.Host == "0.0.0.0" && !cfg.Admin.Enabled {
				issues = append(issues, "server is bound to 0.0.0.0 without admin authentication")
			}
			if cfg.Admin.Enabled && cfg.Admin.Username == "" {
				issues = append(issues, "admin is enabled but username is empty")
			}
			if cfg.Admin.Enabled && cfg.Admin.PasswordHash == "" && cfg.Admin.Token == "" {
				issues = append(issues, "admin is enabled but no password_hash or token is set")
			}
			if cfg.Hardened && !cfg.Admin.Enabled {
				issues = append(issues, "hardened mode requires admin.enabled=true")
			}
			hasTLS := cfg.Server.TLSCert != "" && cfg.Server.TLSKey != ""
			if cfg.Server.Host == "0.0.0.0" && cfg.Admin.Enabled && !hasTLS {
				issues = append(issues, "admin auth is enabled but TLS is not configured for 0.0.0.0")
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "bind:        %s\n", cfg.Server.Host)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "admin auth:  %v\n", cfg.Admin.Enabled)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "tls:         %v\n", hasTLS)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "rate limit:  %v\n", cfg.RateLimit.Enabled)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "hardened:    %v\n", cfg.Hardened)

			if len(issues) > 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nissues:")
				for _, issue := range issues {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", issue)
				}
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nno issues detected")
			}

			return nil
		},
	}
	return cmd
}

func init() {
	// Ensure the security command is registered when root.go is loaded.
	securityCmd := newSecurityCommand()
	_ = securityCmd
}
