package card

import (
	"sort"
	"sync"

	"github.com/jqiris/saki/utils"
)

// CMap 牌面=>数量
type CMap struct {
	Mux   *sync.RWMutex
	tiles map[int]int
	list  []int
}

// NewCMap 初始化一个TileMap
func NewCMap() *CMap {
	return &CMap{
		Mux:   &sync.RWMutex{},
		tiles: make(map[int]int),
		list:  make([]int, 0),
	}
}

// SetTiles 初始化手牌
func (cm *CMap) SetTiles(tiles []int) {
	cm.Mux.Lock()
	defer cm.Mux.Unlock()
	for _, tile := range tiles {
		cm.tiles[tile]++
	}
	cm.list = tiles
}

// GetTileMap 读取所有牌的列表
// 这里不会主动加锁，在外面用的话，如果用于range，需要手动加锁
func (cm *CMap) GetTileMap() map[int]int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	return cm.tiles
}

// AddTile 添加手牌
func (cm *CMap) AddTile(tile, cnt int) {
	cm.Mux.Lock()
	defer cm.Mux.Unlock()
	cm.tiles[tile] += cnt
	for i := 0; i < cnt; i++ {
		cm.list = append(cm.list, tile)
	}
}

// DelTile 删除手牌
func (cm *CMap) DelTile(tile, cnt int) bool {
	cm.Mux.Lock()
	defer cm.Mux.Unlock()
	if cm.tiles[tile] > cnt {
		cm.tiles[tile] -= cnt
		for i := 0; i < cnt; i++ {
			cm.list = utils.SliceDel(cm.list, tile)
		}
	} else if cm.tiles[tile] == cnt {
		delete(cm.tiles, tile)
		for i := 0; i < cnt; i++ {
			cm.list = utils.SliceDel(cm.list, tile)
		}
	} else {
		return false
	}
	return true
}

// GetList 获得列表
func (cm *CMap) GetList() []int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	tiles := []int{}
	for _, tile := range cm.list {
		tiles = append(tiles, tile)
	}
	return tiles
}

// ToSortedSlice 转成slice并排序
func (cm *CMap) ToSortedSlice() []int {
	tiles := cm.GetList()
	sort.Ints(tiles)
	return tiles
}

// GetUnique 获取独立的牌
func (cm *CMap) GetUnique() []int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	tiles := []int{}
	for tile := range cm.tiles {
		tiles = append(tiles, tile)
	}
	return tiles
}

// GetTileCnt 获取某张牌的数量
func (cm *CMap) GetTileCnt(tile int) int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	return cm.tiles[tile]
}

func (cm *CMap) GetNumTiles(num int) []int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	tiles := []int{}
	for tile, cnum := range cm.tiles {
		if cnum == num {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

func (cm *CMap) Pop() int {
	cm.Mux.Lock()
	defer cm.Mux.Unlock()
	length := len(cm.list)
	if length == 0 {
		return -1
	}
	tile := cm.list[length-1]
	cm.list = cm.list[:length-1]
	cm.tiles[tile]--
	return tile
}

func (cm *CMap) Peek() int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	length := len(cm.list)
	if length == 0 {
		return -1
	}
	return cm.list[length-1]
}

func (cm *CMap) GetIndexTile(index int) int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	length := len(cm.list)
	if index >= length {
		return -1
	}
	return cm.list[index]
}

func (cm *CMap) GetListNum() int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	return len(cm.list)
}

func (cm *CMap) GetIndexType(index int) int {
	if index >= len(cm.list) {
		return -1
	}
	tile := cm.list[index]
	if IsDot(tile) {
		return 0
	}
	if IsBAM(tile) {
		return 1
	}
	if IsCrak(tile) {
		return 2
	}
	return -1
}
