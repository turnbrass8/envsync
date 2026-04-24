package encrypt_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/encrypt"
)

const passphrase = "supersecret"

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	original := "my-secret-value"
	ct, err := encrypt.Encrypt(original, passphrase)
	if err != nil {
		t.Fatalf("Encrypt: unexpected error: %v", err)
	}
	if ct == original {
		t.Fatal("ciphertext should differ from plaintext")
	}
	got, err := encrypt.Decrypt(ct, passphrase)
	if err != nil {
		t.Fatalf("Decrypt: unexpected error: %v", err)
	}
	if got != original {
		t.Errorf("expected %q, got %q", original, got)
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	// AES-GCM uses a random nonce so two encryptions must differ.
	ct1, _ := encrypt.Encrypt("value", passphrase)
	ct2, _ := encrypt.Encrypt("value", passphrase)
	if ct1 == ct2 {
		t.Error("expected two encryptions of the same value to differ")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	ct, _ := encrypt.Encrypt("secret", passphrase)
	_, err := encrypt.Decrypt(ct, "wrongpass")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	_, err := encrypt.Encrypt("value", "")
	if err != encrypt.ErrEmptyPassphrase {
		t.Errorf("expected ErrEmptyPassphrase, got %v", err)
	}
}

func TestDecrypt_EmptyPassphrase(t *testing.T) {
	_, err := encrypt.Decrypt("somedata", "")
	if err != encrypt.ErrEmptyPassphrase {
		t.Errorf("expected ErrEmptyPassphrase, got %v", err)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := encrypt.Decrypt("not-valid-base64!!!", passphrase)
	if err != encrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_TruncatedCiphertext(t *testing.T) {
	ct, _ := encrypt.Encrypt("hello", passphrase)
	// Truncate the base64 string to simulate a short ciphertext.
	truncated := ct[:max(1, len(ct)/4)]
	_, err := encrypt.Decrypt(truncated, passphrase)
	if err == nil {
		t.Error("expected error for truncated ciphertext")
	}
}

func TestEncrypt_EmptyPlaintext(t *testing.T) {
	ct, err := encrypt.Encrypt("", passphrase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := encrypt.Decrypt(ct, passphrase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestEncrypt_OutputIsBase64(t *testing.T) {
	ct, _ := encrypt.Encrypt("test", passphrase)
	if strings.ContainsAny(ct, " \t\n") {
		t.Error("ciphertext should not contain whitespace")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
