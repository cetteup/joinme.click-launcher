package titles

import (
	"encoding/hex"
	"fmt"
	"github.com/cetteup/joinme.click-launcher/internal"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

// CryptUnprotectData implementation adapted from https://stackoverflow.com/questions/33516053/windows-encrypted-rdp-passwords-in-golang

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

const (
	CRYPTPROTECT_UI_FORBIDDEN uint32 = 0x1
	DocumentsFolder                  = "Documents"
	GlobalConFile                    = "Global.con"
	ProfileConFile                   = "Profile.con"
	DefaultUserConKey                = "GlobalSettings.setDefaultUser"
	ProfileNickConKey                = "LocalProfile.setGamespyNick"
	ProfilePasswordConKey            = "LocalProfile.setPassword"
)

var (
	crypt32            = syscall.NewLazyDLL("crypt32.dll")
	cryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	kernel32           = syscall.NewLazyDLL("Kernel32.dll")
	localFree          = kernel32.NewProc("LocalFree")
)

func GetDefaultUserProfileCon(game string) (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	profilesFolder := filepath.Join(homeDir, DocumentsFolder, game, "Profiles")

	// Get preferred profile from Global.con
	globalCon, err := ReadParseConFile(filepath.Join(profilesFolder, GlobalConFile))
	if err != nil {
		return nil, err
	}
	defaultUser, ok := globalCon[DefaultUserConKey]
	if !ok {
		return nil, fmt.Errorf("global.con does not reference a default profile")
	}

	return ReadParseConFile(filepath.Join(profilesFolder, defaultUser, ProfileConFile))
}

// GetEncryptedProfileConLogin Extract profile name and encrypted password from a parsed Profile.con file
func GetEncryptedProfileConLogin(profileCon map[string]string) (string, string, error) {
	nickname, ok := profileCon[ProfileNickConKey]
	if !ok {
		return "", "", fmt.Errorf("profile.con does not contain a gamespy nickname")
	}
	encryptedPassword, ok := profileCon[ProfilePasswordConKey]
	if !ok {
		return "", "", fmt.Errorf("profile.con does not contain an encrypted password")
	}

	return nickname, encryptedPassword, nil
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

func ReadParseConFile(path string) (map[string]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConFile(content), nil
}

func ParseConFile(content []byte) map[string]string {
	lines := strings.Split(string(content), "\r\n")

	parsed := map[string]string{}
	for _, line := range lines {
		elements := strings.SplitN(line, " ", 2)

		// TODO do something other than ignoring any invalid lines hereß
		if len(elements) == 2 {
			// Add key, value or append to value
			value := strings.ReplaceAll(elements[1], "\"", "")
			current, exists := parsed[elements[0]]
			if exists {
				value = strings.Join([]string{current, value}, ",")
			}
			parsed[elements[0]] = value
		}
	}

	return parsed
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

func buildOriginURL(offerIDs []string, args []string) string {
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

func getValidMod(installPath string, modBasePath string, givenMod string, supportedMods ...string) (string, error) {
	var mod string
	for _, supportedMod := range supportedMods {
		if strings.ToLower(givenMod) == strings.ToLower(supportedMod) {
			mod = supportedMod
			break
		}
	}

	if mod == "" {
		return "", fmt.Errorf("mod not supported: %s", givenMod)
	}

	modPath := filepath.Join(installPath, modBasePath, mod)
	installed, err := internal.IsValidDirPath(modPath)
	if err != nil {
		return "", fmt.Errorf("failed to determine whether %s mod is installed: %e", mod, err)
	}
	if !installed {
		return "", fmt.Errorf("mod not installed: %s", mod)
	}

	return mod, nil
}
