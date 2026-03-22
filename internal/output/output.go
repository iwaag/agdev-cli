package output

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteJSON(w io.Writer, payload any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func WriteText(w io.Writer, text string) error {
	_, err := fmt.Fprintln(w, text)
	return err
}
