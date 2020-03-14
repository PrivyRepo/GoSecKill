package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type uints []uint32

//返回切片长度
func (x uints) Len() int {
	return len(x)
}

func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//当hash环上没有数据时,提示错误
var errEmpty = errors.New("Hash 环没有数据")

//创建结构体保存一致性hash信息
type Consistent struct {
	//hash环,key为哈希值,值存放节点的信息
	circle       map[uint32]string
	sortedHashes uints
	VirtualNode  int
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		//初始化变量
		circle: make(map[uint32]string),
		//设置虚拟节点个数
		VirtualNode: 20,
	}
}

//自动生成Key值
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

//获取hash位置
func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		//声明一个数组长度为64
		var srcatch [64]byte
		//拷贝数据到数组
		copy(srcatch[:], key)
		//使用IEEE多项式返回数据的CRC-32校验和
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	//判断切片的容量
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}
	//添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	//对所有节点hash值进行排序
	sort.Sort(hashes)
}

//向hash环添加节点
func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string) {
	//循环虚拟节点,设置副本
	for i := 0; i < c.VirtualNode; i++ {
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	//更新排序
	c.updateSortedHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

//删除节点
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

//顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用二分查找
	i := sort.Search(len(c.sortedHashes), f)
	//如果超出范围则设置i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

//根据数据标识获取节点信息
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.Unlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	//计算hash值
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
