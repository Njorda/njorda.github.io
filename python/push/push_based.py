# The data is pushed out the execution chain(however the data is always pulled from the sink)
# the data is pushed up execution step by step until we reach the top
# where the data is collected

from typing import Callable
from functools import partial

table ={}
table["main"] = [1, 2, 3, 4, 5]

def scan(fn: Callable, table_name: str):
    for i in table[table_name]:
        fn(i)

def filter(fn: Callable, min: int, val: int):
    if val > min:
        fn(val)

def map(fn: Callable, const: int, val: int):
        fn(val*const)

def collect(out:list, val: int):
     out.append(val)

def main():
    out = list()
    col = partial(collect, out)
    mapy = partial(map, col, 2)
    filtery = partial(filter, mapy, 2)
    scan(filtery, "main")
    print(out)

if __name__ == "__main__":
     main()