package libparsex

import "runtime/debug"

func GetVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "(dev)"
}
