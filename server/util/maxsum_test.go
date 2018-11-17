package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_deserialize(t *testing.T) {
	var nilNode *TreeNode
	assert.Equal(t, nilNode, deserialize(&[]string{}))

	root := deserialize(&[]string{"1", "2", "#", "5", "#", "#", "3", "#", "#"})
	assert.Equal(t, 1, root.val)
	assert.Equal(t, 2, root.left.val)
	assert.Equal(t, 3, root.right.val)
	assert.Equal(t, nilNode, root.left.left)
	assert.Equal(t, 5, root.left.right.val)
}

func Test_GetMaxSum(t *testing.T) {
	assert.Equal(t, 0, GetMaxSum(""))
	assert.Equal(t, 8, GetMaxSum("1,2,#,5,#,#,3,#,#"))
	assert.Equal(t, 4, GetMaxSum("1,2,#,#,3"))
	assert.Equal(t, 9, GetMaxSum("1,2,#,5,1,#,#,#,3,6,#,#,7,#,#"))
	assert.Equal(t, 14, GetMaxSum("1,2,#,5,1,#,#,3,#,#,3,6,4,#,#,#,7,#,#"))
}
