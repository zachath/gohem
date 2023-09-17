# gohem
A tool for scraping properties listed on hemnet.se

## Commands
### bostad
Scrape a property to JSON format.
```
gohem bostad [url]
```

Flags:

-f : Write to file.
### bostader
Scrape properties including in a search to JSON format. Essentially a repeated `bostad` command.
```
gohem bostader [url]
```

Flags:

-f : Write to file.
## License

[MIT](https://choosealicense.com/licenses/mit/)