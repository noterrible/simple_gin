package gee

import (
	"strings"
)

type node struct {
	pattern  string  //待匹配路由 /p/:expl
	part     string  //路由中的一部分，如 :expl
	children []*node //子节点，
	isWild   bool    //是否精确匹配，part含有:、*为true
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if part[0] == ':' && child.isWild && child.part != part { //防止出现/:id和/:name,这是不允许的,同级只能有一个动态路由，服务端是没办法一个请求执行两个动态路由的
			panic("已有相同路径参数的路由")
		}
		if child.part == part { //直接精准定位，会要求使用者必须写成:  /:id,/:id/a,不能写成/:id,/:name/a
			return child
		}
	}
	return nil
}
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children { //优先搜索静态节点
		if child.part == part {
			nodes = append(nodes, child)
		}
	}
	for _, child := range n.children { //没有静态节点，搜索动态
		if child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入结点，插入成功pattern赋值为路径值。插入过程中，没有相应结点创建结点，有“：”或“*”则为模糊节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 搜索结点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //结点高度等于路径长度或带有*，则搜索到结点
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		//匹配子节点
		result := child.search(parts, height+1) //搜索子节点
		if result != nil {
			return result
		}
	}

	return nil
}
