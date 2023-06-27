package utils
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func CreateExif(fileName string, packname string, author string) *string {

	jsonData := map[string]interface{}{
		"sticker-pack-id":       "fdb.my.id 1456",
		"sticker-pack-name":      packname,
		"sticker-pack-publisher": author,
		"android-app-store-link": "https://play.google.com/store/apps/details?id=com.rayark.cytus2",
		"ios-app-store-link":     "https://apps.apple.com/app/id625334537",
		"emojis":                 []string{"👋"},
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	littleEndian := []byte{
		0x49, 0x49, 0x2a, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57,
		0x07, 0x00,
	}

	bytes := []byte{0x00, 0x00, 0x16, 0x00, 0x00, 0x00}

	len := len(jsonBytes)
	var last string

	if len > 256 {
		len = len - 256
		bytes = append([]byte{0x01}, bytes...)
	} else {
		bytes = append([]byte{0x00}, bytes...)
	}

	if len < 16 {
		last = fmt.Sprintf("0%x", len)
	} else {
		last = fmt.Sprintf("%x", len)
	}

	buf2, err := hex.DecodeString(last)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	buf3 := bytes
	buf4 := jsonBytes

	buffer := append(littleEndian, buf2...)
	buffer = append(buffer, buf3...)
	buffer = append(buffer, buf4...)

	err = os.WriteFile(fileName, buffer, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return &fileName
}
