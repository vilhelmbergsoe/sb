---
title: Creating my website
date: 28/08/2022
archive: false
---

*The blog post management explanation told here is outdated. Each blog post is now parsed at start up from a file. You can check out the code [here](https://github.com/vilhelmbergsoe/sb)*
\
\
### The Beginning

I've had a domain for a while now, but haven't gotten around to building a portfolio site until now.
\
I wanted to do something a little more unique than just finding a nice [Hugo](https://gohugo.io/) theme, generate a static HTML page and call it a day.
\
Therefore, I decided to create [what you're looking at now](https://github.com/vilhelmbergsoe/sb/). A nice little personal website with blog functionality and an admin panel for creating, deleting and updating blog posts.
\
\
### The logistics

I wanted the template to be very minimal and found [this](https://minwiz.com/) nice template for a minimal responsive website.
\
I added minimal changes to it including: go templating, simple JavaScript functions for crud functionality and changed a few colors.
\
The entire website is hosted through go's [net/http](https://pkg.go.dev/net/http/) with the [gorilla HTTP router](https://github.com/gorilla/mux/).
\
The application's architecture is inspired by [Matt Ryer's blog post](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html).
\
I added markdown functionality and HTML sanitizing with [blackfriday](https://github.com/russross/blackfriday/) and [bluemonday](https://github.com/microcosm-cc/bluemonday) respectively.
\
\
### Usage

The actual template is very easy to use. Example instructions are on the GitHub README [here](https://github.com/vilhelmbergsoe/sb/).
\
There are only two endpoints: `/` and `/admin`. The admin endpoint is protected by basic authentication using nothing but go's standard library, gorilla's HTTP router and SQLite!
\
When you start out, you need to add a user to the SQLite database for administration purposes, and you can do that through a simple shell script included in the repository under `/tools/createuser`.
\
After that, you can start customizing the HTML pages and creating blog posts through the admin panel.
\
I hope to showcase some of my future projects on this website and hope this was an interesting read 😀
\
You can check out the repository [here](https://github.com/vilhelmbergsoe/sb/).
