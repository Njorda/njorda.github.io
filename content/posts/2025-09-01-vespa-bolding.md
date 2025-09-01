---
layout: post
title: "Vespa Search Engine: Bolding and Dynamic Snippets"
subtitle: "How to visualize matches for users"
date: 2025-09-01
author: "Niklas Hansson"
URL: "/2025/09/01/vespa_highlights_bolding"
---

This post is a follow-up to the [previous post](/2025/08/23/vespa_search_podcast_index), which includes instructions for setting up the cluster if you want to try this yourself.

Great search results don’t just rank the right documents—they also show users why a result matches. Vespa can both:
- **Highlight matches (“bolding”)**: wrap matching query terms in `<hi>` tags.
- **Create dynamic snippets**: extract fragments around matches to show relevant context, also with `<hi>` tags for highlights.

In Vespa, [dynamic snippets](https://docs.vespa.ai/en/document-summaries.html#dynamic-snippets) are generated from the source field by extracting the most relevant fragments around matching terms. The XML tags `<hi>` mark the matched tokens.

Dynamic summaries and bolding are only supported for fields fetched as `index` and/or `summary` values. They do not work for values fetched from in-memory `attribute`s. If you add dynamic snippets or bolding to an attribute-backed summary field, you’ll see a warning like:

> Will fetch the disk summary value of field 'description' in document summaries [default] since this summary field uses a dynamic summary value (snippet/bolding): Dynamic summaries and bolding is not supported with summary values fetched from in-memory attributes yet. If you want to see partial updates to this attribute, remove any bolding and dynamic snippeting from this field.

To avoid changing the behavior of your existing fast, attribute-only summary, define a new document summary that inherits the minimal one and adds generated fields for snippets/bolding.

### Add a document summary for snippets and bolding

```
document-summary bolded inherits minimal {
    summary title_bold {
        source: title
        bolding: on
    }
    summary description_dynamic {
        source: description
        dynamic
    }
}
```

We can test this with: 


```bash
vespa query \
'yql=select * from podcast where title contains "Vespa" or description contains "RAG"' \
'hits=1' \
'ranking=podcast-search' \
'input.query(q)=100' \
'language=en' \
'presentation.summary=bolded'
```

which returns: 


```json
{
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 506
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
                "relevance": 14.373663603488678,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 9.149635037339992,
                        "bm25(title)": 12.543736596020679,
                        "query(q)": 100.0,
                        "freshness": 7.224947569518041E-9
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "<sep />of agentic AI, search architecture, retrieval-augmented generation (<hi>RAG</hi>), and scalable enterprise applications. Whether you're a CTO<sep />",
                    "titl_bold": "<hi>Vespa</hi> Voice",
                    "title": "Vespa Voice",
                    "description": "Welcome to Vespa Voice, the podcast where AI leaders, search pioneers, and enterprise innovators converge. Each episode dives deep into the evolving landscape of AI, featuring candid conversations with experts shaping the future of agentic AI, search architecture, retrieval-augmented generation (RAG), and scalable enterprise applications. Whether you're a CTO driving digital transformation, a CIO reimagining data strategy, or an engineer building next-gen ML and search systems, this is your signal for what's next in intelligent infrastructure.",
                    "newest_item_pubdate": 1745337276,
                    "url": "https://feeds.blubrry.com/feeds/3953602.xml",
                    "popularity_score": 0.0
                }
            }
        ]
    }
}
```

When to use bolding vs snippets

- **Titles (short fields)**: use bolding only. Snippets add little value when the field is already concise.
- **Descriptions (long, verbose fields)**: use dynamic snippets to show only the most relevant parts—similar to how web search engines present results.

You can tune how snippets are generated (window size, max occurrences, separators, etc.) in [`services.xml`](https://docs.vespa.ai/en/reference/services.html). See the [dynamic snippet configuration](https://docs.vespa.ai/en/document-summaries.html#dynamic-snippet-configuration). Defaults typically work well to start.

