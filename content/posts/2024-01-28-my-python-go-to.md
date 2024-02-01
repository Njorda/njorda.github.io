---
layout: post
title: "Python tricks"
subtitle: "These are the python tricks I like to find fast"
date: 2024-01-28
author: "Niklas Hansson"
URL: "/2024/01/21/python-tricks.md"
---

The blog post will just be a long list of short summaries of python tricks, that I find good and like to remember(and find fast). Lets go!

# Table of Contents
1. [Example](#example)
2. [Example2](#example2)
3. [Third Example](#third-example)
4. [Fourth Example](#fourth-examplehttpwwwfourthexamplecom)


## List
### List comprehensions

List comprehensions are a neat one liner python trick. Not only is it handy it also sometimes faster. 

```python
a = [1, 2, 3, 4, 5]
b = [val**2 for val in a]
```


### Lazy initialize list

To initialize a list with x values that have the same number:


```python
a=[1]*1000
```

we can do the same for lists of x values: 


```python
a = [1,2,3]
b = a *10
b[1, 2, 3, 1, 2, 3, 1, 2, 3, ...]
```

## Strings

### Concatenate strings
```python
a = "hello "
b = "world"
print(a + b)
```

### Selecting chars in string

Strings can be index as a list

```python
a = "hello"
print(a[1])
```

### Get length of string

To get the length of a string use `len`

```python
a = "hello"
print(len(a))
```


### Replace in string

```python
"Hello world".replace("world", "John")
```

### Count string or substring

```python
word = "Hello world"
print(word.count(" "))
```

### To repeat string

```python
word  = "Hello world" * 3
print(word)
```

### Split

```python
a = "hello world"
b = a.split(" ")
print(a[0])
```

### Remove white spaces

```python


```


## Collections
## Python general

### Analyze memory

In order to check the memory:

```
import sys
a=10
print(sys.getsizeof(a))
```

### Swapping values

```python
a = 10
b = 20
a, b = b, a 
```

### Invert dict

The easiest way to invert a dict is to use a dict comprehension: 

```
a = {‘a’: 1, ‘b’: 2,‘c’: 3,‘d’: 4,‘e’: 5,‘f’: 6, ‘g’: 7}
b = {v: k for k, v in dict1.items()}
print(b)
```

### Combine two python list

In order to combine to lists `extends` is the way to go, if you use `append` it will result in that you get a list with a list within. 

```python
a = [1, 2, 3]
b = [2, 4, 6]

a.extends(b)
print(a)
```

### Transpose matrix: 

Cool trick using `zip` for transposing a matrix:

```python
mat = [[8, 9, 10], [11, 12, 13]]
mut = zip(*mat)
trans = list(mut)
```


### Sorting

Python offers build in sorting through some of the data structures but also through `sorted` which is a function available. 

To use it to sort a list: 

```python
a = [1,2,5,3,2,1,8]
b = sorted(a)
print(b)
```

If you have a multi layered data structure such a dict you can handle it like this: 


```python
a = [{"company": "Data&IT", "employees": 6}, 
        {"company": "BK", "employees": 800}, 
        {"company": "ICA", "employees": 9000}
        ]
b = sorted(a, key=lambda x: x["employees"])
print(b)
```

Default it is ascending order, to get id descending set `reverse=True`

```python
a = [1,2,5,3,2,1,8]
b = sorted(a, reverse=True)
print(b)
```

### Context Managment - With

### Slots

### Ignore exceptions -  contextlib.suppress

### All/Any

As the name sounds `all` checks if all values in a iterable data structure is true. 

```python
a = [1, 2, 3, 4, 5]
b = [val > 4 for val in a]
print(all(b))
```

`any` instead checks if any of the values are true. 


```
a = [1, 2, 3, 4, 5]
b = [val > 4 for val in a]
print(any(b))
```

### Lambda functions

The keyword `lambda` allows for creating small anonymous functions and has the following syntax: 

```python
lambda arguments : expression
```

example: 

```python
x = lambda a : a + 10
print(x(5))
```

### Generators/yield

Python offers generators as a way to access the value at runtime instead of calculating the whole list first. This allows for saving on memory and is especially useful working with large objects where downstream processes will not need all at once. Using a comprehension expression it looks like this: 

```python
a = [1, 2, 3, 4, 5]
b = (val**2 for val in a)
print(next(b))
```

the important difference between the list comprehension and a generator comprehension is the square brackets vs the optional brackets. It can also be create using the `yield` keyword and a normal function. 

```python
def square_list(n: list) -> list:
    for val in n: 
        yield n 

a = [1, 2, 3, 4, 5]
print(next(square_list(a)))
```

### import itertools

To unpack evenly structure lists of list python offers `itertools` which is part of the standard library

```python
import itertools
a = [[1, 2], [3, 4], [5, 6]]
b = list(itertools.chain.from_iterable(a))
print(b)
```

It will unpack one level. 


### from collections import

The collections package is part of the standard lib and offers multiple useful tools. 

#### defaultdict

`defaultdict` offers the possibility to define default return values for a dict. Normally the return would be `KeyError!`
with a defaultdict a value to be set when trying to access a value that is not set. The default value can be one of: `List`, `Set` or `Int`. However you can explicitly always set a key to any kind like a normal dict. 

```python
from collections import defaultdict

_map = defaultdict(list)
_map['a'].append(1)
_map['b'].append(2)
print(_map['c'])
_map['d'] = ""
```

#### ordereddict

#### Counter

Counter allows for counting any hashable object(What is a hashable object? Lets not discuss it to much but if you stay with native types you are good. If not check out __hash__). Say we want to count the distinct values in a array and the the three most common ones: 

```python
from collections import Counter
a = [1, 1, 1, 2, 3, 4, 11, 13, 1, 1, 1, 2]
counter = Counter(a)

most_common = counter.most_common(2)
print(most_common) # [(9, 6), (10, 3)]
print(most_common[0]) # (9, 6)
print(most_common[0][0]) # 9
```

It can also be used to check anagrams:

```python
from collections import Counter

def is_anagram(str1, str2):
    return Counter(str1) == Counter(str2)
print(is_anagram(‘taste’, ‘state))
print(is_anagram(‘beach’, ‘peach’))
```

#### deque



## The Zen of Python

[Link](https://peps.python.org/pep-0020/)



