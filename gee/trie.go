package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild { // child.isWild表示该节点是一个通配符节点，可以匹配任何部分。
			return child // 满足以上任一条件，则说明找到了一个匹配成功的子节点，因此返回该节点
		}
	}
	return nil // 表示没有找到匹配的子节点
}

// 所有匹配成功的节点，用于查找
// 在当前节点的子节点列表中查找与给定部分（part）匹配的所有节点，并将它们存储在一个切片中返回。
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// ..
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// 检查当前节点是否已经是路径的最后一个部分（len(parts) == height），
	// 或者当前节点的部分是一个通配符节点（strings.HasPrefix(n.part, "*")）。
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果当前节点的 pattern 属性为空，则返回 nil，表示未找到匹配的节点。
		// 如果当前节点的 pattern 属性不为空，则返回当前节点 n，表示找到了匹配的节点。
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 在路由树中搜索与给定路径部分匹配的节点，并返回匹配的节点。
	part := parts[height]
	children := n.matchChildren(part) // 获取当前路径部分 parts[height] 匹配的所有子节点 children

	for _, child := range children {
		result := child.search(parts, height+1) // 对每个子节点进行递归搜索，将搜索深度 height 加一
		// 如果子节点的递归搜索返回了一个非空节点 result
		// 则说明已经找到了匹配的节点，直接返回该节点。
		if result != nil {
			return result
		}
	}

	return nil // // 未找到匹配的节点。
}
