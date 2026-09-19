# Website

The website is basically a web server that query the companies database and returns them to users.

There is no frontend framework used for the website in order to keep things simple (and fast). Instead, the website consist of html pages without any data and JS scripts to request and insert the data afterwards.

## Structure

`public` is the directory containing the website assets, like  `html`, `css` and `js` files.

`templates` is the directory for the `html` templates. The templates are compiled by the `generator.go` and then put in `public` directory.

## Generator

In order to keep things not repeating too much, the `Generator` was created for reusing `html` parts.

All the pages are compiled beforehand when running `go run -tags "fts5" ./cmd/website` (if `WebsiteGeneratorConfig.Enabled` is set to `true`).

All the `html` files from the root of the `templates` folder are considered website pages, and all the html files from `templates/partials` are considered reusable parts.

Golang [html/template](https://pkg.go.dev/html/template) is used for rendering the pages.