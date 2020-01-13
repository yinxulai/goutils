package enhancer

import (
	"testing"
	"time"
)

func TestAutoParse(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
	}{
		{"ansic", time.ANSIC},
		{"unixdate", time.UnixDate},
		{"rubydate", time.RubyDate},
		{"rfc822", time.RFC822},
		{"rfc822z", time.RFC822Z},
		{"rfc850", time.RFC850},
		{"rfc1123", time.RFC1123},
		{"rfc1123z", time.RFC1123Z},
		{"rfc3339", time.RFC3339},
		{"rfc3339nano", time.RFC3339Nano},
		{"stamp", time.Stamp},
		{"stampmilli", time.StampMilli},
		{"stampmicro", time.StampMicro},
		{"stampnano", time.StampNano},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp := AutoParse(tt.timeStr)
			wantTime, _ := time.Parse(tt.timeStr, tt.timeStr)
			if !wantTime.Equal(gotResp) {
				t.Errorf("AutoParse(%s) = %v, want %v", tt.name, gotResp, wantTime)
			}

		})
	}
}
