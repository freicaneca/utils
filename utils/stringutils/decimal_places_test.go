package stringutils

import (
	"testing"
	"utils/utils/testutils"
)

func TestDecimal(t *testing.T) {

	t.Run("format crazy decimal places", func(t *testing.T) {

		got := FormatValue(15, 4)
		testutils.AssertString(t, got, "0.0015")

		got = FormatValue(-15, 1)
		testutils.AssertString(t, got, "-1.5")

		got = FormatValue(-7200, 0)
		testutils.AssertString(t, got, "-7200")

		got = FormatValue(-7242, 5)
		testutils.AssertString(t, got, "-0.07242")

		got = FormatValue(7242, 1)
		testutils.AssertString(t, got, "724.2")

		got = FormatValue(300, 5)
		testutils.AssertString(t, got, "0.00300")

		got = FormatValue(0, 0)
		testutils.AssertString(t, got, "0")

		got = FormatValue(0, 1)
		testutils.AssertString(t, got, "0.0")

		got = FormatValue(5, 1)
		testutils.AssertString(t, got, "0.5")

	})

	t.Run("parse crazy decimal places", func(t *testing.T) {

		got, err := ParseValue("0.00", 4)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 0)

		got, err = ParseValue("1.20000000", 5)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 120000)

		got, err = ParseValue("4432.99", 1)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 44329)

		got, err = ParseValue("0.83", 4)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 8300)

		got, err = ParseValue("11", 1)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 110)

		got, err = ParseValue("11.3", 4)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 113000)

		got, err = ParseValue("0.3", 4)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 3000)

		// negativos
		got, err = ParseValue("-1.20", 4)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -12000)

		got, err = ParseValue("-4432.99", 4)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -44329900)

		got, err = ParseValue("-0.83", 1)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -8)

		got, err = ParseValue("-11", 0)
		testutils.AssertError(t, err, nil)
		testutils.AssertInt(t, got, -11)

		got, err = ParseValue("-11.3", 1)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -113)

		got, err = ParseValue("-0.3", 1)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -3)

	})

	t.Run("parse cents", func(t *testing.T) {

		got, err := ParseValue("0.00", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 0)

		got, err = ParseValue("1.20", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 120)

		got, err = ParseValue("4432.99", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 443299)

		got, err = ParseValue("0.83", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 83)

		got, err = ParseValue("11", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 1100)

		got, err = ParseValue("11.3", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 1130)

		got, err = ParseValue("0.3", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, 30)

		// negativos
		got, err = ParseValue("-1.20", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -120)

		got, err = ParseValue("-4432.99", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -443299)

		got, err = ParseValue("-0.83", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -83)

		got, err = ParseValue("-11", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -1100)

		got, err = ParseValue("-11.3", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -1130)

		got, err = ParseValue("-0.3", 2)
		testutils.AssertBool(t, err == nil, true)
		testutils.AssertInt(t, got, -30)

	})

	t.Run("format cents", func(t *testing.T) {

		got := FormatValue(15, 2)
		testutils.AssertString(t, got, "0.15")

		got = FormatValue(-15, 2)
		testutils.AssertString(t, got, "-0.15")

		got = FormatValue(-7200, 2)
		testutils.AssertString(t, got, "-72.00")

		got = FormatValue(-7242, 2)
		testutils.AssertString(t, got, "-72.42")

		got = FormatValue(7242, 2)
		testutils.AssertString(t, got, "72.42")

		got = FormatValue(300, 2)
		testutils.AssertString(t, got, "3.00")

		got = FormatValue(0, 2)
		testutils.AssertString(t, got, "0.00")

	})

}
