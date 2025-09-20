---
layout: post
title: "Vespa Search Engine: Ranking"
subtitle: "How to rank podcasts"
date: 2025-08-30
author: "Niklas Hansson"
URL: "/2025/08/30/vespa_podcast_ranking"
---

This blog post is a follow-up to the [previous post](/2025/08/23/vespa_search_podcast_index) where we focus more on the ranking function. We will assume you have a Vespa cluster up and running (when writing this, we are using the Docker setup from the previous blog post).

The ranking profile in the schema is the following:

```
rank-profile podcast-search {

    inputs {
        query(q) string
    }

    function freshness() {
        expression: exp(-1 * age(newest_item_pubdate)/(3600*24*7)) + attribute(popularity_score)/9
    }

    # https://docs.vespa.ai/en/reference/schema-reference.html#match-phase
    # Attribute that decides which documents are a match if the match 
    # phase estimates that there will be more than max-hits hits.
    match-phase {
        attribute: newest_item_pubdate
        order: descending
        max-hits: 1000
    }

    # https://docs.vespa.ai/en/tutorials/hybrid-search.html
    # Add time factor for aging content 
    # https://docs.vespa.ai/en/nativerank.html#putting-our-features-together-into-a-ranking-expression
    # figure out how to use the chunks in the first phase
    # https://pyvespa.readthedocs.io/en/latest/examples/multilingual-multi-vector-reps-with-cohere-cloud.html
    first-phase {
        expression: bm25(title) + bm25(description) 
    }

    match-features {
        bm25(title)
        bm25(description)
        freshness()
        query(q)
    }  
    
}
```

To start off, a good place to learn more about ranking in Vespa is [this part of the docs](https://docs.vespa.ai/en/getting-started-ranking.html).

Let's understand the ranking profile by trying it out. We will use the [Vespa CLI](https://docs.vespa.ai/en/vespa-cli.html) to execute the queries:

```bash
vespa query \
'yql=select title, description from podcast where true' \
'ranking=podcast-search' \
'input.query(q)="Vespa Voice"' \
'hits=1'
```

Something I struggled with for a long time is this core understanding:

**The values used in the ranking function need to be calculated in the selection phase, otherwise they will stay 0.** This means that the fields used in the ranking phase MUST be involved in the document selection criteria.

This was not clear to me initially and something I struggled to understand from the documentation. Let's test it so you can see it with your own eyes:

```bash
vespa query \
'yql=select title, description from podcast where true' \
'hits=2' \
'ranking=podcast-search' \
'input.query(q)=100'
```

This returns:

```json
{
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 1210
        },
        "coverage": {
            "coverage": 0,
            "documents": 1199,
            "degraded": {
                "match-phase": true,
                "timeout": false,
                "adaptive-timeout": false,
                "non-ideal-state": false
            },
            "full": false,
            "nodes": 1,
            "results": 1,
            "resultsFull": 0
        },
        "children": [
            {
                "id": "index:podcast/0/a87ff679f5e06f0080686efc",
                "relevance": 0.0,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 0.0,
                        "bm25(title)": 0.0,
                        "query(q)": 100.0,
                        "freshness": 1.0
                    },
                    "title": "IdiotSpeakShow",
                    "description": "Podcast by IdiotSpeakShow"
                }
            },
            {
                "id": "index:podcast/0/c4ca42387caa2ab73ba16a03",
                "relevance": 0.0,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 0.0,
                        "bm25(title)": 0.0,
                        "query(q)": 100.0,
                        "freshness": 1.0
                    },
                    "title": "Christianity Questions and Answers",
                    "description": "Dr. Mark Alan Williams and friends answer questions about the Christian faith: questions about the God, Jesus, the Bible, eternity, belief, religion the reasonableness of faith and others."
                }
            }
        ]
    }
}
```

Here we're trying to get the `title` and `description` using `where true` so we match anything. However, we can see in the `matchfeatures` that the scores for different parts of our ranking function:

```
bm25(title) + bm25(description) 
```

The values for `bm25(description)` and `bm25(title)` are both 0.

If we instead query with a specific text match:

```bash
vespa query \
'yql=select title, description from podcast where title contains "Vespa Voice"' \
'hits=1' \
'ranking=podcast-search' \
'input.query(q)=100'
```

We now get:

```json
"matchfeatures": {
    "bm25(description)": 0.0,
    "bm25(title)": 21.227861931727304,
    "query(q)": 100.0,
    "freshness": 8.405524628287946E-9
}
```

Now the BM25 score for the title is calculated because we're actually matching against the title field. However, we don't have a score for the description because it wasn't involved in the selection. 

To get both scores calculated, we need to include both fields in our selection criteria:

```bash
vespa query \
'yql=select * from podcast where title contains "Vespa Voice" or description contains "Vespa Voice"' \
'hits=10' \
'ranking=podcast-search' \
'input.query(q)=100' \
'presentation-summary=bolded'
```

This gives us what we want - both BM25 scores will be calculated properly.

Maybe it's just me, but this was so unclear from the documentation. I thought the scores were calculated separately from the filter conditions, but that is not the case at all. For an example of how this would look with embeddings, see [here](https://docs.vespa.ai/en/embedding.html#embedding-a-query-text).

Finally, note that `input.query(q)=100` isn't doing anything meaningful in our ranking function - we just added it as a feature to demonstrate how you can pass parameters from the query into your ranking expressions. 