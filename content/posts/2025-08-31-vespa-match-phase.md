---
layout: post
title: "Vespa Search Engine: Match-phase"
subtitle: "Query degradation and the downsides of match-phase"
date: 2025-08-31
author: "Niklas Hansson"
URL: "/2025/08/31/vespa_match_phase"
---

In this post we will debug a couple of queries to understand the performance and effects of Vespa’s match-phase. This is a follow-up to the [previous post](/2025/08/30/vespa_podcast_ranking). We assume you have the same setup if you want to run the queries yourself.

If we use a previous query as an example:

```bash
vespa query \
'yql=select title, description from podcast where title contains "Vespa AI Search" or description contains "Vespa AI search"' \
'hits=10' \
'ranking=podcast-search' \
'input.query(q)=100'
```

While working on this, I realized something important:

```bash
vespa query \
'yql=select title, description from podcast where title contains "Vespa voice" or description contains "RAG" or description contains "AI"' \
'hits=1' \
'ranking=podcast-search' \
'input.query(q)=100' \
'language=en'
{
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 6119
        },
        "coverage": {
            "coverage": 7,
            "documents": 341462,
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
                "id": "index:podcast/0/68a733e0528d46be938e9d86",
                "relevance": 2.8612187560885154,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 14.306093780442577,
                        "bm25(title)": 0.0,
                        "query(q)": 100.0,
                        "freshness": 0.22501963346236087
                    },
                    "title": "Intelligence Artificielle - Data Driven 101 - Le podcast IA & Data 100% en français",
                    "description": "Sur Data Driven 101, on s’intéresse aux applications pratiques de l'Intelligence Artificielle et de la data dans toute leur diversité avec un objectif : démystifier ces concepts.Dans ce podcast IA & Data (https://datadriven101.tech/) 100% en français, Marc Sanselme reçoit des professionnels de fonctions et d’horizons variés pour nous parler de leurs aventures, leurs succès, leurs échecs, leurs espoirs, leurs techniques, leurs astuces, leurs histoires et leurs convictions.De la Business Intelligence à la Generative AI (LLM, RAG, Agents...) ou à la Computer Vision, toutes les thématiques liées à l'IA sont décortiquées épisode après épisode par Marc Sanselme et ses invités issus de la French tech et d'ailleurs.Marc Sanselme est un spécialiste en Intelligence artificielle (https://scopeo.ai/marc-sanselme/) et dirige la société Draft'n run, studio de développement no-code d'IA sur mesure (https://draftnrun.com/).Équipe : Clémence Reliat, Jean-Christophe Corvisier, Marc SanselmeHébergé "
                }
            }
        ]
    }
}
```

Compare to:

```bash
vespa query \
'yql=select title, description from podcast where title contains "Vespa voice" or description contains "RAG"' \
'hits=1' \
'ranking=podcast-search' \
'input.query(q)=100' \
'language=en'
{
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 491
        },
        "coverage": {
            "coverage": 100,
            "documents": 4625508,
            "full": true,
            "nodes": 1,
            "results": 1,
            "resultsFull": 1
        },
        "children": [
            {
                "id": "index:podcast/0/465a493ab2d24bcd8222dfd6",
                "relevance": 23.0577889391953,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 9.149635037339992,
                        "bm25(title)": 21.227861931727304,
                        "query(q)": 100.0,
                        "freshness": 7.821608121651948E-9
                    },
                    "title": "Vespa Voice",
                    "description": "Welcome to Vespa Voice, the podcast where AI leaders, search pioneers, and enterprise innovators converge. Each episode dives deep into the evolving landscape of AI, featuring candid conversations with experts shaping the future of agentic AI, search architecture, retrieval-augmented generation (RAG), and scalable enterprise applications. Whether you're a CTO driving digital transformation, a CIO reimagining data strategy, or an engineer building next-gen ML and search systems, this is your signal for what's next in intelligent infrastructure."
                }
            }
        ]
    }
}
```

We get very different results. The interesting parts are:

- We only added a extra `contains "AI"`
- Different top result.
- The first query yields a much lower relevance score. Why didn’t the “correct” one win like in the second example?

The culprit: match-phase and graceful degradation

If we compare the coverage, adding the `contains "AI"` term makes the search [degrade](https://docs.vespa.ai/en/graceful-degradation.html#match-phase-degradation) due to [match-phase](https://docs.vespa.ai/en/reference/schema-reference.html#match-phase):

```json
        "coverage": {
            "coverage": 7,
            "documents": 341462,
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
```

What’s happening:

- Adding the broad term `AI` inflates the estimated number of matches in such a way it triggers the match-pahase filtering.
- Because the rank profile is configured with a `match-phase`, Vespa first selects a subset of candidate documents using the match-phase attribute (e.g., a recency/freshness signal) before full ranking.
- Only that subset is then ranked with your relevance function (e.g., BM25 + other features). If the best document isn’t in the subset, it cannot win—leading to a lower score and a different top result.

This explains both the degraded coverage and the unexpected document ordering.

How to mitigate:

- Tighten the query: Avoid very broad terms if they drastically expand the dataset. Prefer phrases or more specific terms, or separate exploratory queries from precise ones.
- Tune match-phase (schema): Choose a match-phase attribute that correlates with final relevance for your use case; consider increasing `max-hits` or disabling match-phase for smaller/static corpora when latency allows. See the [schema reference](https://docs.vespa.ai/en/reference/schema-reference.html#match-phase).
- Use diversity: [Match-phase diversity](https://docs.vespa.ai/en/result-diversity.html#match-phase-diversity) can spread the selected candidates across groups (e.g., by `publisher`, `series`, or another attribute) so a single cluster of similar content doesn’t dominate early selection.
- Observe coverage: Always inspect `coverage.degraded.match-phase`. If true and coverage is low, you likely filtered too early.

Example (rank-profile sketch):

```text
rank-profile podcast-search inherits default {
  match-phase {
    attribute: freshness   # attribute used to select early candidates
    ascending: false       # newer first if freshness is “higher is newer”
    max-hits: 10000        # tune based on corpus size / latency budget
  }
  # ...rest of ranking config...
}
```

If your data is small and relatively static, it can be reasonable to remove or relax match-phase for this profile so the full candidate set is ranked. Otherwise, ensure the match-phase attribute and limits reflect what you actually want to bias early—then validate with `coverage` and result quality.
