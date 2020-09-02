/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2019/2/17
 * @time 14:25
 * @version V1.0
 * Description:
 */

package mapGraph

/*
   符合Graph的逻辑，但是由于底层使用HashMap存储，所以遍历顺序受Hash影响，存在与存入顺序不一致的情况
*/
import (
	"container/list"
)

type MapGraph map[interface{}]map[interface{}]bool

func New() *MapGraph {
	return &MapGraph{}
}

func (g *MapGraph) AddEdge(v interface{}, u interface{}) {
	if _, ok := (*g)[v]; !ok {
		(*g)[v] = map[interface{}]bool{}
	}
	(*g)[v][u] = true
	if _, ok := (*g)[u]; !ok {
		(*g)[u] = map[interface{}]bool{}
	}
	(*g)[u][v] = true
}

func (g *MapGraph) Len() int {
	return len(*g)
}

func (g *MapGraph) BFS(begin interface{}, visit func(interface{})) map[interface{}]int {
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
		for c := range (*g)[top] {
			if _, ok := dist[c]; !ok {
				dist[c] = d
				queue.PushBack(c)
				visit(c)
			}
		}
	}
	return dist
}

func (g *MapGraph) DFS(begin interface{}, visit func(interface{})) {
	dist := map[interface{}]bool{}
	g.dfs_inner(begin, visit, dist)
}

func (g *MapGraph) dfs_inner(begin interface{}, visit func(interface{}), visited map[interface{}]bool) {
	if visited[begin] == true {
		return
	}
	visited[begin] = true
	for i := range (*g)[begin] {
		if visited[i] != true {
			visit(i)
			g.dfs_inner(i, visit, visited)
		}
	}
}
