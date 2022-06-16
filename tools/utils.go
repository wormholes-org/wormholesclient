package tools

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"
	"strings"
)

func SignHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// RecoverAddress recover the address from sig
func RecoverAddress(msg string, sigStr string) (common.Address, error) {
	if !strings.HasPrefix(sigStr, "0x") &&
		!strings.HasPrefix(sigStr, "0X") {
		return common.Address{}, fmt.Errorf("signature must be started with 0x or 0X")
	}
	sigData := hexutil.MustDecode(sigStr)
	if len(sigData) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := hashMsg([]byte(msg))
	fmt.Println("msg ", msg)
	fmt.Println("sigdebug hash=", hexutil.Encode(hash))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

// hashMsg return the hash of plain msg
func hashMsg(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func PriKeyToAddress(priKey string) (account common.Address, fromKey *ecdsa.PrivateKey, err error) {
	fromKey, err = crypto.HexToECDSA(priKey)
	if err != nil {
		return common.Address{}, nil, err
	}
	publicKey := fromKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, nil, err
	}
	account = crypto.PubkeyToAddress(*publicKeyECDSA)
	return
}

func ToHex16(src string) string {
	data := []byte(src)
	rs := hex.EncodeToString(data)
	return "0x" + rs
}

func CheckAddress(name, value string) error {
	if !strings.HasPrefix(value, "0X") && !strings.HasPrefix(value, "0x") {
		return xerrors.Errorf("%s is not string of 0x", name)
	}

	if len(value) != 42 {
		return xerrors.Errorf("the len of %s must be 42", name)
	}
	return nil
}

func CheckHex(name, value string) error {
	if !strings.HasPrefix(value, "0X") && !strings.HasPrefix(value, "0x") {
		return xerrors.Errorf("%s is not string of 0x", name)
	}
	return nil
}

func CheckFlag(name, value string) error {
	if value != "0" && value != "1" {
		return xerrors.Errorf("%s is not the need flag", name)
	}
	return nil
}
