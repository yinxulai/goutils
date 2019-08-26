package pattern

import "testing"

func TestIPV4(t *testing.T) {
	tests := []struct {
		value  string
		result bool
	}{
		{"", false},
		{"1.2.1.3", true},
		{"131.256.13.3", false},
		{"255.255.255.255", true},
	}
	for _, tt := range tests {
		if result := IPV4.MatchString(tt.value); result != tt.result {
			t.Errorf("Test pattern IPV4, expect %v, got %v, args %s", tt.result, result, tt.value)
		}
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		value  string
		result bool
	}{
		{"test@test.com", true},
		{"test.test@test.com", true},
		{"test.test@test.test.com", true},
		{"", false},
		{"@test", false},
		{"test.com", false},
		{"@test.com", false},
	}
	for _, tt := range tests {
		if result := Email.MatchString(tt.value); result != tt.result {
			t.Errorf("Test pattern Email, expect %v, got %v, args %s", tt.result, result, tt.value)
		}
	}
}

func TestPhone(t *testing.T) {
	tests := []struct {
		value  string
		result bool
	}{
		{"", false},
		{"1000000000", false},
		{"19736497865", true},
		{"10000000000", true},
		{"20000000000", false},
		{"100000000000", false},
	}
	for _, tt := range tests {
		if result := Phone.MatchString(tt.value); result != tt.result {
			t.Errorf("Test pattern Phone, expect %v, got %v, args %s", tt.result, result, tt.value)
		}
	}
}

func TestNickname(t *testing.T) {
	tests := []struct {
		value  string
		result bool
	}{
		{"", false},
		{"t", false},
		{"test", true},
		{"!@#12", true},
		{"⭐️⭐️⭐️", true},
	}
	for _, tt := range tests {
		if result := Nickname.MatchString(tt.value); result != tt.result {
			t.Errorf("Test pattern Nickname, expect %v, got %v, args %s", tt.result, result, tt.value)
		}
	}
}

func TestUsername(t *testing.T) {
	tests := []struct {
		value  string
		result bool
	}{
		{"", false},
		{"t", false},
		{"test", false},
		{"!@#12", false},
		{"⭐️⭐️⭐️", false},
		{"test12321", true},
		{"testTT12321", true},
	}
	for _, tt := range tests {
		if result := Username.MatchString(tt.value); result != tt.result {
			t.Errorf("Test pattern Username, expect %v, got %v, args %s", tt.result, result, tt.value)
		}
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		value  string
		result bool
	}{
		{"", false},
		{"aaaBBB", true},
		{"aaBB11", true},
		{"aaBB11.", true},
		{"aaBB11.", true},
		{"aaBB11.", true},
		{"aaBB11.", true},
		{"aaBB11@#", true},
		{"aaBB11@#';d1[]", true},
	}
	for _, tt := range tests {
		if result := Password.MatchString(tt.value); result != tt.result {
			t.Errorf("Test pattern Password, expect %v, got %v, args %s", tt.result, result, tt.value)
		}
	}
}
