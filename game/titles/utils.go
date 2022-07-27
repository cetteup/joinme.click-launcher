package titles

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unicode"
	"unsafe"

	"github.com/cetteup/joinme.click-launcher/internal"
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
	// ProfileNumberMaxLength BF2 only uses 4 digit profile numbers
	ProfileNumberMaxLength = 4
)

var (
	crypt32            = syscall.NewLazyDLL("crypt32.dll")
	cryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	kernel32           = syscall.NewLazyDLL("Kernel32.dll")
	localFree          = kernel32.NewProc("LocalFree")
)

func GetDefaultUserProfileCon(gameFolderName string) (map[string]string, error) {
	profileNumber, err := GetDefaultUserProfileNumber(gameFolderName)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to extract default profile number from global.con: %s", err)
	}

	profileConPath, err := GetProfileConFilePath(gameFolderName, profileNumber)
	if err != nil {
		return map[string]string{}, err
	}

	profileCon, err := ReadParseConFile(profileConPath)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to read profile.con for current default profile (%s): %s", profileNumber, err)
	}

	return profileCon, nil
}

func GetProfilesFolderPath(gameFolderName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, DocumentsFolder, gameFolderName, "Profiles"), nil
}

func GetDefaultUserProfileNumber(gameFolderName string) (string, error) {
	profilesFolder, err := GetProfilesFolderPath(gameFolderName)
	if err != nil {
		return "", err
	}

	// Get preferred profile from Global.con
	globalCon, err := ReadParseConFile(filepath.Join(profilesFolder, GlobalConFile))
	if err != nil {
		return "", err
	}

	profileNumber, ok := globalCon[DefaultUserConKey]
	if !ok || profileNumber == "" {
		return "", fmt.Errorf("reference to default profile is missing/empty")
	}
	// Since BF2 only uses 4 digits for the profile number, 16 bits is plenty to store it
	if _, err := strconv.ParseInt(profileNumber, 10, 16); err != nil || len(profileNumber) > ProfileNumberMaxLength {
		return "", fmt.Errorf("reference to default profile is not a valid profile number: %s", profileNumber)
	}

	return profileNumber, nil
}

// GetEncryptedProfileConLogin Extract profile name and encrypted password from a parsed Profile.con file
func GetEncryptedProfileConLogin(profileCon map[string]string) (string, string, error) {
	nickname, ok := profileCon[ProfileNickConKey]
	// GameSpy nick property is present but empty for local/singleplayer profiles
	if !ok || nickname == "" {
		return "", "", fmt.Errorf("gamespy nickname is missing/empty")
	}
	encryptedPassword, ok := profileCon[ProfilePasswordConKey]
	if !ok || encryptedPassword == "" {
		return "", "", fmt.Errorf("encrypted password is missing/empty")
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

func GetProfileConFilePath(gameFolderName string, profileNumber string) (string, error) {
	profilesFolder, err := GetProfilesFolderPath(gameFolderName)
	if err != nil {
		return "", err
	}
	return filepath.Join(profilesFolder, profileNumber, ProfileConFile), nil
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
		if strings.EqualFold(givenMod, supportedMod) {
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
