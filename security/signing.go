package security

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"

	"OpenZeppelin/fortify-node/protocol"
)

// LoadKey loads the node private key.
func LoadKey() (*keystore.Key, error) {
	f, err := os.OpenFile("/passphrase", os.O_RDONLY, 400)
	if err != nil {
		return nil, err
	}

	pw, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	passphrase := string(pw)

	files, err := ioutil.ReadDir("/.keys")
	if err != nil {
		return nil, err
	}

	if len(files) != 1 {
		return nil, errors.New("there must be only one key in key directory")
	}

	keyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", "/.keys", files[0].Name()))
	if err != nil {
		return nil, err
	}

	return keystore.DecryptKey(keyBytes, passphrase)
}

func SignAlert(key *keystore.Key, alert *protocol.Alert) (*protocol.SignedAlert, error) {
	b, err := proto.Marshal(alert)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256(b)
	sig, err := crypto.Sign(hash, key.PrivateKey)
	if err != nil {
		return nil, err
	}
	signature := fmt.Sprintf("0x%s", hex.EncodeToString(sig))
	return &protocol.SignedAlert{
		Alert: alert,
		Signature: &protocol.Signature{
			Signature: signature,
			Algorithm: "ECDSA",
		},
	}, nil
}

// NewTransactOpts creates new opts with the private key.
func NewTransactOpts(key *keystore.Key) *bind.TransactOpts {
	return bind.NewKeyedTransactor(key.PrivateKey)
}
