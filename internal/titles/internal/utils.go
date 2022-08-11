package internal

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"syscall"
	"unicode"
	"unsafe"

	"github.com/cetteup/joinme.click-launcher/pkg/refractor_config_handler"
)

// CryptUnprotectData implementation adapted from https://stackoverflow.com/questions/33516053/windows-encrypted-rdp-passwords-in-golang

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

const (
	CRYPTPROTECT_UI_FORBIDDEN  uint32 = 0x1
	globalConKeyDefaultUserRef        = "GlobalSettings.setDefaultUser"
	profileConKeyGamespyNick          = "LocalProfile.setGamespyNick"
	profileConKeyPassword             = "LocalProfile.setPassword"
	// ProfileNumberMaxLength BF2 only uses 4 digit profile numbers
	ProfileNumberMaxLength = 4
)

var (
	crypt32            = syscall.NewLazyDLL("crypt32.dll")
	cryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	kernel32           = syscall.NewLazyDLL("Kernel32.dll")
	localFree          = kernel32.NewProc("LocalFree")
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

func DecryptProfileConPassword(enc string) (string, error) {
	data, err := hex.DecodeString(enc)
	if err != nil {
		return "", err
	}

	dec, err := Decrypt(data)
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

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}

	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])

	return d
}

func Decrypt(data []byte) ([]byte, error) {
	pDataIn := uintptr(unsafe.Pointer(NewBlob(data)))
	var pDataOut DATA_BLOB
	r, _, err := cryptUnprotectData.Call(pDataIn, 0, 0, 0, 0, uintptr(CRYPTPROTECT_UI_FORBIDDEN), uintptr(unsafe.Pointer(&pDataOut)))

	if r == 0 {
		return nil, err
	}

	defer func() {
		_, _, _ = localFree.Call(uintptr(unsafe.Pointer(pDataOut.pbData)))
	}()

	return pDataOut.ToByteArray(), nil
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
