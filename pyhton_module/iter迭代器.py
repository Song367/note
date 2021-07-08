# 常用迭代器  迭代器： 包括__next__() 和 __iter__()方法就可以称为迭代器   当数据迭代完后，会返回StopIteration
a = [1,2,3]
b = iter(a)   # 通过iter将a 转为迭代器
c = a.__iter__()
print(next(b))
print(next(b))
print(next(b))
print(next(c))

# 自建迭代器
class SuperIter(object):
    def __init__(self):
        self.a, self.b = 0, 1

    def __iter__(self):
        return self

    def __next__(self):
        self.a,self.b=self.b,self.a+self.b

        if self.a > 10:
            raise StopIteration
        return self.a


su = SuperIter()
for i in su:
    print(i)



# open  打开文件时，f 就是一个迭代器

# dropwhile是itemtools当中的一个函数
# 它可以接收一个我们自定义的过滤函数和迭代器重新生成一个新的迭代器
# 过滤掉头部加了#注释的部分

from itertools import dropwhile
with open('xxxx.txt') as f:
    for line in dropwhile(lambda line: line.startswith('#'), f):
        print(line)

# 获取文本第三行开始获取，[3:]
from itertools import islice
with open('xxxx.txt') as f:
    for line in islice(f, 3, None):
        print(line)


# ------------------------------  permutations  返回数组的，多种排列 . 返回的数据都是元组  --------------------------------
items = ['a', 'b', 'c']
from itertools import permutations

for p in permutations(items):
    print(p)

# 还可传入一个参数，指定获取前几位。
for p in permutations(items,2):
    print(p)


# ----------------------------------   combinations 组合 元素与元素之间的组合------------------------------------------------------

item = ['a','b','c']
from itertools import combinations
for c in combinations(item,2):
    print(c)


# --------------------------------------   chain 元素拼接  -----------------------------------

from itertools import chain
nums = [1, 2, 3]
chars = ['a', 'b', 'c']

for i in chain(nums, chars):
    print(i)


# -------------------------------------   heapq 归并  在chain 的基础上还可以将数据进行排序    -------------------------------------

a = [1, 3, 5]
b = [2, 4, 6]
import heapq
for c in heapq.merge(a, b):
    print(c)

