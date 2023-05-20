# The data is pulled up where the operator on top start by
# asking for data and thus all the operators below will ask
# until we reach the sink that will pull it up and we work our
# way up again.

from typing import Callable, Iterator
from functools import partial

table ={}
table["main"] = [1, 2, 3, 4, 5]

def scan(table_name: str):
    for i in table[table_name]:
        yield i

def filter(iter: Iterator, min: int) -> int:
    for i in iter:
        if i > min:
            yield i

def map(iter: Iterator, const: int) -> int:
    for i in iter:
        yield i *const

def collect(iter: Iterator, out: list):
    for i in iter:
        out.append(i)

def main():
    out = list()
    collect(map(filter(scan("main"), 2), 2), out)
    print(out)


if __name__ == "__main__":
     main()