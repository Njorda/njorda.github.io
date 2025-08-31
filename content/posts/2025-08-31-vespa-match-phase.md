---
layout: post
title: "Vespa Search Engine: Highlight search results"
subtitle: "Dynamic bolding"
date: 2025-08-31
author: "Niklas Hansson"
URL: "/2025/08/30/vespa_highlights_bolding"
---

This blog post is a follow-up to the [previous post](/2025/08/23/vespa_podcast_ranking). But in this blog post we will focuse on how to make the result of search more clear to user why id matched and on what. This is called[Dynamic snippets](https://docs.vespa.ai/en/document-summaries.html#dynamic-snippets) it is an awesome feature that is not that well explained in the docs so here we go. 

If we use a previous query as an example: 



```bash
vespa query \
'yql=select title, description from podcast where title contains "Vespa AI Search" or description contains "Vespa AI search"' \
'hits=10' \
'ranking=podcast-search' \
'input.query(q)=100'
```

However I realised one thing whild working on this: 

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

compare to: 


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

We get a very different result. The interesting parts are: 

1) We get different results
2) The first query gives a lot lower ranking score, why did not the correct one get chosen as in the second example. 


If we look futher we can set an important difference if we compare the coverage. When we add the contains AI term the search becomes [`degrade`](https://docs.vespa.ai/en/graceful-degradation.html#match-phase-degradation) due to the [match-phase](https://docs.vespa.ai/en/reference/schema-reference.html#match-phase): 

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

This is due to that the Vespa estimates about how many rows will be returned when we have a `match-phase` set in the schema and if more then the hit limit we will use the `match-phase` selection criteria to filter the results. In our case this is the difference between the querys, when we add that `or contains "AI"` the estimation is that we will hit a lot more and thus we start to filter based upon recency ranking. Since more data is small and static we will remove this criteria. Futher allways watch out for the degraded criteria. 