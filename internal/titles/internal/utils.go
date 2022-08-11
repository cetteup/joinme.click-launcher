package internal

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unsafe"

	"github.com/cetteup/joinme.click-launcher/pkg/refractor_config_handler"
	"golang.org/x/sys/windows"
)

// CryptUnprotectData implementation adapted from https://stackoverflow.com/questions/33516053/windows-encrypted-rdp-passwords-in-golang
// and https://git.zx2c4.com/wireguard-windows/tree/conf/dpapi/dpapi_windows.go?h=v0.5.3

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

const (
	globalConKeyDefaultUserRef = "GlobalSettings.setDefaultUser"
	profileConKeyGamespyNick   = "LocalProfile.setGamespyNick"
	profileConKeyPassword      = "LocalProfile.setPassword"
	// ProfileNumberMaxLength BF2 only uses 4 digit profile numbers
	ProfileNumberMaxLength = 4
)

func GetDefaultUserProfileCon(configHandler *refractor_config_handler.Handler) (*refractor_config_handler.Config, error) {
	profileNumber, err := GetDefaultUserProfileNumber(configHandler)
	if err != nil {
		return nil, err
	}

	profileCon, err := configHandler.ReadProfileConfig(refractor_config_handler.GameBf2, profileNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to read Profile.con for current default profile (%s): %s", profileNumber, err)
	}

	return profileCon, nil
}

func GetDefaultUserProfileNumber(configHandler *refractor_config_handler.Handler) (string, error) {
	globalCon, err := configHandler.ReadGlobalConfig(refractor_config_handler.GameBf2)
	if err != nil {
		return "", fmt.Errorf("failed to read Global.con: %s", err)
	}

	defaultUserRef, err := globalCon.GetValue(globalConKeyDefaultUserRef)
	if err != nil {
		return "", fmt.Errorf("reference to default profile is missing from Global.con")
	}
	// Since BF2 only uses 4 digits for the profile number, 16 bits is plenty to store it
	if _, err := strconv.ParseInt(defaultUserRef.String(), 10, 16); err != nil || len(defaultUserRef.String()) > ProfileNumberMaxLength {
		return "", fmt.Errorf("reference to default profile in Global.con is not a valid profile number: %s", defaultUserRef.String())
	}

	return defaultUserRef.String(), nil
}

// GetEncryptedProfileConLogin Extract profile name and encrypted password from a parsed Profile.con file
func GetEncryptedProfileConLogin(profileCon *refractor_config_handler.Config) (string, string, error) {
	nickname, err := profileCon.GetValue(profileConKeyGamespyNick)
	// GameSpy nick property is present but empty for local/singleplayer profiles
	if err != nil || nickname.String() == "" {
		return "", "", fmt.Errorf("gamespy nickname is missing/empty")
	}
	encryptedPassword, err := profileCon.GetValue(profileConKeyPassword)
	if err != nil || encryptedPassword.String() == "" {
		return "", "", fmt.Errorf("encrypted password is missing/empty")
	}

	return nickname.String(), encryptedPassword.String(), nil
}

func EncryptProfileConPassword(plain string) (string, error) {
	enc, err := Encrypt([]byte(plain+"\x00"), "This is the description string.")
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(enc), nil
}

func DecryptProfileConPassword(enc string) (string, error) {
	data, err := hex.DecodeString(enc)
	if err != nil {
		return "", err
	}

	dec, _, err := Decrypt(data)
	if err != nil {
		return "", err
	}

	clean := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, string(dec))

	return clean, nil
}

func newBlob(data []byte) *windows.DataBlob {
	if len(data) == 0 {
		return &windows.DataBlob{}
	}

	return &windows.DataBlob{
		Size: uint32(len(data)),
		Data: &data[0],
	}
}

func blobToByteArray(blob windows.DataBlob) []byte {
	bytes := make([]byte, blob.Size)
	copy(bytes, unsafe.Slice(blob.Data, blob.Size))
	return bytes
}

func Encrypt(data []byte, description string) ([]byte, error) {
	dataIn := newBlob(data)
	var dataOut windows.DataBlob
	name, err := windows.UTF16PtrFromString(description)
	if err != nil {
		return nil, err
	}

	if err = windows.CryptProtectData(dataIn, name, nil, uintptr(0), nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &dataOut); err != nil {
		return nil, err
	}

	defer func() {
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(dataOut.Data)))
	}()

	return blobToByteArray(dataOut), nil
}

func Decrypt(data []byte) ([]byte, string, error) {
	dataIn := newBlob(data)
	var dataOut windows.DataBlob
	name, err := windows.UTF16PtrFromString("")
	if err != nil {
		return nil, "", err
	}

	if err = windows.CryptUnprotectData(dataIn, &name, nil, uintptr(0), nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &dataOut); err != nil {
		return nil, "", err
	}

	defer func() {
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(name)))
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(dataOut.Data)))
	}()

	return blobToByteArray(dataOut), windows.UTF16PtrToString(name), nil
}

func BuildOriginURL(offerIDs []string, args []string) string {
	params := url.Values{
		"offerIds":  {strings.Join(offerIDs, ",")},
		"authCode":  {},
		"cmdParams": {url.PathEscape(strings.Join(args, " "))},
	}
	u := url.URL{
		Scheme:   "origin2",
		Path:     "game/launch",
		RawQuery: params.Encode(),
	}
	return u.String()
}
