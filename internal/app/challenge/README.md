# Project Structure

This project layout follow standards from https://github.com/golang-standards/project-layout so make sure to check it out before adding new files.

# Application Description

App configures the web server to serve GET request /suggestions%q=partialOrCompleteCityName&latitude=12.34&longitude=56.78. The function Initialize must be called first to configure the web server and the cityRepository then the function Run can be called to launch the web server.

A GET request /suggestions is handled as follow:
- The request is first parsed into a cityQuery which contains the requested city name and optionally the latitude and/or the longitude.
- Then the function FindRankedSuggestionsFor is called on the cityRepository.
- FindRankedSuggestionsFor will search through city records from ./data/cities_canada-usa.tsv to find city names matching the cityQuery.name. It is the function matchQueryName role to determine if the record maches the query.  
- If a match is found a score is then calculated to determine how good this suggestion is compared to others.
- The score is computed by computeScoreFor which return the product of:
  - the ratio of matching character,
  - the weight of the distance between the query latitude and the record latitude,
  - and the weight of the distance between the query longitude and the record longitude,
- Those parameter values are between 0 and 1. When the latitude or the longitude is not present their values are 1 so that they do not affect the score product.
- When all the meaningfull suggestions have been identified they are sorted by their score in descending order so that the best suggestions come first. Those suggestions are returned by FindRankedSuggestionsFor.
- Finally, the returned suggestions, defined in suggestions.go, are converted to JSON format and sent as the reply to the GET request /suggestions.


