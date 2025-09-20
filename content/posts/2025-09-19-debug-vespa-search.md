---
layout: post
title: "Debug vespa search"
subtitle: "Fix slow prefix search"
date: 2025-09-19
author: "Niklas Hansson"
URL: "/2025/09/01/vespa_highlights_bolding"
---

This post is a follow-up to the [previous post](/2025/08/16/vespa-odcast-search), which includes instructions for setting up the cluster if you want to try this yourself.

As often happens you first make it work and then you need to make it fast or good or what ever. In our case the prefix search is slow, as in really slow. 

The following is the prefix searh discussed [here](2025/08/16/vespa-odcast-search):

```bash
vespa query \
  'yql=select id, title from podcast where prefix contains ({maxEditDistance: 4, prefix: true} "All-in")' \
  'hits=5' \
  'timeout=10s' \
  'presentation.summary=minimal' \
  'timing=true' \
  'trace.timestamps=true'
{
    "timing": {
        "querytime": 0.335,
        "summaryfetchtime": 0.001,
        "searchtime": 0.337
    },
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 61
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
                "id": "index:podcast/0/f7da8dee1354896e16fde095",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "24325",
                    "title": "ALL-IN BALLIN’"
                }
            },
            {
                "id": "index:podcast/0/bd0403aab7a4b137f873978e",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "1015378",
                    "title": "All-In with Chamath, Jason, Sacks & Friedberg"
                }
            },
            {
                "id": "index:podcast/0/1edffbba6f79b8df236696f9",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "2481852",
                    "title": "all-in.de audionews"
                }
            },
            {
                "id": "index:podcast/0/fee423f36823ad385e4f6b69",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "2991462",
                    "title": "All-Inclusive Sports Podcast"
                }
            },
            {
                "id": "index:podcast/0/21d4ea261a288b67af28e6bc",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "3374917",
                    "title": "All-In Cast"
                }
            }
        ]
    }
}
```

It takes roughly 300 ms, which is a lot more than I expected. Especially compared to a full text search which is 70 ms: 

```bash
vespa query \
  'yql=select * from podcast where userQuery()' \
  'query=All-in with Chamath' \
  'ranking=podcast-search' \
  'default-index=content_search' \
  'timing=true'
{
    "timing": {
        "querytime": 0.07100000000000001,
        "summaryfetchtime": 0.002,
        "searchtime": 0.074
    },
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 120578
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
                "id": "id:podcast:podcast::1015378",
                "relevance": 21.461054369895148,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 15.38908424942882,
                        "bm25(title)": 18.383237520009384,
                        "query(q)": 0.0,
                        "freshness": 1.0027671093758963
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "<sep />degenerate gamblers & besties <hi>Chamath</hi> Palihapitiya, Jason<sep />David Friedberg cover <hi>all</hi> things economic, tech, political<sep />",
                    "titl_bold": "<hi>All</hi>-<hi>In</hi> <hi>with</hi> <hi>Chamath</hi>, Jason, Sacks & Friedberg",
                    "documentid": "id:podcast:podcast::1015378",
                    "title_prefix": "All-In with Chamath, Jason, Sacks & Friedberg",
                    "title_array": [
                        "All-In",
                        "with",
                        "Chamath,",
                        "Jason,",
                        "Sacks",
                        "&",
                        "Friedberg"
                    ],
                    "id": "1015378",
                    "url": "https://allinchamathjason.libsyn.com/rss",
                    "title": "All-In with Chamath, Jason, Sacks & Friedberg",
                    "last_update": 1754781504,
                    "link": "http://allinchamathjason.libsyn.com/website",
                    "image_url": "https://static.libsyn.com/p/assets/a/9/c/b/a9cb4d1dadb1ea21/all-in_logo.png",
                    "popularity_score": 9.0,
                    "description": "Industry veterans, degenerate gamblers & besties Chamath Palihapitiya, Jason Calacanis, David Sacks & David Friedberg cover all things economic, tech, political, social & poker."
                }
            },
            {
                "id": "id:podcast:podcast::2185926",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 5.1787157709201365E-136
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "More on the division happening <hi>in</hi> the Methodist denomination and  Interview <hi>with</hi> Nina Roesner author of \"<hi>With</hi> <hi>All</hi> Due Respect\"",
                    "titl_bold": "More on the division happening <hi>in</hi> the Methodist denomination and  Interview <hi>with</hi> Nina Roesner author of \"<hi>With</hi> <hi>All</hi> Due Respect\"",
                    "documentid": "id:podcast:podcast::2185926",
                    "title_prefix": "More on the division happening in the Methodist denomination and  Interview with Nina Roesner author of \"With All Due Respect\"",
                    "title_array": [
                        "More",
                        "on",
                        "the",
                        "division",
                        "happening",
                        "in",
                        "the",
                        "Methodist",
                        "denomination",
                        "and",
                        "Interview",
                        "with",
                        "Nina",
                        "Roesner",
                        "author",
                        "of",
                        "\"With",
                        "All",
                        "Due",
                        "Respect\""
                    ],
                    "id": "2185926",
                    "url": "https://afr.net/api/v2/podcast/podcastsbyidrss/?id=43798",
                    "title": "More on the division happening in the Methodist denomination and  Interview with Nina Roesner author of \"With All Due Respect\"",
                    "last_update": 1617458773,
                    "link": "https://afr.net/podcasts/airing-the-addisons/2019/october/more-on-the-division-happening-in-the-methodist-denomination-and-interview-with-nina-roesner-author-of-with-all-due-respect/",
                    "image_url": "https://afr.net/media/4719/ata_800x500.png",
                    "popularity_score": 0.0,
                    "description": "More on the division happening in the Methodist denomination and  Interview with Nina Roesner author of \"With All Due Respect\""
                }
            },
            {
                "id": "id:podcast:podcast::2215057",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 2.30113918712327E-198
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "Pastor Joseph talks <hi>with</hi> Mound Bayou Mississippi<sep />Johnson about being a disciple of Christ <hi>in</hi> <hi>all</hi> aspects of life",
                    "titl_bold": "Pastor Joseph talks <hi>with</hi> Mound Bayou Mississippi Mayor and Pastor Darryl Johnson about being a disciple of Christ <hi>in</hi> <hi>all</hi> aspects of life",
                    "documentid": "id:podcast:podcast::2215057",
                    "title_prefix": "Pastor Joseph talks with Mound Bayou Mississippi Mayor and Pastor Darryl Johnson about being a disciple of Christ in all aspects of life",
                    "title_array": [
                        "Pastor",
                        "Joseph",
                        "talks",
                        "with",
                        "Mound",
                        "Bayou",
                        "Mississippi",
                        "Mayor",
                        "and",
                        "Pastor",
                        "Darryl",
                        "Johnson",
                        "about",
                        "being",
                        "a",
                        "disciple",
                        "of",
                        "Christ",
                        "in",
                        "all",
                        "aspects",
                        "of",
                        "life"
                    ],
                    "id": "2215057",
                    "url": "https://afr.net/api/v2/podcast/podcastsbyidrss/?id=38644",
                    "title": "Pastor Joseph talks with Mound Bayou Mississippi Mayor and Pastor Darryl Johnson about being a disciple of Christ in all aspects of life",
                    "last_update": 1617500586,
                    "link": "https://afr.net/podcasts/the-hour-of-intercession/2016/december/pastor-joseph-talks-with-mound-bayou-mississippi-mayor-and-pastor-darryl-johnson-about-being-a-disciple-of-christ-in-all-aspects-of-life/",
                    "image_url": "https://afr.net/media/4713/thoi_800x500.jpg",
                    "popularity_score": 0.0,
                    "description": "Pastor Joseph talks with Mound Bayou Mississippi Mayor and Pastor Darryl Johnson about being a disciple of Christ in all aspects of life"
                }
            },
            {
                "id": "id:podcast:podcast::2223605",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 2.25264933250383E-117
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "<hi>All</hi> <hi>In</hi> Healing <hi>with</hi> Kimberlie Carlson<sep />space <hi>in</hi> our body through inflammation, rooted <hi>in</hi> low vibrational emotions<sep />",
                    "titl_bold": "<hi>All</hi> <hi>In</hi> Healing <hi>with</hi> Kimberlie Carlson",
                    "documentid": "id:podcast:podcast::2223605",
                    "title_prefix": "All In Healing with Kimberlie Carlson",
                    "title_array": [
                        "All",
                        "In",
                        "Healing",
                        "with",
                        "Kimberlie",
                        "Carlson"
                    ],
                    "id": "2223605",
                    "url": "https://www.spreaker.com/show/4207963/episodes/feed",
                    "title": "All In Healing with Kimberlie Carlson",
                    "last_update": 1754435128,
                    "link": "https://kimberliecarlson.com/",
                    "image_url": "https://d3wo5wojvuv7l.cloudfront.net/t_rss_itunes_square_1400/images.spreaker.com/original/377b10ec43267b9d0c6cf5f55658380d.jpg",
                    "popularity_score": 0.0,
                    "description": "All In Healing with Kimberlie Carlson: Health ~ Harmony ~ Happiness<br /><br />Begin your journey to total healing with All In Healing with Kimberlie Carlson. Kimberlie’s deepest passion in life is to set free as many people as possible from physical and emotional patterns of struggle, addiction, and dis-ease. These things take up space in our body through inflammation, rooted in low vibrational emotions. Your continued suffering is rooted deep within you, a resistance to higher vibrations. It takes courage to continue deeper and deeper until the root is discovered and every aspect of an issue healed. Once caught in a never-ending cycle of pain, struggle, and isolation, Kimberlie discovered new and unbelievably thorough methods of dealing with allergies, fatigue, migraines, and so much more.<br /><br />Learn about YOUR layers and how to work through them to banish hay fever, chronic infections, and brain fog for good!"
                }
            },
            {
                "id": "id:podcast:podcast::2225174",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 5.373970568581307E-189
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "<sep /> and Jerry's big push for <hi>all</hi> things homosexual <hi>in</hi> Australia and <hi>in</hi> the U.S.A. Man living <hi>with</hi> two women raising a family",
                    "titl_bold": "HR2 - Tiger Woods arrested on DUI charges, Ben and Jerry's big push for <hi>all</hi> things homosexual <hi>in</hi> Australia and <hi>in</hi> the U.S.A. Man living <hi>with</hi> two women raising a family",
                    "documentid": "id:podcast:podcast::2225174",
                    "title_prefix": "HR2 - Tiger Woods arrested on DUI charges, Ben and Jerry's big push for all things homosexual in Australia and in the U.S.A. Man living with two women raising a family",
                    "title_array": [
                        "HR2",
                        "-",
                        "Tiger",
                        "Woods",
                        "arrested",
                        "on",
                        "DUI",
                        "charges,",
                        "Ben",
                        "and",
                        "Jerry's",
                        "big",
                        "push",
                        "for",
                        "all",
                        "things",
                        "homosexual",
                        "in",
                        "Australia",
                        "and",
                        "in",
                        "the",
                        "U.S.A.",
                        "Man",
                        "living",
                        "with",
                        "two",
                        "women",
                        "raising",
                        "a",
                        "family"
                    ],
                    "id": "2225174",
                    "url": "https://afr.net/api/v2/podcast/podcastsbyidrss/?id=36615",
                    "title": "HR2 - Tiger Woods arrested on DUI charges, Ben and Jerry's big push for all things homosexual in Australia and in the U.S.A. Man living with two women raising a family",
                    "last_update": 1617513550,
                    "link": "https://afr.net/podcasts/airing-the-addisons/2017/may/hr2-tiger-woods-arrested-on-dui-charges-ben-and-jerry-s-big-push-for-all-things-homosexual-in-australia-and-in-the-u-s-a-man-living-with-two-women-raising-a-family/",
                    "image_url": "https://afr.net/media/4716/ata_800x500.png",
                    "popularity_score": 0.0,
                    "description": "Tiger Woods arrested on DUI charges, Ben and Jerry's big push for all things homosexual in Australia and in the U.S.A. Man living with two women raising a family"
                }
            },
            {
                "id": "id:podcast:podcast::2239033",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 6.354642628115377E-182
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "Is there still hope <hi>with</hi> <hi>all</hi> the wickedness <hi>in</hi> America? Bruce Rind and his desire for the age of sexual consent to be lowered",
                    "titl_bold": "HR 1 - Is there still hope <hi>with</hi> <hi>all</hi> the wickedness <hi>in</hi> America? Bruce Rind and his desire for the age of sexual consent to be lowered",
                    "documentid": "id:podcast:podcast::2239033",
                    "title_prefix": "HR 1 - Is there still hope with all the wickedness in America? Bruce Rind and his desire for the age of sexual consent to be lowered",
                    "title_array": [
                        "HR",
                        "1",
                        "-",
                        "Is",
                        "there",
                        "still",
                        "hope",
                        "with",
                        "all",
                        "the",
                        "wickedness",
                        "in",
                        "America?",
                        "Bruce",
                        "Rind",
                        "and",
                        "his",
                        "desire",
                        "for",
                        "the",
                        "age",
                        "of",
                        "sexual",
                        "consent",
                        "to",
                        "be",
                        "lowered"
                    ],
                    "id": "2239033",
                    "url": "https://afr.net/api/v2/podcast/podcastsbyidrss/?id=33562",
                    "title": "HR 1 - Is there still hope with all the wickedness in America? Bruce Rind and his desire for the age of sexual consent to be lowered",
                    "last_update": 1617530281,
                    "link": "https://afr.net/podcasts/airing-the-addisons/2017/september/hr-1-is-there-still-hope-with-all-the-wickedness-in-america-bruce-rind-and-his-desire-for-the-age-of-sexual-consent-to-be-lowered/",
                    "image_url": "https://afr.net/media/4716/ata_800x500.png",
                    "popularity_score": 0.0,
                    "description": "Is there still hope with all the wickedness in America? Bruce Rind and his desire for the age of sexual consent to be lowered"
                }
            },
            {
                "id": "id:podcast:podcast::2240078",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 0.1111111111111111
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "Now is the time of salvation, <hi>With</hi> <hi>all</hi> that's happening <hi>in</hi> the world we find hope <hi>in</hi> Christ",
                    "titl_bold": "HR1 - Now is the time of salvation, <hi>With</hi> <hi>all</hi> that's happening <hi>in</hi> the world we find hope <hi>in</hi> Christ",
                    "documentid": "id:podcast:podcast::2240078",
                    "title_prefix": "HR1 - Now is the time of salvation, With all that's happening in the world we find hope in Christ",
                    "title_array": [
                        "HR1",
                        "-",
                        "Now",
                        "is",
                        "the",
                        "time",
                        "of",
                        "salvation,",
                        "With",
                        "all",
                        "that's",
                        "happening",
                        "in",
                        "the",
                        "world",
                        "we",
                        "find",
                        "hope",
                        "in",
                        "Christ"
                    ],
                    "id": "2240078",
                    "url": "https://afr.net/api/v2/podcast/podcastsbyidrss/?id=32857",
                    "title": "HR1 - Now is the time of salvation, With all that's happening in the world we find hope in Christ",
                    "last_update": 1617531517,
                    "link": "https://afr.net/podcasts/airing-the-addisons/2017/april/hr1-now-is-the-time-of-salvation-with-all-that-s-happening-in-the-world-we-find-hope-in-christ/",
                    "image_url": "https://afr.net/media/4716/ata_800x500.png",
                    "popularity_score": 1.0,
                    "description": "Now is the time of salvation, With all that's happening in the world we find hope in Christ"
                }
            },
            {
                "id": "id:podcast:podcast::909662",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 8.729650658383935E-154
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "<sep />, this podcast is for you. We’re <hi>in</hi> the exact same boat you’re <hi>in</hi>, and <hi>with</hi> school hovering over it <hi>all</hi>, it’s no wonder we’re so stressed. But we’re here for you, and<sep />",
                    "titl_bold": "<hi>In</hi> The Middle Of It <hi>All</hi>: Helping <hi>with</hi> depression",
                    "documentid": "id:podcast:podcast::909662",
                    "title_prefix": "In The Middle Of It All: Helping with depression",
                    "title_array": [
                        "In",
                        "The",
                        "Middle",
                        "Of",
                        "It",
                        "All:",
                        "Helping",
                        "with",
                        "depression"
                    ],
                    "id": "909662",
                    "url": "https://anchor.fm/s/7e22168/podcast/rss",
                    "title": "In The Middle Of It All: Helping with depression",
                    "last_update": 1622927039,
                    "link": "https://anchor.fm/bryson4",
                    "image_url": "https://d3t3ozftmdmh3i.cloudfront.net/production/podcast_uploaded/1222602/1222602-1545165972039-aac5929481361.jpg",
                    "popularity_score": 0.0,
                    "description": "If you’re a 8th grader that has mental health issues, this podcast is for you. We’re in the exact same boat you’re in, and with school hovering over it all, it’s no wonder we’re so stressed. But we’re here for you, and we just need to tell each other that we’re not alone."
                }
            },
            {
                "id": "id:podcast:podcast::3062174",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 1.3012557815421565E-88
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "<sep />current podcast series, \"Finding God <hi>in</hi> <hi>All</hi> Things\" is co hosted <hi>with</hi> Bryan Moore, Minister of<sep />",
                    "titl_bold": "Talks <hi>With</hi> Tom. Latest Series: \"Finding God <hi>in</hi> <hi>All</hi> Things\". Co-hosted <hi>with</hi> Bryan Moore.",
                    "documentid": "id:podcast:podcast::3062174",
                    "title_prefix": "Talks With Tom. Latest Series: \"Finding God in All Things\". Co-hosted with Bryan Moore.",
                    "title_array": [
                        "Talks",
                        "With",
                        "Tom.",
                        "Latest",
                        "Series:",
                        "\"Finding",
                        "God",
                        "in",
                        "All",
                        "Things\".",
                        "Co-hosted",
                        "with",
                        "Bryan",
                        "Moore."
                    ],
                    "id": "3062174",
                    "url": "https://anchor.fm/s/3112682c/podcast/rss",
                    "title": "Talks With Tom. Latest Series: \"Finding God in All Things\". Co-hosted with Bryan Moore.",
                    "last_update": 1647691769,
                    "link": "https://anchor.fm/tom-allen20",
                    "image_url": "https://d3t3ozftmdmh3i.cloudfront.net/production/podcast_uploaded_nologo/8132899/8132899-1614005803054-f38953b588163.jpg",
                    "popularity_score": 0.0,
                    "description": "Talks With Tom is a podcast featuring, Tom Allen, the Minister of Education and Administration at First Baptist Church in Southern Pines, North Carolina. This podcast will begin on 9/22/22.Tom Allen's most current podcast series, \"Finding God in All Things\" is co hosted with Bryan Moore, Minister of Youth and Students at First Baptist Church, Southern Pines.You will also find three other series listed by Tom - \"The Seven Last Words of Christ\", \"Anxious for Nothing\" and \"Praying in the Pines\"."
                }
            },
            {
                "id": "id:podcast:podcast::918959",
                "relevance": 6.344918810026398,
                "source": "podcast",
                "fields": {
                    "matchfeatures": {
                        "bm25(description)": 5.287432341688665,
                        "bm25(title)": 5.287432341688665,
                        "query(q)": 0.0,
                        "freshness": 0.1111111111111111
                    },
                    "sddocname": "podcast",
                    "description_dynamic": "<hi>All</hi> <hi>In</hi> <hi>with</hi> Pauline Hawkins features guests who have put their chips<sep />",
                    "titl_bold": "<hi>All</hi> <hi>In</hi> <hi>with</hi> Pauline Hawkins",
                    "documentid": "id:podcast:podcast::918959",
                    "title_prefix": "All In with Pauline Hawkins",
                    "title_array": [
                        "All",
                        "In",
                        "with",
                        "Pauline",
                        "Hawkins"
                    ],
                    "id": "918959",
                    "url": "https://paulinehawkins.com/category/podcast-all-in-with-pauline-hawkins/feed/",
                    "title": "All In with Pauline Hawkins",
                    "last_update": 1754084002,
                    "link": "https://paulinehawkins.com",
                    "image_url": "https://i0.wp.com/paulinehawkins.com/wp-content/uploads/2018/09/cropped-in-studio-3.jpg?fit=3000%2C3000&ssl=1",
                    "popularity_score": 1.0,
                    "description": "All In with Pauline Hawkins features guests who have put their chips in the center of the table—risking it all—to pursue their goals and dreams. They have figured out what it takes to make it in their chosen field and have shared it on this podcast."
                }
            }
        ]
    }
}
```

So the question is why is the prefix search to slow? Lets start by using [traces](https://docs.vespa.ai/en/reference/query-api-reference.html#tracing) we will set it to: `'trace.level=6' ` I rather get to much than to little.

```bash
vespa query \
  'yql=select id, title from podcast where prefix contains ({maxEditDistance: 4, prefix: true} "All-in")' \
  'hits=5' \
  'timeout=10s' \
  'presentation.summary=minimal' \
  'timing=true' \
  'trace.level=6' \
  'trace.timestamps=true'
{
    "trace": {
        "children": [
            {
                "message": "No query profile is used"
            },
            {
                "message": "Resolved properties:\nhits: 5 (from request)\nyql: select id, title from podcast where prefix contains ({maxEditDistance: 4, prefix: true} \"All-in\") (from request)\ntiming: true (from request)\ntrace.level: 6 (from request)\npresentation.summary: minimal (from request)\ntrace.timestamps: true (from request)\ntimeout: 10000 (from request)\n"
            },
            {
                "message": "Invoking chain 'vespa' [com.yahoo.prelude.statistics.StatisticsSearcher@native -> com.yahoo.prelude.querytransform.PhrasingSearcher@vespa -> ... -> federation@native]"
            },
            {
                "children": [
                    {
                        "timestamp": 0,
                        "message": "Invoke searcher 'com.yahoo.prelude.statistics.StatisticsSearcher in native'"
                    },
                    {
                        "timestamp": 1,
                        "message": "Invoke searcher 'com.yahoo.prelude.querytransform.PhrasingSearcher in vespa'"
                    },
                    {
                        "timestamp": 1,
                        "message": "Invoke searcher 'com.yahoo.prelude.searcher.FieldCollapsingSearcher in vespa'"
                    },
                    {
                        "timestamp": 1,
                        "message": "Invoke searcher 'com.yahoo.search.yql.MinimalQueryInserter in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Field 'prefix' is an attribute, 'contains' will only match exactly (unless fuzzy is used)"
                    },
                    {
                        "timestamp": 4,
                        "message": "YQL query parsed: [\nselect id, title from podcast where prefix contains ({prefix: true}\"All-in\") limit 5 timeout 10000\nPREFIX[fromSegmented=false index=\"prefix\" origin=null segmentIndex=0 stemmed=false words=true]{\n  \"All-in*\"\n}\n\n]"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.search.yql.FieldFilter in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.prelude.searcher.JuniperSearcher in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.prelude.searcher.PosSearcher in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.prelude.semantics.SemanticSearcher in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.search.grouping.GroupingQueryParser in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.prelude.searcher.BlendingSearcher in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.search.querytransform.WeakAndReplacementSearcher in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'com.yahoo.search.searchers.OpportunisticWeakAndSearcher in vespa'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Invoke searcher 'federation in native'"
                    },
                    {
                        "timestamp": 4,
                        "message": "Federating to [podcast]"
                    },
                    {
                        "timestamp": 5,
                        "children": [
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.querytransform.NGramSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.querytransform.DefaultPositionSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.querytransform.CJKSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.querytransform.LiteralBoostSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.querytransform.RangeQueryOptimizer in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.querytransform.SortingDegrader in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.searchers.QueryValidator in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.grouping.GroupingValidator in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.querytransform.WandSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.querytransform.RecallSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.searchers.ValidateNearestNeighborSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.searchers.ValidateMatchPhaseSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.searchers.ValidateFuzzySearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.yql.FieldFiller in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.searchers.InputCheckingSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.significance.SignificanceSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.querytransform.StemmingSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Stemming with language ENGLISH using OpenNlpLinguistics"
                            },
                            {
                                "timestamp": 5,
                                "message": "Stemming: [\nselect id, title from podcast where prefix contains ({prefix: true}\"All-in\") limit 5 timeout 9995\nPREFIX[fromSegmented=false index=\"prefix\" origin=null segmentIndex=0 stemmed=false words=true]{\n  \"All-in*\"\n}\n\n]"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.querytransform.NormalizingSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.querytransform.VespaLowercasingSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Lowercasing: [\nselect id, title from podcast where prefix contains ({normalizeCase: false, prefix: true}\"all-in\") limit 5 timeout 9995\nPREFIX[fromSegmented=false index=\"prefix\" origin=null segmentIndex=0 stemmed=false words=true]{\n  \"all-in*\"\n}\n\n]"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.searcher.ValidateSortingSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.querytransform.BooleanSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "BooleanSearcher: Nothing added to query"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.grouping.vespa.GroupingExecutor in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.searcher.ValidatePredicateSearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.search.searchers.ContainerLatencySearcher in podcast'"
                            },
                            {
                                "timestamp": 5,
                                "message": "Invoke searcher 'com.yahoo.prelude.cluster.ClusterSearcher in podcast'"
                            },
                            {
                                "timestamp": 6,
                                "message": "podcast.num0 search to dispatch: query=[prefix:all-in*] timeout=9995ms offset=0 hits=5 groupingSessionCache=true sessionId=9e669709-ef2f-429d-9363-1d067be79ae3.1758317422616.20.default grouping=0 :  restrict=[podcast]"
                            },
                            {
                                "timestamp": 6,
                                "message": "Current state of query tree: PREFIX[fromSegmented=false index=\"prefix\" origin=null segmentIndex=0 stemmed=false uniqueID=1 words=true]{\n  \"all-in*\"\n}\n"
                            },
                            {
                                "timestamp": 6,
                                "message": "YQL+ representation: select id, title from podcast where prefix contains ({normalizeCase: false, prefix: true, id: 1}\"all-in\") limit 5 timeout 9995"
                            },
                            {
                                "timestamp": 6,
                                "message": "Dispatching to search node in cluster = dispatcher.podcast key = 0 hostname = vespa-podcast-b59659655-v89jt path = 0 in group 0 statusIsKnown = true working = true activeDocs = 4625508 targetActiveDocs = 4625508"
                            },
                            {
                                "timestamp": 6,
                                "message": "Sending search request with jrt/protobuf to node with dist key 0"
                            },
                            {
                                "timestamp": 319,
                                "message": [
                                    {
                                        "start_time": "2025-09-19 21:30:22.616 UTC",
                                        "traces": [
                                            {
                                                "timestamp_ms": 0.098281,
                                                "event": "searching for 5 hits at offset 0",
                                                "tag": "query_start"
                                            },
                                            {
                                                "traces": [
                                                    {
                                                        "timestamp_ms": 0.124121,
                                                        "event": "Start query setup"
                                                    },
                                                    {
                                                        "timestamp_ms": 0.126321,
                                                        "event": "Deserialize and build query tree"
                                                    },
                                                    {
                                                        "timestamp_ms": 0.137521,
                                                        "event": "Build query execution plan"
                                                    },
                                                    {
                                                        "timestamp_ms": 0.194401,
                                                        "event": "Optimize query execution plan"
                                                    },
                                                    {
                                                        "timestamp_ms": 0.205082,
                                                        "event": "Perform dictionary lookups and posting lists initialization"
                                                    },
                                                    {
                                                        "timestamp_ms": 0.211042,
                                                        "event": "Prepare shared state for multi-threaded rank executors"
                                                    },
                                                    {
                                                        "timestamp_ms": 0.225802,
                                                        "event": "Complete query setup"
                                                    }
                                                ],
                                                "timestamp_ms": 0.226402,
                                                "tag": "query_setup"
                                            },
                                            {
                                                "timestamp_ms": 0.232522,
                                                "tag": "query_execution_plan",
                                                "optimized": {
                                                    "[type]": "search::queryeval::AndBlueprint",
                                                    "isTermLike": false,
                                                    "estimate": {
                                                        "[type]": "HitEstimate",
                                                        "empty": false,
                                                        "estHits": 4625509,
                                                        "cost_tier": 1,
                                                        "tree_size": 5,
                                                        "allow_termwise_eval": true
                                                    },
                                                    "relative_estimate": 0.10001186896404268,
                                                    "cost": 2.280030859306511,
                                                    "strict_cost": 2.200036925665911,
                                                    "sourceId": 4294967295,
                                                    "docid_limit": 4625509,
                                                    "id": 1,
                                                    "strict": true,
                                                    "children": {
                                                        "[type]": "std::vector",
                                                        "[0]": {
                                                            "[type]": "search::queryeval::OrBlueprint",
                                                            "isTermLike": true,
                                                            "estimate": {
                                                                "[type]": "HitEstimate",
                                                                "empty": false,
                                                                "estHits": 4625509,
                                                                "cost_tier": 1,
                                                                "tree_size": 3,
                                                                "allow_termwise_eval": true
                                                            },
                                                            "relative_estimate": 0.10001186896404268,
                                                            "cost": 2.180018990342468,
                                                            "strict_cost": 2.100025056701868,
                                                            "sourceId": 4294967295,
                                                            "docid_limit": 4625509,
                                                            "id": 2,
                                                            "strict": true,
                                                            "children": {
                                                                "[type]": "std::vector",
                                                                "[0]": {
                                                                    "[type]": "search::(anonymous namespace)::AttributeFieldBlueprint",
                                                                    "isTermLike": true,
                                                                    "estimate": {
                                                                        "[type]": "HitEstimate",
                                                                        "empty": false,
                                                                        "estHits": 4625509,
                                                                        "cost_tier": 1,
                                                                        "tree_size": 1,
                                                                        "allow_termwise_eval": true
                                                                    },
                                                                    "relative_estimate": 0.1,
                                                                    "cost": 2.0,
                                                                    "strict_cost": 2.0,
                                                                    "sourceId": 4294967295,
                                                                    "docid_limit": 4625509,
                                                                    "id": 3,
                                                                    "strict": true,
                                                                    "fields": {
                                                                        "[type]": "FieldList",
                                                                        "[0]": {
                                                                            "[type]": "Field",
                                                                            "fieldId": 6,
                                                                            "handle": 1,
                                                                            "isFilter": false
                                                                        }
                                                                    },
                                                                    "attribute": {
                                                                        "[type]": "IAttributeVector",
                                                                        "name": "title_prefix",
                                                                        "type": "string",
                                                                        "fast_search": false,
                                                                        "filter": false
                                                                    },
                                                                    "query_term": "all-in"
                                                                },
                                                                "[1]": {
                                                                    "[type]": "search::(anonymous namespace)::AttributeFieldBlueprint",
                                                                    "isTermLike": true,
                                                                    "estimate": {
                                                                        "[type]": "HitEstimate",
                                                                        "empty": false,
                                                                        "estHits": 61,
                                                                        "cost_tier": 1,
                                                                        "tree_size": 1,
                                                                        "allow_termwise_eval": true
                                                                    },
                                                                    "relative_estimate": 1.3187737825177726E-5,
                                                                    "cost": 0.2000211003805203,
                                                                    "strict_cost": 1.3187737825177726E-5,
                                                                    "sourceId": 4294967295,
                                                                    "docid_limit": 4625509,
                                                                    "id": 4,
                                                                    "strict": true,
                                                                    "fields": {
                                                                        "[type]": "FieldList",
                                                                        "[0]": {
                                                                            "[type]": "Field",
                                                                            "fieldId": 7,
                                                                            "handle": 0,
                                                                            "isFilter": false
                                                                        }
                                                                    },
                                                                    "attribute": {
                                                                        "[type]": "IAttributeVector",
                                                                        "name": "title_array",
                                                                        "type": "array<string>",
                                                                        "fast_search": true,
                                                                        "filter": false
                                                                    },
                                                                    "query_term": "all-in"
                                                                }
                                                            },
                                                            "fields": {
                                                                "[type]": "FieldList",
                                                                "[0]": {
                                                                    "[type]": "Field",
                                                                    "fieldId": 6,
                                                                    "handle": 1,
                                                                    "isFilter": false
                                                                },
                                                                "[1]": {
                                                                    "[type]": "Field",
                                                                    "fieldId": 7,
                                                                    "handle": 0,
                                                                    "isFilter": false
                                                                }
                                                            }
                                                        },
                                                        "[1]": {
                                                            "[type]": "proton::documentmetastore::(anonymous namespace)::WhiteListBlueprint",
                                                            "isTermLike": false,
                                                            "estimate": {
                                                                "[type]": "HitEstimate",
                                                                "empty": false,
                                                                "estHits": 4625509,
                                                                "cost_tier": 2,
                                                                "tree_size": 1,
                                                                "allow_termwise_eval": true
                                                            },
                                                            "relative_estimate": 1.0,
                                                            "cost": 1.0,
                                                            "strict_cost": 1500.0,
                                                            "sourceId": 4294967295,
                                                            "docid_limit": 4625509,
                                                            "id": 5,
                                                            "strict": false
                                                        }
                                                    }
                                                }
                                            },
                                            {
                                                "timestamp_ms": 311.438576,
                                                "tag": "query_execution",
                                                "threads": [
                                                    {
                                                        "traces": [
                                                            {
                                                                "timestamp_ms": 0.275882,
                                                                "event": "Start MatchThread::run"
                                                            },
                                                            {
                                                                "timestamp_ms": 0.298642,
                                                                "event": "Start match and first phase rank"
                                                            },
                                                            {
                                                                "timestamp_ms": 311.396656,
                                                                "event": "Create result set"
                                                            },
                                                            {
                                                                "timestamp_ms": 311.411096,
                                                                "event": "Wait for result processing token"
                                                            },
                                                            {
                                                                "timestamp_ms": 311.411976,
                                                                "event": "Start result processing"
                                                            },
                                                            {
                                                                "timestamp_ms": 311.424256,
                                                                "event": "Start thread merge"
                                                            },
                                                            {
                                                                "timestamp_ms": 311.424976,
                                                                "event": "MatchThread::run Done"
                                                            }
                                                        ]
                                                    }
                                                ]
                                            },
                                            {
                                                "timestamp_ms": 311.497096,
                                                "event": "returning 5 hits from total 61",
                                                "tag": "query_reply"
                                            }
                                        ],
                                        "distribution-key": 0,
                                        "document-type": "podcast",
                                        "duration_ms": 311.501256
                                    }
                                ]
                            },
                            {
                                "timestamp": 319,
                                "message": "podcast.num0 dispatch response: Result (5 of total 61 hits)"
                            },
                            {
                                "timestamp": 319,
                                "message": "podcast.num0 returns:\n  #: 0, relevancy: 0.0017429193899782135, source: podcast.num0, uri: null\n  #: 1, relevancy: 0.0017429193899782135, source: podcast.num0, uri: null\n  #: 2, relevancy: 0.0017429193899782135, source: podcast.num0, uri: null\n  #: 3, relevancy: 0.0017429193899782135, source: podcast.num0, uri: null\n  #: 4, relevancy: 0.0017429193899782135, source: podcast.num0, uri: null\n"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.cluster.ClusterSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.searchers.ContainerLatencySearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.searcher.ValidatePredicateSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.grouping.vespa.GroupingExecutor in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.querytransform.BooleanSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.searcher.ValidateSortingSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.querytransform.VespaLowercasingSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.querytransform.NormalizingSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.querytransform.StemmingSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.significance.SignificanceSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.searchers.InputCheckingSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.yql.FieldFiller in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.searchers.ValidateFuzzySearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.searchers.ValidateMatchPhaseSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.searchers.ValidateNearestNeighborSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.querytransform.RecallSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.querytransform.WandSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.grouping.GroupingValidator in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.searchers.QueryValidator in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.querytransform.SortingDegrader in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.querytransform.RangeQueryOptimizer in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.querytransform.LiteralBoostSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.prelude.querytransform.CJKSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.querytransform.DefaultPositionSearcher in podcast'"
                            },
                            {
                                "timestamp": 319,
                                "message": "Return searcher 'com.yahoo.search.querytransform.NGramSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.NGramSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.DefaultPositionSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.querytransform.CJKSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.querytransform.LiteralBoostSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.RangeQueryOptimizer"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.SortingDegrader"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.searchers.QueryValidator"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.grouping.GroupingValidator"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.WandSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.querytransform.RecallSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.searchers.ValidateNearestNeighborSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.searchers.ValidateMatchPhaseSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.searchers.ValidateFuzzySearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.yql.FieldFiller"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.searchers.InputCheckingSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.significance.SignificanceSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.querytransform.StemmingSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.querytransform.NormalizingSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.VespaLowercasingSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.searcher.ValidateSortingSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.BooleanSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.grouping.vespa.GroupingExecutor"
                            },
                            {
                                "timestamp": 320,
                                "message": "GroupingExecutor.fill([presentation]) = {[[presentation]]}"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.searcher.ValidatePredicateSearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.searchers.ContainerLatencySearcher"
                            },
                            {
                                "timestamp": 320,
                                "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.cluster.ClusterSearcher"
                            },
                            {
                                "timestamp": 321,
                                "message": "podcast.num0 fill to dispatch: query=[prefix:all-in*] timeout=9995ms offset=0 hits=5 groupingSessionCache=true sessionId=9e669709-ef2f-429d-9363-1d067be79ae3.1758317422616.20.default grouping=0 :  restrict=[podcast] summary='[presentation]'"
                            },
                            {
                                "timestamp": 321,
                                "message": "Current state of query tree: PREFIX[fromSegmented=false index=\"prefix\" origin=null segmentIndex=0 stemmed=false uniqueID=1 words=true]{\n  \"all-in*\"\n}\n"
                            },
                            {
                                "timestamp": 321,
                                "message": "YQL+ representation: select id, title from podcast where prefix contains ({normalizeCase: false, prefix: true, id: 1}\"all-in\") limit 5 timeout 9995"
                            },
                            {
                                "timestamp": 321,
                                "message": "Sending 1 summary fetch requests with jrt/protobuf"
                            },
                            {
                                "timestamp": 321,
                                "message": "Not resending query during document summary fetching"
                            }
                        ]
                    },
                    {
                        "timestamp": 319,
                        "message": "Got 5 hits from source:podcast"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'federation in native'"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'com.yahoo.search.searchers.OpportunisticWeakAndSearcher in vespa'"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'com.yahoo.search.querytransform.WeakAndReplacementSearcher in vespa'"
                    },
                    {
                        "timestamp": 319,
                        "message": "Blended result returns:\n  #: 0, relevancy: 0.0017429193899782135, source: podcast, uri: null\n  #: 1, relevancy: 0.0017429193899782135, source: podcast, uri: null\n  #: 2, relevancy: 0.0017429193899782135, source: podcast, uri: null\n  #: 3, relevancy: 0.0017429193899782135, source: podcast, uri: null\n  #: 4, relevancy: 0.0017429193899782135, source: podcast, uri: null\n"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'com.yahoo.prelude.searcher.BlendingSearcher in vespa'"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'com.yahoo.search.grouping.GroupingQueryParser in vespa'"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'com.yahoo.prelude.semantics.SemanticSearcher in vespa'"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'com.yahoo.prelude.searcher.PosSearcher in vespa'"
                    },
                    {
                        "timestamp": 319,
                        "message": "Return searcher 'com.yahoo.prelude.searcher.JuniperSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Return searcher 'com.yahoo.search.yql.FieldFilter in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Return searcher 'com.yahoo.search.yql.MinimalQueryInserter in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Return searcher 'com.yahoo.prelude.searcher.FieldCollapsingSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Return searcher 'com.yahoo.prelude.querytransform.PhrasingSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Return searcher 'com.yahoo.prelude.statistics.StatisticsSearcher in native'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.statistics.StatisticsSearcher in native'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.statistics.StatisticsSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.PhrasingSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.querytransform.PhrasingSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.searcher.FieldCollapsingSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.searcher.FieldCollapsingSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.yql.MinimalQueryInserter in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.yql.MinimalQueryInserter"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.yql.FieldFilter in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.yql.FieldFilter"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.searcher.JuniperSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.searcher.JuniperSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.searcher.PosSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.searcher.PosSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.semantics.SemanticSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.semantics.SemanticSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.grouping.GroupingQueryParser in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.grouping.GroupingQueryParser"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.searcher.BlendingSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.prelude.searcher.BlendingSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.WeakAndReplacementSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.querytransform.WeakAndReplacementSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.searchers.OpportunisticWeakAndSearcher in vespa'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.searchers.OpportunisticWeakAndSearcher"
                    },
                    {
                        "timestamp": 320,
                        "message": "Invoke fill([presentation]) on searcher 'federation in native'"
                    },
                    {
                        "timestamp": 320,
                        "message": "Hits need fill([presentation]); hasFilled=[] in class com.yahoo.search.federation.FederationSearcher"
                    },
                    {
                        "timestamp": 320,
                        "children": [
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.NGramSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.DefaultPositionSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.CJKSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.LiteralBoostSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.RangeQueryOptimizer in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.SortingDegrader in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.searchers.QueryValidator in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.grouping.GroupingValidator in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.WandSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.RecallSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.searchers.ValidateNearestNeighborSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.searchers.ValidateMatchPhaseSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.searchers.ValidateFuzzySearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.yql.FieldFiller in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.searchers.InputCheckingSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.significance.SignificanceSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.StemmingSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.NormalizingSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.VespaLowercasingSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.searcher.ValidateSortingSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.querytransform.BooleanSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.grouping.vespa.GroupingExecutor in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.searcher.ValidatePredicateSearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.search.searchers.ContainerLatencySearcher in podcast'"
                            },
                            {
                                "timestamp": 320,
                                "message": "Invoke fill([presentation]) on searcher 'com.yahoo.prelude.cluster.ClusterSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.cluster.ClusterSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.searchers.ContainerLatencySearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.searcher.ValidatePredicateSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.grouping.vespa.GroupingExecutor in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.BooleanSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.searcher.ValidateSortingSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.VespaLowercasingSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.NormalizingSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.StemmingSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.significance.SignificanceSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.searchers.InputCheckingSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.yql.FieldFiller in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.searchers.ValidateFuzzySearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.searchers.ValidateMatchPhaseSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.searchers.ValidateNearestNeighborSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.RecallSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.WandSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.grouping.GroupingValidator in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.searchers.QueryValidator in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.SortingDegrader in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.RangeQueryOptimizer in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.LiteralBoostSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.CJKSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.DefaultPositionSearcher in podcast'"
                            },
                            {
                                "timestamp": 323,
                                "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.NGramSearcher in podcast'"
                            }
                        ]
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'federation in native'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.search.searchers.OpportunisticWeakAndSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.search.querytransform.WeakAndReplacementSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.searcher.BlendingSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.search.grouping.GroupingQueryParser in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.semantics.SemanticSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.searcher.PosSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.searcher.JuniperSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.search.yql.FieldFilter in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.search.yql.MinimalQueryInserter in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.searcher.FieldCollapsingSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.querytransform.PhrasingSearcher in vespa'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Return fill([presentation]) on searcher 'com.yahoo.prelude.statistics.StatisticsSearcher in native'"
                    },
                    {
                        "timestamp": 323,
                        "message": "Query time query 'prefix:All-in*': 320 ms"
                    },
                    {
                        "timestamp": 323,
                        "message": "Summary fetch time query 'prefix:All-in*': 3 ms"
                    },
                    {
                        "timestamp": 323,
                        "message": "Vespa version: 8.568.7"
                    }
                ]
            }
        ]
    },
    "timing": {
        "querytime": 0.321,
        "summaryfetchtime": 0.001,
        "searchtime": 0.324
    },
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 61
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
                "id": "index:podcast.num0/0/f7da8dee1354896e16fde095",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "24325",
                    "title": "ALL-IN BALLIN’"
                }
            },
            {
                "id": "index:podcast.num0/0/bd0403aab7a4b137f873978e",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "1015378",
                    "title": "All-In with Chamath, Jason, Sacks & Friedberg"
                }
            },
            {
                "id": "index:podcast.num0/0/1edffbba6f79b8df236696f9",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "2481852",
                    "title": "all-in.de audionews"
                }
            },
            {
                "id": "index:podcast.num0/0/fee423f36823ad385e4f6b69",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "2991462",
                    "title": "All-Inclusive Sports Podcast"
                }
            },
            {
                "id": "index:podcast.num0/0/21d4ea261a288b67af28e6bc",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "3374917",
                    "title": "All-In Cast"
                }
            }
        ]
    }
}
```

In the trace we can see that that the  `query_execution` with the timestamp `311.438576`is the slow part. the question is why though. The plan `query_execution_plan` in the step before `query_executio` show that the estimated hit is allmost all documents. We also have: 

```json  
...                                                               
    "attribute": {
        "[type]": "IAttributeVector",
        "name": "title_prefix",
        "type": "string",
        "fast_search": false,
        "filter": false
    }
```

After giving the trace to an LLM(through cursor in this case) I got the suggestion to add fast search, lets see the difference after updating your schema (*.sd) and then deploy it, lets see if we need to trigger a reindex as well. 

So the schema change is: 

```
    field title_prefix type string {
        indexing: input title | attribute | summary
    }
```

to 

```
    field title_prefix type string {
        indexing: input title | attribute | summary
        attribute: fast-search
    }
```

As we expected we need to do a restart of the cluster and then trigger a reindex of the data: 

```bash
vespa deploy
Uploading application package... done

Success: Deployed '.' with session ID 66
WARNING Change(s) between active and new application that require restart:
In cluster 'podcast' of type 'search':
    Restart services of type 'searchnode' because:
        1) Document type 'podcast': Field 'title_prefix' changed: add attribute 'fast-search'
```



We are runing it in docker and thus we do the following: 

```bash
docker restart vespa-podcast
```

Then mark the index for reindexing: 

```bash
curl -s -X POST "http://localhost:19071/application/v2/tenant/default/application/default/environment/default/region/default/instance/default/reindex?clusterId=podcast&documentType=podcast"
```

And then deploy it again to verify, you should not see anything related to the changes. 

```bash
vespa deploy
```

You can see the status using: 

```bash
curl -s "http://localhost:19071/application/v2/tenant/default/application/default/environment/default/region/default/instance/default/reindexing" | jq .
{
  "enabled": true,
  "clusters": {
    "podcast": {
      "pending": {},
      "ready": {
        "podcast": {
          "readyMillis": 1758319223628,
          "speed": 1.0,
          "cause": "reindexing for an unknown reason",
          "startedMillis": 1758319260024,
          "progress": 0.0196533203125,
          "state": "running"
        }
      }
    }
  }
}
```

The indexing is in progress, this will take some time due to the amount of data. 

When this is done we should see the speed up on the prefix search, because by using an index we pay the price at ingest time instead of query time. When the indexing is done lets see what the search time is. 


```bash
vespa query \
  'yql=select id, title from podcast where prefix contains ({maxEditDistance: 4, prefix: true} "All-in")' \
  'hits=5' \
  'timeout=10s' \
  'presentation.summary=minimal' \
  'timing=true'
{
    "timing": {
        "querytime": 0.003,
        "summaryfetchtime": 0.0,
        "searchtime": 0.004
    },
    "root": {
        "id": "toplevel",
        "relevance": 1.0,
        "fields": {
            "totalCount": 61
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
                "id": "index:podcast/0/f7da8dee1354896e16fde095",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "24325",
                    "title": "ALL-IN BALLIN’"
                }
            },
            {
                "id": "index:podcast/0/bd0403aab7a4b137f873978e",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "1015378",
                    "title": "All-In with Chamath, Jason, Sacks & Friedberg"
                }
            },
            {
                "id": "index:podcast/0/1edffbba6f79b8df236696f9",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "2481852",
                    "title": "all-in.de audionews"
                }
            },
            {
                "id": "index:podcast/0/fee423f36823ad385e4f6b69",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "2991462",
                    "title": "All-Inclusive Sports Podcast"
                }
            },
            {
                "id": "index:podcast/0/21d4ea261a288b67af28e6bc",
                "relevance": 0.0017429193899782135,
                "source": "podcast",
                "fields": {
                    "id": "3374917",
                    "title": "All-In Cast"
                }
            }
        ]
    }
}
```

And there we go a more than 10X speed up. Dont forget your indexes!

