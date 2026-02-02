package sliceutils

import (
	"testing"
	"utils/utils/testutils"
)

func TestSlice(t *testing.T) {

	t.Run("intersection string", func(t *testing.T) {

		got := IntersectionString([]string{}, []string{})
		testutils.AssertInt(t, len(got), 0)

		got = IntersectionString(
			[]string{},
			[]string{"baba"},
			[]string{"baba"},
		)
		testutils.AssertInt(t, len(got), 0)

		got = IntersectionString(
			[]string{"baba"},
			[]string{"baba"},
			[]string{"baba"},
		)
		testutils.AssertInt(t, len(got), 1)
		testutils.AssertStruct(t, got, []string{"baba"})

		got = IntersectionString(
			[]string{"baba", "baba"},
			[]string{"baba"},
			[]string{"baba", "baba"},
		)
		testutils.AssertInt(t, len(got), 1)
		testutils.AssertStruct(t, got, []string{"baba"})

		got = IntersectionString(
			[]string{"baba", "bobo", "baba"},
			[]string{"baba", "bobo"},
			[]string{"baba", "baba", "bobo"},
		)
		testutils.AssertInt(t, len(got), 2)
		testutils.AssertStruct(t, got, []string{"baba", "bobo"})

	})

	t.Run("process operation", func(t *testing.T) {

		// zero case
		var ini []string
		var values []string

		got := ProcessSliceOperation(
			ini, OpAppend, values,
		)
		testutils.AssertInt(t, len(got), 0)

		got = ProcessSliceOperation(
			ini, OpRemove, values,
		)
		testutils.AssertInt(t, len(got), 0)

		got = ProcessSliceOperation(
			ini, OpOverwrite, values,
		)
		testutils.AssertInt(t, len(got), 0)

		// now some casual tests
		ini = []string{"aa"}
		values = []string{}

		got = ProcessSliceOperation(
			ini, OpAppend, values,
		)
		testutils.AssertStruct(t, got, ini)

		got = ProcessSliceOperation(
			ini, OpRemove, values,
		)
		testutils.AssertStruct(t, got, ini)

		got = ProcessSliceOperation(
			ini, OpOverwrite, values,
		)
		testutils.AssertStruct(t, got, values)

		// now some more casual tests
		ini = []string{"aa"}
		values = []string{"bb"}

		got = ProcessSliceOperation(
			ini, OpAppend, values,
		)
		testutils.AssertStruct(t, got, []string{"aa", "bb"})

		got = ProcessSliceOperation(
			ini, OpRemove, values,
		)
		testutils.AssertStruct(t, got, ini)

		got = ProcessSliceOperation(
			ini, OpOverwrite, values,
		)
		testutils.AssertStruct(t, got, values)

		// now some more casual tests
		ini = []string{"aa", "bb"}
		values = []string{"bb"}

		got = ProcessSliceOperation(
			ini, OpAppend, values,
		)
		testutils.AssertStruct(t, got, []string{"aa", "bb", "bb"})

		got = ProcessSliceOperation(
			ini, OpRemove, values,
		)
		testutils.AssertStruct(t, got, []string{"aa"})

		got = ProcessSliceOperation(
			ini, OpOverwrite, values,
		)
		testutils.AssertStruct(t, got, values)

	})
}
