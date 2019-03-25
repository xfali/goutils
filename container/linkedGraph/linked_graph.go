/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2019/2/17
 * @time 14:23
 * @version V1.0
 * Description:
 */

package linkedGraph

import (
	"container/list"
	"github.com/xfali/goutils/container/linkedSet"
)

type LinkedGraph map[interface{}]*linkedSet.LinkedSet

func New() *LinkedGraph {
	return &LinkedGraph{}
}

func (g *LinkedGraph) AddEdge(v interface{}, u interface{}) {
	if _, ok := (*g)[v]; !ok {
		(*g)[v] = linkedSet.New()
	}
	(*g)[v].PushBack(u, true)
	if _, ok := (*g)[u]; !ok {
		(*g)[u] = linkedSet.New()
	}
	(*g)[u].PushBack(v, true)
}

func (g *LinkedGraph) Len() int {
	return len(*g)
}

func (g *LinkedGraph) BFS(begin interface{}, visit func(interface{})) map[interface{}]int {
	queue := list.New()
	queue.PushBack(begin)
	dist := map[interface{}]int{}
	dist[begin] = 0
	i := 0
	for queue.Len() != 0 {
		e := queue.Front()
		top := e.Value
		queue.Remove(e)
		i++

		d := dist[top] + 1
		(*g)[top].Foreach(func(c interface{}) bool {
			if _, ok := dist[c]; !ok {
				dist[c] = d
				queue.PushBack(c)
				visit(c)
			}
			return false
		})
	}
	return dist
}

func (g *LinkedGraph) DFS(begin interface{}, visit func(interface{})) {
	dist := map[interface{}]bool{}
	g.dfs_inner(begin, visit, dist)
}

func (g *LinkedGraph) dfs_inner(begin interface{}, visit func(interface{}), visited map[interface{}]bool) {
	if visited[begin] == true {
		return
	}
	visited[begin] = true
	(*g)[begin].Foreach(func(i interface{}) bool {
		if visited[i] != true {
			visit(i)
			g.dfs_inner(i, visit, visited)
		}
		return false
	})
}
