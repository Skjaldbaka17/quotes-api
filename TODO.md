# TODOS

* Make the search better (Indexes / special tsVector columns / rankings etc) and more understandable. Note: you also need to improve the initializations in /setup-quotes-db to comply with the search/initialize indexes and views etc.
    * https://www.opsdash.com/blog/postgres-full-text-search-golang.html 
    * https://medium.com/@bencagri/implementing-multi-table-full-text-search-with-gorm-632518257d15
    * https://www.freecodecamp.org/news/fuzzy-string-matching-with-postgresql/
    * https://www.compose.com/articles/mastering-postgresql-tools-full-text-search-and-phrase-search/ 
* Implement GetQuotesById (multiple quotes route) and GetAuthorsById (Multiple authors route)
* Add categories-search (Motivational etc)
* Add Search-"scroll", User is searching and is scrolling through her search and wants next batch of results matching her search.
* Setup Testing (Unit, remember test-driven dev!)
* Setup Swagger for api-docs
* Setup AWS server
