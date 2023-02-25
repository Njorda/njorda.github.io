# Poor mans datalake 

---
layout:     post 
title:      "Poor mans datalake"
subtitle:   "DuckDb"
date:       2023-02-01
author:     "Niklas Hansson"
URL: "/2023/02/01/Trition_with_post_and_pre_processing/"
---

This post is a deep dive playing with DuckDB doing a twist on [Build a poor man’s data lake from scratch with DuckDB](https://dagster.io/blog/duckdb-data-lake#-the-limitations-of-duckdb) where we will do the following changes:
- Use [Minio](https://github.com/minio/minio) instead of S3 
- [DBT](https://github.com/dbt-labs/dbt-core) instead of dagster.

We will host it on Kubernetes and set it up so it all run locally.


