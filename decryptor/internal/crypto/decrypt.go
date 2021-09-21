package crypto

import (
	"encoding/hex"
	"errors"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/nacl/box"
	"log"
	"voting2021/decryptor/internal"
)

func DecryptVoteMessage(encryptedChoice *internal.EncryptedChoice, privateKey string) ([]uint32, error) {

	var _message, _ = hex.DecodeString(encryptedChoice.EncryptedMessage)
	var _nonse, _ = hex.DecodeString(encryptedChoice.Nonce)
	var _public_key, _ = hex.DecodeString(encryptedChoice.PublicKey)
	var _private_key, _ = hex.DecodeString(privateKey)

	var Nonse [24]byte
	var PublicKey [32]byte
	var SecretKey [32]byte

	copy(Nonse[:], _nonse[:24])
	copy(PublicKey[:], _public_key[:32])
	copy(SecretKey[:], _private_key[:32])

	decrypted, valid := box.Open(nil, _message, &Nonse, &PublicKey, &SecretKey)
	if !valid {
		return nil, errors.New("unable to decrypt: invalid keys or message: ")
	}

	choices := &Choices{}
	if err := proto.Unmarshal(decrypted[3:], choices); err != nil {
		log.Fatalln("Failed to parse Choices proto", err)
		return nil, err
	}
	return choices.GetData(), nil
}
