package helpers

// ArchReference map used to convert runtime.GOOS to often used abbreviations
var ArchReference = map[string][]string{
	"amd64": {
		"64",
		"x86_64",
		"x64",
		"amd64",
	},
	//TODO: Add x86
	//TODO: Add arm
}

const FileChmod = 0o644
const FileExecutableChmod = 0o755
