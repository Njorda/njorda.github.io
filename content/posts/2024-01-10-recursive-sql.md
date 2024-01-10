---
layout: post
title: "Recursive queries"
subtitle: "How to traverse graphs in SQL"
date: 2024-01-10
author: "Niklas Hansson"
URL: "/2024/01/10/recursive-sql"
---

I coupe of weeks ago I came across that Postgres supports recursive queries and thus we will take a deep dive in this blog post on the concept of recursion in SQL. The docs for postgres and recursive queries can be found [here](https://www.postgresql.org/docs/current/queries-with.html#QUERIES-WITH-RECURSIVE) if you are not familiar with [CTE](https://www.postgresql.org/docs/current/queries-with.html) I recommend you to start with that. In this blog post we will start with exploring this concept in Postgres before we jump to [DuckDB](https://duckdb.org/).

# RECURSIVE queries in postgres

In order to have a postgres instance to play around with we will use docker, we will not deep dive to much in the setup but rather leave that for your own exploration: 


Start a postgres instance: 

```bash
docker run --rm --name postgresql -p 5432:5432 \
-e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin \
-e POSTGRES_DB=test \
-d postgres:16.1
```

Connect using psql:

```bash
docker exec -it postgresql psql -d test -U admin
```

You should now have a postgres instance with in which you can play around. 

So what is a `Recursive` sql query? A recursive SQL query allow for a query to reference is own output. This is done using the `WITH` key word in combination with `RECURSIVE`, it is important to notice here that the keyword `RECURSIVE` changes the `WITH` statement from being a syntactic feature to instead introducing a new behaviour. Personally I think this introduces some confusion around `RECURSIVE`. Here is an example of a recursive query: 

```SQL
WITH RECURSIVE t(n) AS (
    VALUES (1)
  UNION ALL
    SELECT n+1 FROM t WHERE n < 100
)
SELECT sum(n) FROM t;
```

A recursive query follows the form, of `WITH` followed by a query and then a `UNION ALL` term, the `UNION` term can be though of as the recursive part. In this case the query will sum all integers from 0 to 100, which is 5050. This example is slightly to simple to properly understand recursion and how it can be used so lets make a some what more realistic example. I will take the example from by duckdb [here](https://duckdb.org/docs/sql/query_syntax/with.html#recursive-ctes). 

```SQL
CREATE TABLE tag (id INT, name VARCHAR, subclassof INT);
INSERT INTO tag VALUES
 (1, 'U2',     5),
 (2, 'Blur',   5),
 (3, 'Oasis',  5),
 (4, '2Pac',   6),
 (5, 'Rock',   7),
 (6, 'Rap',    7),
 (7, 'Music',  9),
 (8, 'Movies', 9),
 (9, 'Art', NULL);
```

Feel free to run `SELECT * FROM tag;` to check out the table. Using recursive we could try to answere question such as, what is the path between `2Pac` and `Art`?

```SQL
WITH RECURSIVE steps(id, name, subclassof) AS (
    SELECT id, name, subclassof 
    FROM tag
    WHERE name='2Pac'
  UNION ALL
    SELECT tag.id, tag.name, tag.subclassof
    FROM tag, steps
    WHERE steps.subclassof=tag.id
)
SELECT * FROM steps;
```

The output should be:

```sql
id | name  | subclassof
----+-------+------------
  4 | 2Pac  |          6
  6 | Rap   |          7
  7 | Music |          9
  9 | Art   |
(4 rows)
```

For postgres even though the terminology is recursive the query is evaluated iteratively. This should give you an example on how to traverse a graph like structure using SQL. An alternative do doing this using postgres is to use DuckDB for the compute and graph traversal. 


# RECURSIVE queries from duckdb on postgres. 


The first step is to create a new shell inside the docker container. 

```bash
docker exec -it postgresql /bin/bash
```

```bash
apt-get update
apt-get install curl -y
apt-get install unzip
curl -OL https://github.com/duckdb/duckdb/releases/download/v0.9.2/duckdb_cli-linux-amd64.zip
unzip duckdb_cli-linux-amd64.zip
./duckdb
```

Then it is time to install and load the postgres connector: 

```SQL
D INSTALL postgres_scanner;
D LOAD postgres_scanner;
```

The next step is to connect to the postgres instance: 

```SQL
D ATTACH 'host=localhost port=5432 dbname=test connect_timeout=10 user=admin password=admin' AS postgres_db (TYPE postgres);
D USE postgres_db; -- Otherwise the use postgres_db.public.tag;
D SELECT * FROM tag;
D SELECT * FROM postgres_db.public.tag;
```


We can then run `RECURSIVE` query through duckdb: 


```SQL
WITH RECURSIVE steps(id, name, subclassof) AS (
    SELECT id, name, subclassof 
    FROM tag
    WHERE name='2Pac'
  UNION ALL
    SELECT tag.id, tag.name, tag.subclassof
    FROM tag, steps
    WHERE steps.subclassof=tag.id
)
SELECT * FROM steps;
```

which gives: 

```SQL
┌───────┬─────────┬────────────┐
│  id   │  name   │ subclassof │
│ int32 │ varchar │   int32    │
├───────┼─────────┼────────────┤
│     4 │ 2Pac    │          6 │
│     6 │ Rap     │          7 │
│     7 │ Music   │          9 │
│     9 │ Art     │            │
└───────┴─────────┴────────────┘
```


Thats it, now you are ready to do `RECURSIVE` queries. 
