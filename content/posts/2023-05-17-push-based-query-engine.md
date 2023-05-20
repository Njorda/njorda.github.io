---
layout: post
title: "Push based query engine"
subtitle: "Push based vs pull based query engine for OLAP"
date: 2023-04-14
author: "Niklas Hansson"
URL: "/2023/05/17push_based_query_engine"
---

In this blog post we will dive down in to the difference between [push based vs pull based query engines](https://liuyehcf.github.io/resources/paper/Push-vs-Pull-Based-Loop-Fusion-in-Query-Engines.pdf). As simple as is sounds push based is based upon that data is pushed from the sink through the different operators, this is used by [snowflake and argued to be superior for OLAP](https://info.snowflake.net/rs/252-RFO-227/images/Snowflake_SIGMOD.pdf) which we will dive deeper into. Pull based have been around for a longer time and is based upon that data is pulled from the sink up through the operators, this is also known as the[Volcano Iterator Model](https://db.in.tum.de/~grust/teaching/ws0607/MMDBMS/DBMS-CPU-5.pdf)


This blog post is heavily inspired by [this](https://justinjaffray.com/query-engines-push-vs.-pull/) awesome blog post on the same topic.

## Code examples

### Push based
 
```python 
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
```


### Pull based
```python
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
```


## Does it matter? 

According to [snowflake push-based is superior](https://info.snowflake.net/rs/252-RFO-227/images/Snowflake_SIGMOD.pdf), it is also discussed in this [video from CMU Advanced Databases / Spring 2023](https://www.youtube.com/watch?v=ck3PkXTOueU&t=1040s).

According to the authors the advantages with push based execution engines are:

|  Push-based execution improves cache efficiency, because it removes control flow logic from tight loops. It also enables Snowflake to efficiently process DAG-shaped plans, as opposed to just trees, creating additional opportunities for sharing and pipelining of intermediate results. 

## Cache efficiency

In order to achive the maximum fron hardware knowledge about how computers works are crucial. Leaving control to the [operating system(OS) will result in lower performance](https://www.youtube.com/watch?v=ck3PkXTOueU&list=PLSE8ODhjZXjYzlLMbX3cR0sxWnRM7CLFn&index=7)
and there is today work in building database specific OS. The following sections, will give a brief introduction to CPU architecture and what is relevant for databases. 

For maximum peformance code have to be written for computers not humans. 

### Out-of-order execution

Out-of-order execution refers to that CPU:s execute instructions in the order by the available of the input data rather than the order of the program. This results in that the CPU can avoid idle periods.

CPU instructions are organised in to pipelines where each pipeline([pipelines is a technique for instruction level parallelism](https://en.wikipedia.org/wiki/Instruction_pipelining)) stage is one task. The goal with Pipelining is to keep all parts the CPU busy at all times by dividing the instructions into a series of sequences. By doing so the delays of 
instructions that take more than one cycle can be hidden. Further super scalar CPU:s, (CPU:s that implement instruction-level-parallelism, refers to that, during a single clock cycle multiple instructions can be executed by different execution units on the processor.). 

Out-of-order execution can be broken up to three parts: 

1) Speculative execution - Means that we will speculativly execute the instructions potentially without knowing if it is the exact correct code path or not. This is done in order to avoid having the CPU stall.

2) Branch prediction - predict the instructions which will be most likely to be used next. 

3) Dataflow analysis - reorder the instructions so they are aligned for optimal execution, discarding the original order. 

However this can be come problematic when there are dependencies between instructions since it can not push into the same pipeline. It will also result in issues when the branch prediction was wrong, all speculative work needs to be removed and the pipelines flushed. This will result in lower performance due to the wasted CPU cycles. 

| The most executed branching code in a DBMS is
the filter operation during a sequential scan.
But this is (nearly) impossible to predict correctly.[Source]([https://15721.courses.cs.cmu.edu/spring2023/slides/06-execution.pdf] )

The best solution is try to write code where branching can be avoided. 


## DuckDB case study

[Duckdb](https://github.com/duckdb/duckdb) started out as a pull based system but have since changed to push based execution. This decision where done based upon the analysis that it would result in: 
1) less code duplication, due to that every operator have to ask for data. 
2) hard to do parallel union operators
3) hard to use inlined operators else where. 

for more details check out the [issue](https://github.com/duckdb/duckdb/issues/1583)

## Conclusion

According to [Andy Pavlo](https://15721.courses.cs.cmu.edu/spring2023/slides/06-execution.pdf). The following are the trade offs between push based and pull based. 

Pull)
- Easy to control output via LIMIT.
- Parent operator blocks until its child returns with a tuple.
- Additional overhead because operators' next functions are
implemented as virtual functions.
- Branching costs on each next invocation.

Bottom-to-Top (Push)
- Allows for tighter control of caches/registers in pipelines.
- Difficult to control output via LIMIT.
- Difficult to implement Sort-Merge Join.

Further in the paper [Push vs. Pull-Based Loop Fusion in Query Engines](https://arxiv.org/pdf/1610.09166.pdf) argue for that neither out performs(execution speed of a query) the other when implemented correctly. However [Duckdb switched to push based for easy of use](#duckdb-case-study) and [compiling a push-based query makes for simpler code according](https://arxiv.org/pdf/1610.09166.pdf).

## RANDOM STUFF HERE


I want to build a small tool that can help you learn things faster and ask question based upon some material(Q/A system). The end goal would be to build a small tool you could use to do this based upon any course material such as: 
    - articles(pdf:s)
    - presentations(power points)
    - videos(later on)
I would also like it to be able to generate content based upon some topics(think of a take home exame) with your type of writing, so we would have to accept some input text. 

I think it would be cool if it could run locally(but calling OpenAI remote). But maybe small tool using a shell och just a small webpage, using streamlit or something similar. 



