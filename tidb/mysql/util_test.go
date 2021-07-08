package mysql

import (
	"encoding/hex"
	"testing"

	"Garyen-go/tidb/util"
)

func TestCalcPassword(t *testing.T) {
	/*
		// **** JDBC ****
		seed:
			@jx=d_3z42;sS$YrS)p|
		hex:
			406a783d645f337a34323b73532459725329707c
		pass:
			kingshard
		scramble:
			fbc71db5ac3d7b51048d1a1d88c1677f34bcca11
	*/
	test := RandomBuf(20)
	hexTest := hex.EncodeToString(test)
	t.Logf("rnd seed: %s, %s", util.String(test), hexTest)

	seed := util.Slice("@jx=d_3z42;sS$YrS)p|")
	hexSeed := hex.EncodeToString(seed)

	t.Logf("seed: %s equal %s, pass: %v", "406a783d645f337a34323b73532459725329707c", hexSeed, "406a783d645f337a34323b73532459725329707c" == hexSeed)
	scramble := CalcPassword(seed, util.Slice("targaryen"))

	hexScramble := hex.EncodeToString(scramble)
	t.Logf("scramble: %s equal %s, pass: %v", "fbc71db5ac3d7b51048d1a1d88c1677f34bcca11", hexScramble, "fbc71db5ac3d7b51048d1a1d88c1677f34bcca11" == hexScramble)
}
