package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_deserialize(t *testing.T) {
	root, err := deserialize(&[]string{})
	assert.Nil(t, root)
	assert.Nil(t, err)

	root, err = deserialize(&[]string{""})
	assert.Nil(t, root)
	assert.NotNil(t, err)

	root, err = deserialize(&[]string{"9223372036854775808"})
	assert.Nil(t, root)
	assert.NotNil(t, err)

	root, _ = deserialize(&[]string{"1", "2", "#", "5", "#", "#", "3", "#", "#"})
	assert.Equal(t, 1, root.val)
	assert.Equal(t, 2, root.left.val)
	assert.Equal(t, 3, root.right.val)
	assert.Nil(t, root.left.left)
	assert.Equal(t, 5, root.left.right.val)
}

func Test_GetMaxSum(t *testing.T) {
	max, err := GetMaxSum("")
	assert.Equal(t, int64(0), max)
	assert.NotNil(t, err)

	max, err = GetMaxSum("0,0,0,0")
	assert.Equal(t, int64(0), max)
	assert.Nil(t, err)

	max, err = GetMaxSum("1,2,#,5,#,#,3,#,#")
	assert.Equal(t, int64(8), max)
	assert.Nil(t, err)

	max, err = GetMaxSum("1,2,#,5,#,#,3")
	assert.Equal(t, int64(8), max)
	assert.Nil(t, err)

	max, err = GetMaxSum("1,2,#,5,1,#,#,#,3,6,#,#,7,#,#")
	assert.Equal(t, int64(9), max)
	assert.Nil(t, err)

	max, err = GetMaxSum("1,2,#,5,1,#,#,3,#,#,3,6,4,#,#,#,7,#,#")
	assert.Equal(t, int64(14), max)
	assert.Nil(t, err)

	max, err = GetMaxSum("1,-2,#,#,-3")
	assert.Equal(t, int64(-1), max)
	assert.Nil(t, err)

	max, err = GetMaxSum("0,-9223372036854775809")
	assert.Equal(t, int64(0), max)
	assert.NotNil(t, err)

	max, err = GetMaxSum("0,9223372036854775808")
	assert.Equal(t, int64(0), max)
	assert.NotNil(t, err)
}
