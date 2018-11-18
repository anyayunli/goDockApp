package core

import (
	"errors"
	"math"
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
func deserialize(input *[]string) (*TreeNode, error) {
	if len(*input) == 0 {
		return nil, nil
	}
	var cur = (*input)[0]
	*input = (*input)[1:]
	if cur == "#" {
		return nil, nil
	}

	i, err := strconv.Atoi(cur)
	if err != nil {
		return nil, err
	}
	if i < math.MinInt64 || i > math.MaxInt64 {
		return nil, errors.New("Value overflow")
	}
	var node = TreeNode{val: i}
	if node.left, err = deserialize(input); err != nil {
		return nil, err
	}
	if node.right, err = deserialize(input); err != nil {
		return nil, err
	}
	return &node, nil
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
func GetMaxSum(data string) (int64, error) {
	data = strings.Replace(data, `\s+`, "", -1)
	arr := strings.Split(data, ",")

	logger.Infof("calculating maxsum for tree: %s", data)
	root, err := deserialize(&arr)
	if err != nil {
		return 0, err
	}
	paths := &[][]int{}
	max := 0
	getLongestPaths(root, []int{}, &max, paths)
	// logger.Infof("longest paths: %v", paths)

	maxsum := math.MinInt64
	for _, path := range *paths {
		sum := 0
		for j := range path {
			sum += path[j]
		}
		if sum > maxsum {
			maxsum = sum
		}
	}
	if maxsum < math.MinInt64 || maxsum > math.MaxInt64 {
		return 0, errors.New("Value overflow")
	}
	return int64(maxsum), nil
}
