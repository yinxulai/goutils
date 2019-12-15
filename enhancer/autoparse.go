package enhancer

import "time"

var timeLayoutMap = map[string]string{
	"ansic":       time.ANSIC,
	"unixdate":    time.UnixDate,
	"rubydate":    time.RubyDate,
	"rfc822":      time.RFC822,
	"rfc822z":     time.RFC822Z,
	"rfc850":      time.RFC850,
	"rfc1123":     time.RFC1123,
	"rfc1123z":    time.RFC1123Z,
	"rfc3339":     time.RFC3339,
	"rfc3339nano": time.RFC3339Nano,
	"stamp":       time.Stamp,
	"stampmilli":  time.StampMilli,
	"stampmicro":  time.StampMicro,
	"stampnano":   time.StampNano,
}

// AutoParse 自动匹配 layout 进行 Parse
func AutoParse(timeStr string) time.Time {
	times := map[string]*time.Time{}
	for name, layout := range timeLayoutMap {
		timeD, _ := time.Parse(layout, timeStr)
		times[name] = &timeD
	}

	// 寻找有效且最精确的
	mostAccurate := &time.Time{}
	for _, timeD := range times {
		if timeD != nil && timeD.After(*mostAccurate) {
			mostAccurate = timeD
		}
	}

	return *mostAccurate
}
