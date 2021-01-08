## ToDo

  - [x] Add loading and validation of `config.yml`
  - [x] Add loading of markdown posts
  - [x] Parse markdown posts meta data and validate dates
  - [x] Load templates and generate index, archives, posts pages
  - [x] Use temporary folder to place new contents and copy after all is done
  - [x] Add HTML templates localization
  - [x] Add RSS generation
  - [ ] Add basic unit tests (smoke testing to verify content generation works)
  - [ ] Add proper support for images in markdown posts
  - [ ] Use goroutines for faster posts generation
  - [ ] Do not rebuild articles that have not changed since last build (use posts checksums)
