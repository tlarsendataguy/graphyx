package util

import (
	"bufio"
	"bytes"
	"github.com/danieljoos/wincred"
	"github.com/tlarsendataguy/goalteryx/sdk"
	"golang.org/x/text/encoding/unicode"
	"strings"
)

func GetCredentials(username, password string, provider sdk.Provider) (string, string) {
	if username == `` && password == `` {
		creds, err := wincred.GetGenericCredential(`Graphyx`)
		if err != nil {
			return ``, ``
		}

		builder := strings.Builder{}
		decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
		scanner := bufio.NewScanner(decoder.Reader(bytes.NewReader(creds.CredentialBlob)))
		for scanner.Scan() {
			builder.Write(scanner.Bytes())
		}

		return creds.UserName, builder.String()
	}
	return username, provider.Io().DecryptPassword(password)
}
