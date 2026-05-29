package pdf

import _ "embed"

//go:embed fonts/DejaVuSans.ttf
var dejaVuRegular []byte

//go:embed fonts/DejaVuSans-Bold.ttf
var dejaVuBold []byte
