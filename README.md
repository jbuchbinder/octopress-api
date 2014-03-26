# OCTOPRESS-API

[![Status](https://secure.travis-ci.org/jbuchbinder/octopress-api.png)](http://travis-ci.org/jbuchbinder/octopress-api)

Golang API for Octopress. This allows you to control your Octopress instance(s)
with a REST-ful API.

## Parameters

```
Usage of ./octopress-api:
  -bind=":8888": Port/IP for binding interface
  -git="git": Executable for git command
  -password="password": Password for BASIC auth
  -rake="rake": Executable for rake command
  -username="admin": Username for BASIC auth
```

Paths to Octopress sites should be given as additional parameters. If no sites are
specified, the service will not run.

## API

*Please note that this is a work in progress, until a stable version has been reached.
If you're planning on building a client based on this, please contact 
[@jbuchbinder](https://twitter.com/jbuchbinder) before you start. ;)*

### /api/version (GET)

Returns information about the version of both the software and the current API version.

### /api/1.0/sites (GET)

Returns a list of all available Octopress sites.

### /api/1.0/site/commit/SITE (GET)

Issues a git "commit" request for the specified site. SITE is the "name" parameter
of a site, which is also the map key, returned by the **/api/VERSION/sites** call.

### /api/1.0/site/deploy/SITE (GET)

Issues a generate/deploy request for the specified site. SITE is the "name" parameter
of a site, which is also the map key, returned by the **/api/VERSION/sites** call.

### /api/1.0/post/categories/SITE (GET)

Lists all post categories for the specified Octopress site.

### /api/1.0/post/list/SITE (GET)

Lists all posts and meta information for the specified Octopress site.

### /api/1.0/post/new/SITE/POSTTITLE (GET)

Issues a new_post request, and returns both the filename and post file text.

### /api/1.0/post/update/SITE/SLUG (POST)

Updates a post, based on the post slug, with the post body data.

## Building

```
go get -d
go build
```

## TODO

See [TODO](https://github.com/jbuchbinder/octopress-api/blob/master/TODO.md).

## LICENSE

BSD

