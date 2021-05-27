# TODOS

- [ ] Add icelandic / English Support
- [ ] Validate RequestBody + Tests
- [ ] Error handling + Tests
- [ ] Make error tests (i.e. made-to-fail-tests)
- [ ] Make function for building the SQL query (general for all?)
- [ ] Add categories-search (Motivational etc)
- [ ] Setup AWS server
- [ ] Make Authors Search more efficient (create a similarity-based index ?)
- [ ] Look into and maybe Change into GraphQL
- [x] Add Search-"scroll", User is searching and is scrolling through her search and wants next batch of results matching her search i.e. PAGINATION
- [x] setup testing (unit)
- [x] Implement GetQuotesById (multiple quotes route)
- [x] Clean-up test files (Move some lines into their own functions etc.)
- [x] Setup Swagger for api-docs 
      * https://github.com/go-swagger/go-swagger
      * https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
      * https://www.youtube.com/watch?v=07XhTqE-j8k&t=374s




### Helpful resources for full-text search in postgres

* https://www.opsdash.com/blog/postgres-full-text-search-golang.html 
* https://medium.com/@bencagri/implementing-multi-table-full-text-search-with-gorm-632518257d15
* https://www.freecodecamp.org/news/fuzzy-string-matching-with-postgresql/
* https://www.compose.com/articles/mastering-postgresql-tools-full-text-search-and-phrase-search/ 
