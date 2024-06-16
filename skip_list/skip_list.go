package skip_list

import "math/rand"

type SkipList struct {
	head *Node
}

type Node struct {
	nexts    []*Node
	key, val int
}

// Get 根据 key 读取 val，第二个 bool flag 反映 key 在 skiplist 中是否存在
func (s *SkipList) Get(key int) (int, bool) {
	if node := s.Search(key); node != nil {
		return node.val, true
	}

	return -1, false
}

// Search 从跳表中检索 key 对应的 node
func (s *SkipList) Search(key int) *Node {
	// 每次检索从头部出发
	move := s.head
	// 每次检索从最大高度出发，直到来到首层
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		// 在每一层中持续向右遍历，直到下一个节点不存在或者 key 值大于等于 key
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		// 如果 key 值相等，则找到了目标直接返回
		if move.nexts[level] != nil && move.nexts[level].key == key {
			return move.nexts[level]
		}

		// 没有找到，则下沉至下一层
	}

	// 第一层都没有找到，返回nil
	return nil
}

func (s *SkipList) randomLevel() int {
	var level int
	// 每次投出 1，则层数加 1
	for rand.Intn(2) == 1 {
		level++
	}

	return level
}

// Put 将 key-val 对加入 skiplist
func (s *SkipList) Put(key, val int) {
	if node := s.Search(key); node != nil {
		node.val = val
		return
	}

	level := s.randomLevel()

	for len(s.head.nexts)-1 < level {
		s.head.nexts = append(s.head.nexts, nil)
	}

	newNode := &Node{
		nexts: make([]*Node, level),
		key:   key,
		val:   val,
	}

	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		// 向右遍历，直到右侧节点不存在或者 key 值大于 key
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		// 调整指针关系，完成新节点的插入
		newNode.nexts[level] = move.nexts[level]
		move.nexts[level] = newNode
	}
}

// Del 根据 key 从跳表中删除对应的节点
func (s *SkipList) Del(key int) {
	// 如果 kv 对不存在，则无需删除直接返回
	if node := s.Search(key); node == nil {
		return
	}

	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		// 走到此处意味着右侧节点的 key 值必然等于 key，则调整指针引用
		if move.nexts[level] != nil && move.nexts[level].key == key {
			move.nexts[level] = move.nexts[level].nexts[level]
		}

		// 右侧节点不存在或者 key 值大于 target，则直接跳过
	}

	// 对跳表的最大高度进行更新
	var diff int
	for level := len(s.head.nexts) - 1; level >= 0 && s.head.nexts[level] == nil; level-- {
		diff++
	}

	s.head.nexts = s.head.nexts[:len(s.head.nexts)-diff]
}

// Range 找到 skiplist 当中 ≥ start，且 ≤ end 的 kv 对
func (s *SkipList) Range(start, end int) [][2]int {
	ceilNode := s.ceiling(start)
	if ceilNode == nil {
		return [][2]int{}
	}

	result := make([][2]int, 0)
	move := ceilNode
	for move != nil && move.key <= end {
		result = append(result, [2]int{move.key, move.val})
		move = move.nexts[0]
	}

	return result
}

// Ceiling 找到 skiplist 中，key 值大于等于 target 且最接近于 target 的 key-value 对
func (s *SkipList) Ceiling(target int) ([2]int, bool) {
	node := s.ceiling(target)
	if node == nil {
		return [2]int{}, false
	}

	return [2]int{node.key, node.val}, true
}

// 找到 key 值大于等于 target 且 key 值最接近于 target 的节点
func (s *SkipList) ceiling(target int) *Node {
	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}

		// 如果 key 值等于 target 的 kv 对存在，则直接返回
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// 此时 move 已经对应于在首层 key 值小于 key 且最接近于 key 的节点，其右侧第一个节点即为所寻找的目标节点
	// 如果第一层也没找到到，比如target = 10, 第一层值是[1,2,3]，这时候move.nexts[0]是nil，也是符合要求的
	return move.nexts[0]
}

// Floor 找到 skiplist 中，key 值小于等于 target 且最接近于 target 的 key-value 对
func (s *SkipList) Floor(target int) ([2]int, bool) {
	node := s.floor(target)
	if node == nil {
		return [2]int{}, false
	}

	return [2]int{node.key, node.val}, true
}

// 找到 key 值小于等于 target 且 key 值最接近于 target 的节点
func (s *SkipList) floor(target int) *Node {
	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}

		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	return move
}
