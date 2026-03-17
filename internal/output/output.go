package output

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteSuccess(w io.Writer, jsonMode bool, text string, payload any) error {
	if jsonMode {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	_, err := fmt.Fprintln(w, text)
	return err
}
