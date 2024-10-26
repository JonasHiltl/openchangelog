# Contributing to Openchangelog

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:  

Below are a set of guidlines for contributing to Openchangelog. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Table of Contents
1. [Tech Stack](#tech-stack)
2. [Project structure](#project-structure)
3. [Environment Setup](#environment-setup)
4. [Starting Openchangelog as a contributor](#starting-openchangelog-as-a-contributor)
5. [Creating a PR](#creating-a-pr)

## Tech Stack
- [Go](https://go.dev) The main programming language for writing the Openchangelog server.
- [Templ](https://templ.guide) For building type-safe HTML components.
- [Tailwind](https://tailwindcss.com) For easy styling of the HTML components.
- [sqlc](https://github.com/sqlc-dev/sqlc) For generating type-safe code from sql queries.

## Project structure
Openchangelog is divided into **three** separate packages.  
- The Openchangelog server is located in the repo root, it's entry point is the `cmd/server.go` file. It's the HTTP server that loads the changelog from a config or `sqlite` db and also loads it's articles through a `source` (GitHub or local). Then it parses the markdown files and responds to the user with the rendered HTML changelog.
- The **apitypes** package holds all models that are returned from the api. These are shared between the `go` api client and the Openchangelog server.
- The **api** package is the `go` api client which can be used to interact with the Openchangelog API (mostly needed in multi-tenancy setup).

## Environment Setup
Install [Go](https://go.dev/dl/), [Templ](https://templ.guide/quick-start/installation) and optionally [Air](https://github.com/air-verse/air) for live reloading.  

## Starting Openchangelog as a contributor
Create a `openchangelog.yml` file in the repo root. Have a look at the `openchangelog.example.yml` file for inspiration or just copy it's content fully for a working config.  
Run `templ generate --watch` in the repo root to have `templ` automatically generate go code from the `*.templ` files.  

Inside `internal/handler/web` run `npm run watch` to generate the `base.css` file with tailwind whenever anything changes.  

Now run `air` or `go run cmd/server.go` in the repo root to start Openchangelog with live reloading. Since `base.css` is embedded on server start, `air` sometimes doesn't update the `css` file after it changes. Rerunning `air` fixes this issue.  

Since the changelog page is cached for 5 minutes, you might need to disable the cache in the dev tools to see latest changelog updates.

## Creating a PR
If you've made changes to any `*.templ` files, ensure you run `templ generate` afterward.  
Additionally, after using watch mode, manually run `templ generate` again. Watch mode updates every `*_templ.go` file, even if no actual changes were made. Without this step, many lines may appear modified, even though no `*.templ` files were changed.  

If you changed any tailwind classes, make sure you ran `npm run watch` to generate the new `base.css` file with the tailwind styling.  

After following the above steps, you can create a PR to Openchangelog and a maintainer will review your PR swiftly.