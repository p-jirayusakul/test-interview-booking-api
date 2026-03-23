package config

import "time"

func (l *Loader) GetString(key string) string {
	return l.v.GetString(key)
}

func (l *Loader) GetBool(key string) bool {
	return l.v.GetBool(key)
}

func (l *Loader) GetInt(key string) int {
	return l.v.GetInt(key)
}

func (l *Loader) GetInt32(key string) int32 {
	return l.v.GetInt32(key)
}

func (l *Loader) GetInt64(key string) int64 {
	return l.v.GetInt64(key)
}

func (l *Loader) GetUint8(key string) uint8 {
	return l.v.GetUint8(key)
}

func (l *Loader) GetUint(key string) uint {
	return l.v.GetUint(key)
}

func (l *Loader) GetUint16(key string) uint16 {
	return l.v.GetUint16(key)
}

func (l *Loader) GetUint32(key string) uint32 {
	return l.v.GetUint32(key)
}

func (l *Loader) GetUint64(key string) uint64 {
	return l.v.GetUint64(key)
}

func (l *Loader) GetFloat64(key string) float64 {
	return l.v.GetFloat64(key)
}

func (l *Loader) GetTime(key string) time.Time {
	return l.v.GetTime(key)
}

func (l *Loader) GetDuration(key string) time.Duration {
	return l.v.GetDuration(key)
}

func (l *Loader) GetIntSlice(key string) []int {
	return l.v.GetIntSlice(key)
}

func (l *Loader) GetStringSlice(key string) []string {
	return l.v.GetStringSlice(key)
}

func (l *Loader) GetStringMap(key string) map[string]any {
	return l.v.GetStringMap(key)
}

func (l *Loader) GetStringMapString(key string) map[string]string {
	return l.v.GetStringMapString(key)
}

func (l *Loader) GetStringMapStringSlice(key string) map[string][]string {
	return l.v.GetStringMapStringSlice(key)
}

func (l *Loader) GetSizeInBytes(key string) uint {
	return l.v.GetSizeInBytes(key)
}

func (l *Loader) UnmarshalKey(key string, target any) error {
	return l.v.UnmarshalKey(key, target)
}

func (l *Loader) Get(key string) any {
	return l.v.Get(key)
}

func (l *Loader) IsSet(key string) bool {
	return l.v.IsSet(key)
}
