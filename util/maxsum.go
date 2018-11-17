package util

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "util")

type TreeNode struct {
	val   int
	left  *TreeNode
	right *TreeNode
}

// preorder deserialize
func deserialize(input *[]string) *TreeNode {
	if len(*input) == 0 {
		return nil
	}
	var cur = (*input)[0]
	*input = (*input)[1:]
	if cur == "#" {
		return nil
	}

	i, err := strconv.Atoi(cur)
	if err != nil {
		logger.Errorf("error parsing char: %s", cur)
		return nil
	}
	var node = TreeNode{val: i}
	node.left = deserialize(input)
	node.right = deserialize(input)
	return &node
}

func getLongestPaths(root *TreeNode, path []int, max *int, ans *[][]int) {
	if root == nil {
		return
	}
	path = append(path, root.val)
	if root.left == nil && root.right == nil {
		level := len(path)
		if level >= (*max) {
			if level > (*max) {
				*ans = [][]int{}
			}
			(*max) = level
			tmp := make([]int, len(path))
			copy(tmp, path)
			(*ans) = append(*ans, tmp)
		}
		return
	}
	if root.left != nil {
		getLongestPaths(root.left, path, max, ans)
	}
	if root.right != nil {
		getLongestPaths(root.right, path, max, ans)
	}
}

// GetMaxSum returns the sum of the longest path, and if there are multiple paths that have the
// same longest length, return the largest sum among those sums
func GetMaxSum(data string) int {
	data = strings.Replace(data, `\s+`, "", -1)
	arr := strings.Split(data, ",")
	logger.Infof("calculating maxsum for tree: %s", data)

	root := deserialize(&arr)
	paths := &[][]int{}
	max := 0
	getLongestPaths(root, []int{}, &max, paths)
	logger.Infof("longest paths: %v", paths)

	maxsum := 0
	for _, path := range *paths {
		sum := 0
		for j := range path {
			sum += path[j]
		}
		if sum > maxsum {
			maxsum = sum
		}
	}
	return maxsum
}
