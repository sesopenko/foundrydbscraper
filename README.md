# sesopenko/foundrydbscraper

Scrapes the db file from a foundry module and produces a static website.  Created specifically to scrape the data
from the [Abomination Vaults foundry module](https://foundryvtt.com/packages/pf2e-abomination-vaults) so that I can use
it for machine learning for my own campaign. Mostly just descriptions have been mapped to html files the first pass
to see how the resulting MLL [embedding](https://python.langchain.com/docs/modules/data_connection/text_embedding/)
behaves.

Sharing the module files which aren't licensed under the Open Game License breaks copyright. Please don't contact me
looking for Abomination Vaults content.

## Static Site Structure

```
generated/
    journals.html (index of journals)
    journals/
        <journal_id>.html (one for each journal page)
    journal_pages/
        <journal_page_id>.html (one for each journal page)
```

## Capabilities

* Parses journal and journal page links and links renders links to them.
* Static site hosting on port `:8080` with optional `-s` flag.
* Theoretically could work for other campaign modules if they're structured the same.

## TODO:

* Pathfinder core:
    * Monsters
    * Items
    * Spells

## Long Term TODO

* Actual stat blocks (I forget who people are and where they are, more than stats, so doing stats later)
* Copying images and rendering links to them.

## Requirements

* [go 1.19 or later](https://go.dev/doc/install)
* knowledge of how to build binaries with go ([ChatGPT](https://chat.openai.com) can help you out easily)

## Building

```go build -o foundrydbscraper main.go```

## Running (Linux)

Inside your Foundry installation, you'll find `data/Data/modules/pf2e-abomination-vaults/packs/av.db`.

Get the full file path to this file, and set it as the environment variable `DB_PATH`

```bash
DB_PATH=/home/sesopenko/foundrydata/data/Data/modules/pf2e-abomination-vaults/packs/av.db foundrydbscraper
```

The optional `-s` parameter will make it serve the files on port `:8080`

```bash
DB_PATH=/home/sesopenko/foundrydata/data/Data/modules/pf2e-abomination-vaults/packs/av.db foundrydbscraper -s
```

*console output:*
```
2024/01/12 20:10:37 Serving at http://127.0.0.1:8080
```

## Windows Support

It will likely work on windows, with the right environment variables.

## Apache 2.0 License

This software is licensed under [Apache 2.0](LICENSE.txt)

## Copyright

This software is copyright (c) Sean Esopenko 2024, all rights reserved.