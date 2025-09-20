---
layout: post
title: "Vespa(Search engine): podcast index"
subtitle: "How to load a SQLite db to Vespa search engine"
date: 2025-08-23
author: "Niklas Hansson"
URL: "/2025/08/23/vespa_search_podcast_index"
---

After playing around with the [Podcast index](https://podcastindex.org/apps?appTypes=app) and the amazing [SQLite dump](https://github.com/Podcastindex-org/database), I got stuck and annoyed that the search experience using SQLite and [FTS5](https://www.sqlite.org/fts5.html) is not that nice. I never found a good way to do a fuzzy search for podcasts. After being stuck, I decided that I want to use Vespa instead, which I had played around with before and liked. 

I will not go over Vespa but instead focus on how to: 

- move a SQLite dump to Vespa
- set up the Vespa schemas
- set up the Vespa rankings( we will focuse further on this in follow up posts)


# Set up Vespa

So first, the schema of the SQLite db is as follows; there is only one table called podcast: 

```sql
// CREATE TABLE podcasts (
//     id INTEGER PRIMARY KEY,
//     url TEXT NOT NULL UNIQUE,
//     title TEXT NOT NULL,
//     lastUpdate INTEGER,
//     link TEXT NOT NULL,
//     lastHttpStatus INTEGER,
//     dead INTEGER,
//     contentType TEXT NOT NULL,
//     itunesId INTEGER,
//     originalUrl TEXT NOT NULL,
//     itunesAuthor TEXT NOT NULL,
//     itunesOwnerName TEXT NOT NULL,
//     explicit INTEGER,
//     imageUrl TEXT NOT NULL,
//     itunesType TEXT NOT NULL,
//     generator TEXT NOT NULL,
//     newestItemPubdate INTEGER,
//     language TEXT NOT NULL,
//     oldestItemPubdate INTEGER,
//     episodeCount INTEGER,
//     popularityScore INTEGER,
//     priority INTEGER,
//     createdOn INTEGER,
//     updateFrequency INTEGER,
//     chash TEXT NOT NULL,
//     host TEXT NOT NULL,
//     newestEnclosureUrl TEXT NOT NULL,
//     podcastGuid TEXT NOT NULL,
//     description TEXT NOT NULL,
//     category1 TEXT NOT NULL,
//     category2 TEXT NOT NULL,
//     category3 TEXT NOT NULL,
//     category4 TEXT NOT NULL,
//     category5 TEXT NOT NULL,
//     category6 TEXT NOT NULL,
//     category7 TEXT NOT NULL,
//     category8 TEXT NOT NULL,
//     category9 TEXT NOT NULL,
//     category10 TEXT NOT NULL,
//     newestEnclosureDuration INTEGER
// );
```


To load it up after downloading it from [here](https://github.com/Podcastindex-org/database#database), open it up using:


```bash
sqlite3 --readonly podcastindex.db
.schema
```

I will use Go (just because I like it) for this, but it should be easy to ask your favorite LLM to translate it to any programming language.


Vespa uses a typed schema definition. If you feel we're glossing over the basics of Vespa, that is indeed the case. Vespa has [good documentation](https://docs.vespa.ai), even though it is sometimes hard to find the right content. What we will focus on though is how we construct our search, how we set up [our schema](https://docs.vespa.ai/en/schemas.html), and so on. 


I will not take all the fields in the SQLite db schema but just the ones I care about and intend to keep, in order to reduce size and simplify things as much as possible.

```bash
id INTEGER PRIMARY KEY,
//     url TEXT NOT NULL UNIQUE,
//     title TEXT NOT NULL,
//     lastUpdate INTEGER,
//     link TEXT NOT NULL,
//     lastHttpStatus INTEGER,
//     dead INTEGER,
//     contentType TEXT NOT NULL,
//     itunesId INTEGER,
//     originalUrl TEXT NOT NULL,
//     itunesAuthor TEXT NOT NULL,
//     itunesOwnerName TEXT NOT NULL,
//     explicit INTEGER,
//     imageUrl TEXT NOT NULL,
//     itunesType TEXT NOT NULL,
//     generator TEXT NOT NULL,
//     newestItemPubdate INTEGER,
//     language TEXT NOT NULL,
//     oldestItemPubdate INTEGER,
//     episodeCount INTEGER,
//     popularityScore INTEGER,
//     priority INTEGER,
//     createdOn INTEGER,
//     updateFrequency INTEGER,
//     chash TEXT NOT NULL,
//     host TEXT NOT NULL,
//     newestEnclosureUrl TEXT NOT NULL,
//     podcastGuid TEXT NOT NULL,
//     description TEXT NOT NULL,
```


Based upon the fields in the SQLite db, this is the initial schema we kick it off with: 


```
schema podcast {
    document podcast {


        # Since this is the id I want to have it in memory and thus use attribute, also make it available 
        # for operations such as sort filtering and ranking.
        # I want it indexed and thus I also set index
        # Since it is an ID I want exact matching, I make the exact terminator
        # explicit but the default is actually the same.
        # Attribute link: https://docs.vespa.ai/en/attributes.html
        field id type string {
            indexing: attribute | summary 
            match {
                exact
                exact-terminator: "@@"
            }
        }
  
        field url type string {
            indexing: summary 
        }

        field title type string {
            indexing: attribute | summary | index
            match: text
        }

        field last_update type long {
            indexing: summary
        }

        field link type string {
            indexing: summary
        }

        field last_http_status type int {}
        field dead type bool {
            indexing: attribute
            rank: filter
        }
        field content_type type string {}
        field original_url type string {}
        field itunes_author type string {}
        field itunes_owner_name type string {}
        field explicit type bool {}

        field image_url type string {
            indexing: summary
        }

        field itunes_type type string {}
        field generator type string {}
        field newest_item_pubdate type long {
            indexing: attribute
            attribute: fast-search
        }
        field language type string {}
        field oldest_item_pubdate type long {}
        field episode_count type int {}
        field popularity_score type double {
            indexing: summary | attribute 
        }
        field priority type int {}
        field created_on type long {}
        field update_frequency type int {}
        field chash type string {}
        field podcast_host type string {}
        field newest_enclosure_url type string {}
        field podcast_guid type string {}

        field description type string {
            indexing: summary | attribute | index
            attribute: fast-search
            match: text
        }
    }

    # Vespa seems to first do a filter on the fields that we rank on
    # we can not have the filtering without the ranking it seems like? 
    fieldset prefix {
        fields: title_array, title_prefix
    }

    field title_prefix type string {
        indexing: input title | attribute | summary
    }

    # This will be used for prefix search in on the podcast name 
    # https://docs.vespa.ai/en/text-matching.html#prefix-match
    field title_array type array<string> {
        indexing: input title | trim | split " +" | attribute | summary
        stemming: best
        attribute: fast-search
    }

    # Document Summaries are memory-only operations if all fields are attributes.
    document-summary minimal {
        summary title {}
        summary description {}
        summary newest_item_pubdate {}
        summary url {}
        summary popularity_score {}
    }


    # The match phase is dont on the content nodes, we rank the content based upon recency. 
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
        # Figure out how to use the chunks in the first phase
        # https://pyvespa.readthedocs.io/en/latest/examples/multilingual-multi-vector-reps-with-cohere-cloud.html
        first-phase {
            expression: bm25(title) + bm25(description) * 0.2 + freshness()
        }

        match-features {
            bm25(title)
            bm25(description)
            freshness()
            query(q)
        }  
    }
}

```

The schema contains some inline comments, but let's go over it a bit more. See the Vespa [Schema reference](https://docs.vespa.ai/en/reference/schema-reference.html). We only selected a subset of the fields. My first impression when using Vespa is that it is more low-level than for example Postgres in the way the schemas are set up. It is also important to notice that Vespa is not a relational database or a document store but a search engine, and most of the internals are focused around search. 

The first field we define is the ID, which we keep as a string. Important to notice is that we explicitly set the matching condition to exact since we don't want to do search over the IDs (but actual matches). We also set the field to summary and attribute. This was one of the most unclear things in the schema definition for me, but the best explanation can be found [here](https://docs.vespa.ai/en/schemas.html#indexing) in the table below. Further reading can also be found [here](https://docs.vespa.ai/en/reference/schema-reference.html#indexing): 


>| EXPRESSION | DESCRIPTION |
>| --- | --- |
>| attribute | For structured data: Keep this field in memory in a forward structure. This makes the field available for grouping, sorting and ranking. Attributes may also be searched by complete match (word or exact), or (for numerical fields) by range. Optionally a B-tree in memory can also be created by adding the fast-search option - this improves performance if the attribute is a strong criterion in queries (i.e. filters out many documents). Reference / attribute details |
>| index | For unstructured text: Create a text index for this field. Text matching and all text ranking features become available. Indexes are disk backed and do not need to fit in memory. Reference / index details |
>| summary | Include this field in the document summary in search result sets. Reference / document summary details |

The difference can also be expressed as: 
- attribute indexing is used for database-style string matching <strong>without linguistic processing and where the exact string contents are matched
- index indexing is used with linguistic processing such as tokenization and stemming.

So if you want to use BM25 and similar ranking functions for a field, you will need to add index. This is not of interest for the ID field where we have structured SQL-like queries only.

We set up the id field as follows: 

```
    field id type string {
        indexing: attribute | summary
        match {
            exact
            exact-terminator: "@@"
        }
    }
```

Since we want to filter on it, we keep it as an attribute (an attribute field will be in memory at all times) and have it as a summary field in order to return it.

We than make a small optimization but one that potentially could reduce the CPU load a lot and set the [`rank: filter`](https://docs.vespa.ai/en/reference/schema-reference.html#rank): 

```
    field dead type bool {
        indexing: attribute
        rank: filter
    }
```

This creates a bitvector field that can be used for searching. If other indexes would be available, the query planner would select the best one in order to optimize the query performance. 


The next interesting part is the `title_array`,where we create a [synthetic field](https://docs.vespa.ai/en/reference/indexing-language-reference.html): 

```
    field title_array type array<string> {
        indexing: input title | trim | split " +" | attribute | summary
        stemming: best
        attribute: fast-search
    }
```

We take the column `title`, trim it, then split it. This is done in order to do prefix search on individual terms in the title. I will later show how we do this and also handle spelling errors. The field is not set as an index field, due to the fact that regular index fields do [not support prefix matching](https://docs.vespa.ai/en/text-matching.html).  By default, for attribute fields no index is added. By adding `attribute: fast-search`, a B-tree-like structured index is added to the attribute which [can significantly speed up search](https://docs.vespa.ai/en/performance/feature-tuning.html#when-to-use-fast-search-for-attribute-fields). 


In order to also allow for matching against a title with multiple words, the title array will only work for each word individually and we cannot use the title field since it is both an attribute and an index field. We create a synthetic field from the `title` field called `title_prefix`: 

```
field title_prefix type string {
    indexing: input title | attribute | summary
}

```

We then create a [`document summary`](https://docs.vespa.ai/en/document-summaries.html): 

```
    document-summary minimal {
        summary title {}
        summary description {}
        summary newest_item_pubdate {}
        summary url {}
        summary popularity_score {}
    }
```

This is used as a grouping for the returned fields, the summary, and is a logical grouping of fields we like to get out. Since this is a search engine, there is a distinction between fields used internally and exposed externally. 

> Key to understanding Vespa is to understand what it is not: a relational database, but rather a search engine. 


Next up is probably the most interesting thing, which is the ranking functions. Before we look at the ranking function we build, let's spend some time looking at search as a ranking problem and how we break it down into different phases. We will do this through the lens of Vespa, but the concepts hold outside. In general, it is a pipeline with less and less data but potentially more compute spent per sample in each step in order to get a better ranking. 


1) Matching based upon the incoming filter requirements, such as IDs or Language. These are strict requirements that need to be fulfilled.
2) First reranking - this is done on the content nodes ([more about Vespa setup](https://docs.vespa.ai/en/content/content-nodes.html) can be found here. But in short, we have stateful content nodes storing the data and stateless compute nodes).
3) Second reranking (also done on the content nodes).
4) Global reranking, combining the tuples from all content nodes.

[This is a good resource](https://docs.vespa.ai/en/phased-ranking.html) to read more about it. All the options for the rank-profile can be found [here](https://docs.vespa.ai/en/reference/schema-reference.html#rank-profile).

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
        expression: bm25(title) + bm25(description) * 0.2 + freshness()
    }

    match-features {
        bm25(title)
        bm25(description)
        freshness()
        query(q)
    }  
}
```

We will dive more in to the ranking in a follow up post.


To get vespa setup and running we will use docker: 

```
docker run --detach --name vespa-podcast --hostname vespa-podcast \
  --publish 127.0.0.1:8080:8080 --publish 127.0.0.1:19112:19112 --publish 127.0.0.1:19071:19071 \
  vespaengine/vespa
```

We assume the following folder structure: 

search/
├── schemas/
│   ├── podcast.sd
└── services.xml


Where the service.xml file looks like the following (use it as an example for recreating this blog post):

```xml
<?xml version='1.0' encoding='UTF-8'?>
<services version="1.0" xmlns:deploy="vespa" xmlns:preprocess="properties">

  <container id="search" version="1.0">
    <document-api/>
    <search>
    </search>
    <document-processing/>
    <nodes>
      <node hostalias="node1" />
    </nodes>
  </container>

  <content id="podcast" version="1.0">
    <min-redundancy>2</min-redundancy>
    <documents>
      <document type='podcast' mode="index"/>
    </documents>
    <nodes>
      <node distribution-key='0' hostalias='node1'/>
    </nodes>
  </content>
  
</services>
```

You can then use the [Vespa CLI](https://docs.vespa.ai/en/vespa-cli.html#install-configure-and-run) from the search folder: 

```
vespa deploy
```

# Dump the SQLite database data

Feel free to dump the SQLite db any way you want. Below is a small Go script to do it. 

```go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Podcast struct {
	ID                 string  `json:"id"`
	URL                string  `json:"url"`
	Title              string  `json:"title"`
	LastUpdate         int64   `json:"last_update"`
	Link               string  `json:"link"`
	LastHTTPStatus     int     `json:"last_http_status"`
	Dead               bool    `json:"dead"`
	ContentType        string  `json:"content_type"`
	OriginalURL        string  `json:"original_url"`
	ITunesAuthor       string  `json:"itunes_author"`
	ITunesOwnerName    string  `json:"itunes_owner_name"`
	Explicit           bool    `json:"explicit"`
	ImageURL           string  `json:"image_url"`
	ITunesType         string  `json:"itunes_type"`
	Generator          string  `json:"generator"`
	NewestItemPubDate  int64   `json:"newest_item_pubdate"`
	Language           string  `json:"language"`
	OldestItemPubDate  int64   `json:"oldest_item_pubdate"`
	EpisodeCount       int     `json:"episode_count"`
	PopularityScore    float64 `json:"popularity_score"`
	Priority           int     `json:"priority"`
	CreatedOn          int64   `json:"created_on"`
	UpdateFrequency    int     `json:"update_frequency"`
	Chash              string  `json:"chash"`
	PodcastHost        string  `json:"podcast_host"`
	NewestEnclosureURL string  `json:"newest_enclosure_url"`
	PodcastGUID        string  `json:"podcast_guid"`
	Description        string  `json:"description"`
}

type VespaPut struct {
	Put    string  `json:"put"`
	Fields Podcast `json:"fields"`
}

func main() {
	fmt.Println(len(os.Args))
	if len(os.Args) != 4 {
		log.Fatal(" add 3 args, usage add path to db file and out put path")
	}
	dbPath := os.Args[1]
	outputPath := os.Args[2]
	nbrOfPodcasts, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal("strconv: ", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	podcasts, err := fetchPodcasts(db, nbrOfPodcasts)
	if err != nil {
		log.Fatal("fetchPodcasts: ", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("create: ", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	for _, p := range podcasts {
		putID := fmt.Sprintf("id:podcast:podcast::%s", p.ID)
		record := VespaPut{Put: putID, Fields: p}
		if err := enc.Encode(record); err != nil {
			log.Fatal("encode: ", err)
		}
	}

}

func fetchPodcasts(db *sql.DB, limit int) ([]Podcast, error) {
	rows, err := db.Query(`
		SELECT
			id,
			url,
			title,
			lastUpdate,
			link,
			lastHttpStatus,
			dead,
			contentType,
			originalUrl,
			itunesAuthor,
			itunesOwnerName,
			explicit,
			imageUrl,
			itunesType,
			generator,
			newestItemPubdate,
			language,
			oldestItemPubdate,
			episodeCount,
			popularityScore,
			priority,
			createdOn,
			updateFrequency,
			chash,
			host,
			newestEnclosureUrl,
			podcastGuid,
			description
		FROM podcasts
		ORDER BY id
		LIMIT ?;
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Podcast
	for rows.Next() {
		var (
			idInt              int64
			url                string
			title              string
			lastUpdate         int64
			link               string
			lastHTTPStatus     int
			deadInt            int
			contentType        string
			originalURL        string
			itunesAuthor       string
			itunesOwnerName    string
			explicitInt        int
			imageURL           string
			itunesType         string
			generator          string
			newestItemPubDate  int64
			language           string
			oldestItemPubDate  int64
			episodeCount       int
			popularityScoreInt int
			priority           int
			createdOn          int64
			updateFrequency    int
			chash              string
			host               string
			newestEnclosureURL string
			podcastGUID        string
			description        string
		)

		if err := rows.Scan(
			&idInt,
			&url,
			&title,
			&lastUpdate,
			&link,
			&lastHTTPStatus,
			&deadInt,
			&contentType,
			&originalURL,
			&itunesAuthor,
			&itunesOwnerName,
			&explicitInt,
			&imageURL,
			&itunesType,
			&generator,
			&newestItemPubDate,
			&language,
			&oldestItemPubDate,
			&episodeCount,
			&popularityScoreInt,
			&priority,
			&createdOn,
			&updateFrequency,
			&chash,
			&host,
			&newestEnclosureURL,
			&podcastGUID,
			&description,
		); err != nil {
			return nil, err
		}

		out = append(out, Podcast{
			ID:                 strconv.FormatInt(idInt, 10),
			URL:                url,
			Title:              title,
			LastUpdate:         lastUpdate,
			Link:               link,
			LastHTTPStatus:     lastHTTPStatus,
			Dead:               deadInt != 0,
			ContentType:        contentType,
			OriginalURL:        originalURL,
			ITunesAuthor:       itunesAuthor,
			ITunesOwnerName:    itunesOwnerName,
			Explicit:           explicitInt != 0,
			ImageURL:           imageURL,
			ITunesType:         itunesType,
			Generator:          generator,
			NewestItemPubDate:  newestItemPubDate,
			Language:           language,
			OldestItemPubDate:  oldestItemPubDate,
			EpisodeCount:       episodeCount,
			PopularityScore:    float64(popularityScoreInt),
			Priority:           priority,
			CreatedOn:          createdOn,
			UpdateFrequency:    updateFrequency,
			Chash:              chash,
			PodcastHost:        host,
			NewestEnclosureURL: newestEnclosureURL,
			PodcastGUID:        podcastGUID,
			Description:        description,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

```

The script is slightly longer in order to remap the name of the fields (could also be done by renaming them in the query). Key things to notice:

- The output is JSONL
- Each row has `{"put":"id:podcast:podcast::6","fields":{ .... }}` where fields is the actual data of value we care about. The put field contains the id, the schema and the document as part of it. 


To run it: 

```
go run . /Users/niklashansson/Documents/github/pod/backend/podcastindex_feeds.db dump.jsonl 1000000
```

After this I end up with a file of `5.8GB` 

# Ingest the data into Vespa

```
vespa feed --dump.jsonl
```

The ingest will take a little bit of time so let it run. I kind of felt it is slow but it might be due to the Docker setup and how the CLI ingests it, since it is roughly 4 million podcasts (only). 

This returned:  

```bash
feed: got status 400 ({"pathId":"/document/v1/podcast/podcast/docid/1825377","message":"In document 'id:podcast:podcast::1825377': Could not parse field 'description' of type string: The string field value contains illegal code point 0x3"}) for put id:podcast:podcast::1825377: not retryable
feed: got status 400 ({"pathId":"/document/v1/podcast/podcast/docid/2783072","message":"In document 'id:podcast:podcast::2783072': Could not parse field 'itunes_author' of type string: The string field value contains illegal code point 0x3"}) for put id:podcast:podcast::2783072: not retryable
feed: got status 400 ({"pathId":"/document/v1/podcast/podcast/docid/5535724","message":"In document 'id:podcast:podcast::5535724': Could not parse field 'itunes_author' of type string: The string field value contains illegal code point 0x3"}) for put id:podcast:podcast::5535724: not retryable
feed: got status 400 ({"pathId":"/document/v1/podcast/podcast/docid/6420537","message":"In document 'id:podcast:podcast::6420537': Could not parse field 'itunes_author' of type string: The string field value contains illegal code point 0x8"}) for put id:podcast:podcast::6420537: not retryable
{
  "feeder.operation.count": 4625512,
  "feeder.seconds": 308.485,
  "feeder.ok.count": 4625508,
  "feeder.ok.rate": 14994.269,
  "feeder.error.count": 0,
  "feeder.inflight.count": 0,
  "http.request.count": 4625512,
  "http.request.bytes": 3448131187,
  "http.request.MBps": 11.178,
  "http.exception.count": 0,
  "http.response.count": 4625512,
  "http.response.bytes": 414304888,
  "http.response.MBps": 1.343,
  "http.response.error.count": 4,
  "http.response.latency.millis.min": 1,
  "http.response.latency.millis.avg": 34,
  "http.response.latency.millis.max": 191,
  "http.response.code.counts": {
    "200": 4625508,
    "400": 4
  }
}
```

Indicating that we failed on 4 podcasts, this will not matter for this blog post so we will just skip ahead.

# Query our Vespa search engine

The easiest way to query Vespa is to use the [Vespa CLI](https://docs.vespa.ai/en/vespa-cli.html#queries)


Let's start by trying the prefix matching. We will use [Lenny's Podcast](https://www.lennysnewsletter.com/podcast) as our example. 


As a first example, we can try to match against words within the podcast name, thus `Lenny's` and `Podcast`. 

```bash
vespa query 'select title from podcast where title_array contains ({maxEditDistance: 4, prefix: true} "Lenny'\''s") and id contains "1539262"'
```


```bash
vespa query 'select title from podcast where title_array contains ({maxEditDistance: 4, prefix: true} "Podcast") and id contains "1539262"'
```

However, if we want to match against a prefix that contains both, it will not work to use the title_array field since we match against an array where we split it. 

```bash
 vespa query 'select title from podcast where title_array contains ({maxEditDistance: 4, prefix: true} "Lenny'\''s Podcast") and id contains "1539262"'
```

A good learning here is that the field `title_array` will not help if the user has written `Lenny's Pod` since the title is split. Thus we want to actually match on both the split and the full ones in order to get the best possible hits (of course depending upon your use case, but for my case at least). 

Instead we then need to do: 

```bash
vespa query 'select title from podcast where prefix contains ({maxEditDistance: 4, prefix: true} "Lenny'\''s Podcast") and id contains "1539262"'
```

This works since it uses the `fieldset` for the prefix search which is a grouping for both the `title_array` and the `title_prefix`. 

That is all for this post. In a follow-up, we will dive deeper into the ranking function. 



