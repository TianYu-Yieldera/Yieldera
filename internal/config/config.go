package config
import "os"
func Env(key, def string) string {
    if v := os.Getenv(key); v != "" { return v }
    return def
}
