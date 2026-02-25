package mongorm

func String(s string) *string {
	str := string(s)
	return &str
}
func StringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Int(i int) *int {
	return &i
}
func IntVal(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func Float(f float64) *float64 {
	return &f
}
func FloatVal(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func Bool(b bool) *bool {
	return &b
}
func BoolVal(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func Int64(i int64) *int64 {
	return &i
}
func Int64Val(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func Uint(u uint) *uint {
	return &u
}
func UintVal(u *uint) uint {
	if u == nil {
		return 0
	}
	return *u
}

func Uint64(u uint64) *uint64 {
	return &u
}
func Uint64Val(u *uint64) uint64 {
	if u == nil {
		return 0
	}
	return *u
}

func Byte(b byte) *byte {
	return &b
}
func ByteVal(b *byte) byte {
	if b == nil {
		return 0
	}
	return *b
}

func Rune(r rune) *rune {
	return &r
}
func RuneVal(r *rune) rune {
	if r == nil {
		return 0
	}
	return *r
}
