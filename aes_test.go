package mask

import (
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
		key       []byte
	}{
		{
			name:      "basic test",
			plaintext: []byte{1, 2, 3, 4, 5, 6},
			key:       []byte("0123456789abcdef"), // 16 bytes for AES-128
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := AESEncrypt(tt.key, tt.plaintext)
			if err != nil {
				t.Fatalf("AESEncrypt failed: %v", err)
			}

			t.Log(encrypted)

			// Decrypt
			// decrypted, err := AESDecrypt(tt.key, encrypted)
			// if err != nil {
			// 	t.Fatalf("AESDecrypt failed: %v", err)
			// }

			// // Compare original and decrypted
			// if string(decrypted) != string(tt.plaintext) {
			// 	t.Errorf("Decrypted text does not match original.\nWant: %s\nGot:  %s",
			// 		string(tt.plaintext), string(decrypted))
			// }
		})
	}
}
