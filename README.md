# go-bit-filter

对于bit filter，经典的算法如Bloom过滤算法，也有新近的算法Cuckoo过滤算法，其实现如：

- https://github.com/willf/bloom
- https://github.com/seiflotfy/cuckoofilter

这些都是通用的算法，特别是Cuckoo，对冲突做了很巧妙的处理。




## Example

```go
```
## 冲突的数量

### 只使用一个hash函数（fnv.New32a）

使用2MB的空间，即1.67kw左右的表示空间: 

- 元素个数1000w，其hash冲突率约为6.7%
- 元素个数700w，其hash冲突率约为1.3%
- 元素个数600w或者低于600w时，其hash冲突率约为0%

每次用完得alg.Reset()

