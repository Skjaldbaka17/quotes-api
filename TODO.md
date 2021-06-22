# TODOS

- [ ] automate setting up the EC2 and fetching the code and running server
- [ ] Google Analytics for the site (set it up on google)
- [ ] https certificate
- [ ] Logo

- [ ] Front end for API (Create it in its own repo!) (Find template?)
- [ ] LandingPage (with minor info i.e. used by www.whothefucksaidthat.com + some quotes + tiers/pricing info)
- [ ] SignUp / Login (Add googleLogin?)
- [ ] Move Users Backend to Front End Repo?
- [ ] HomePage for users (History of requests + Tier + upgrade / downgrade tier)
- [ ] Pay with Crypto

- [ ] Separate WhoTheFuckSaidThat.com from the API (i.e. have as its own APP that queries the API!)
- [ ] find / buy an url for the API (quotel.com?)

- [ ] Test/Use Amazon's API Gateway (Cheaper?) //https://www.quora.com/What-is-the-best-and-cheapest-way-of-hosting-REST-API-as-a-startup-I-am-using-AWS-EC2-but-I-am-not-sure-whether-that-is-the-best-option-or-not-for-the-startup-who-has-limited-budget
- [ ] Test/Use ElasticBeans (Cheaper?)
- [ ] Test/Use Container Service (Cheaper?)
- [ ] Test/Use AWS lambda (Cheaper?)
- [ ] User Heroku (if Cheaper)

- [ ] Save History of errors (i.e. error logs)?

 -------------- Further Stuff  --------------

- [ ] Draw up DB-Graph (i.e. how tables are connected to view etc)

- [ ] Optimize ApiKey Validation queries (index created_at dates-column?)

- [ ] Change to Use Gorm to the fullest, oooooooorrr just change returned json to : {"name":"authorName", "id":authorId, "hasIcelandicQuotes":true/false, "nrOfEnglishQuoes":int, "nrOfIcelandicQuotes":int, "quotes":[{"quote": "theQuote", "id":quoteId, "isIcelandic": true/false}]}

- [ ] Insert Quote for created author or for a 'real' author (private and public)
- [ ] update inserted quote (priv and pub)
- [ ] Create new Author (private and public)
- [ ] Update created author (priv and pub)
- [ ] Create new Topic (private and public)
- [ ] update created topic (priv and pub)

- [ ] is random truly random (i.e. does the "random" funcitonality truly return randomly or is it biased towards quotes in the "front" of the DB (i.e. in the front where postgres stores them))
- [ ] Make Authors Search more efficient (create a similarity-based index ?)
- [ ] Sort return list alphabetically Icelandic support

- [ ] Look into and maybe Change rest into GraphQL (neeeeeee, frekar fyrir næsta project)
- [ ] Look into payment for some privileges
- [ ] New crawler for new quotes / authors

 -------------- Done  --------------

- [x] Frontend Look fixes according to Roberto
- [x] More info about author (Wikipedia link + birth-death i.e. for example 1901-2000)
- [x] Buy and setup Domain name
- [x] Make Random query faster
- [x] Setup frontend on / route for displaying a random quote.
- [x] Setup AWS (EC2) server
- [x] Setup RDS-Postgres db on AWS and setup the quotes
- [x] Add .env variables (i.e. for names of tables etc.)
- [x] Coordinate naming convention (apiKey vs api_key vs apikey etc)
- [x] only return keys, in the response-json, that are relevant to the request
- [x] CleanUp DB after tests
- [x] add api key to swagger
- [x] Add authentication for access to the api + Creating apiKeys + Documenting usage + admin access vs normal access
- [x] Add password protection / protected routes capability (at least for SetQuoteOfTheyDay route )
- [x] Save History of requests
- [x] Add Users (GOD vs ...)
- [x] Error handling
- [x] Add error response to Swagger
- [x] Go over Swagger + Clean it up and make pretty
- [x] Clean up Documentation look (Swagger)
- [x] Review /topics for Swagger 
- [x] Review /search for Swagger 
- [x] Review /quotes for Swagger 
- [x] Review /meta for Swagger 
- [x] Review /authors for Swagger 
- [x] Clean up get/set QOD/AOD
- [x] Pagination Everywhere where needed
- [x] Clean up routes files
- [x] add search for topics and searching in topics (maybe just have a single search endpoint with parameters? i.e. want to search for authors / quotes / language inside a specific topic?)
- [x] getQuotes route (combine with getQuotesById and add to it to get quotes from a specific author + add pagination)
- [x] Add Counting each time a quote is accessed / sent from Api (also for topics) - i.e. stats
- [x] Add tests for GetAOD and AODHistory and SetAOD
- [x] Add tests for GetQOD and QODHistory and SetQOD
- [x] Add get and set Author of The Day (plus points for available to set authors for multiple days in one request + plus points for AOD history)
- [x] Add get and set Quote of The Day (plus points for available to set quotes for multiple days in one request + plus points for QOD history)
- [x] quoteoftheday << qod
- [x] Get list of authors (with parameters for pagination and alphabet and languages)
- [x] Get random Author
- [x] Validate RequestBody + Tests
- [x] Add get random Quote
- [x] Use TopicsView instead of searchview (and change name to something more general)
- [x] Add Icelandic / English Support
- [x] Add categories
- [x] English and Icelandic Authors with same name have same author id
- [x] Add Search-"scroll", User is searching and is scrolling through her search and wants next batch of results matching her search i.e. PAGINATION
- [x] setup testing (unit)
- [x] Implement GetQuotesById (multiple quotes route)
- [x] Clean-up test files (Move some lines into their own functions etc.)
- [x] Setup Swagger for api-docs 
      * https://github.com/go-swagger/go-swagger
      * https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
      * https://www.youtube.com/watch?v=07XhTqE-j8k&t=374s
      * https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/handlers/docs.go


Author api: https://quotes.rest




### Helpful resources for full-text search in postgres

* https://www.opsdash.com/blog/postgres-full-text-search-golang.html 
* https://medium.com/@bencagri/implementing-multi-table-full-text-search-with-gorm-632518257d15
* https://www.freecodecamp.org/news/fuzzy-string-matching-with-postgresql/
* https://www.compose.com/articles/mastering-postgresql-tools-full-text-search-and-phrase-search/ 
