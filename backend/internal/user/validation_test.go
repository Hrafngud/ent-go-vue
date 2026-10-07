package user

import (
	"strings"
	"testing"
)

func TestNormalizeInput(t *testing.T) {
	for _, tc := range []struct {
		name           string
		input          UserInput
		editing, valid bool
	}{
		{"valid punctuation", UserInput{" O'Connor-José ", " Person@example.com ", " password "}, false, true},
		{"short password", UserInput{"Name", "a@example.com", "seven77"}, false, false},
		{"unicode password byte minimum", UserInput{"Name", "a@example.com", "😀😀"}, false, true},
		{"unicode password byte maximum", UserInput{"Name", "a@example.com", strings.Repeat("😀", 19)}, false, false},
		{"blank password", UserInput{"Name", "a@example.com", "        "}, false, false},
		{"NUL password", UserInput{"Name", "a@example.com", "password\x00"}, false, false},
		{"bad UTF-8 password", UserInput{"Name", "a@example.com", "password\xff"}, false, false},
		{"optional password", UserInput{"Name", "a@example.com", ""}, true, true},
		{"required password", UserInput{"Name", "a@example.com", ""}, false, false},
		{"markup name", UserInput{"<script>alert(1)</script>", "a@example.com", "password"}, false, false},
		{"control name", UserInput{"Name\x00", "a@example.com", "password"}, false, false},
		{"invalid UTF-8 name", UserInput{"Name\xff", "a@example.com", "password"}, false, false},
		{"long Unicode name", UserInput{strings.Repeat("é", 128), "a@example.com", "password"}, false, false},
		{"empty name", UserInput{"  ", "a@example.com", "password"}, false, false},
		{"display email", UserInput{"Name", "Name <a@example.com>", "password"}, false, false},
		{"invalid local", UserInput{"Name", "a..b@example.com", "password"}, false, false},
		{"invalid domain", UserInput{"Name", "a@-example.com", "password"}, false, false},
		{"local domain", UserInput{"Name", "a@localhost", "password"}, false, false},
		{"long local", UserInput{"Name", strings.Repeat("a", 65) + "@example.com", "password"}, false, false},
		{"long label", UserInput{"Name", "a@" + strings.Repeat("a", 64) + ".com", "password"}, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NormalizeInput(tc.input, !tc.editing)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v, error=%v", tc.valid, err)
			}
		})
	}
	input, err := NormalizeInput(UserInput{" Jose\u0301 ", " Person@example.com ", " password "}, true)
	if err != nil || input.Name != "José" || input.Email != "Person@example.com" || input.Password != " password " {
		t.Fatalf("normalization changed credentials or missed NFC: %#v, %v", input, err)
	}
}
