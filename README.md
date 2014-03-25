# OCTOPRESS-API

[![Status](https://secure.travis-ci.org/jbuchbinder/octopress-api.png)](http://travis-ci.org/jbuchbinder/octopress-api)

Golang API for Octopress. This allows you to control your Octopress instance(s)
with a REST-ful API.

## API

*Please note that this is a work in progress, until a stable version has been reached.
If you're planning on building a client based on this, please contact 
[@jbuchbinder](https://twitter.com/jbuchbinder) before you start. ;)*

### /api/version

Returns information about the version of both the software and the current API version.

### /api/1.0/sites

Returns a list of all available Octopress sites.

### /api/1.0/deploy/SITE

Issues a generate/deploy request for the specified site. SITE is the "name" parameter
of a site, which is also the map key, returned by the **/api/VERSION/sites** call.

### /api/1.0/post/new/SITE/POSTTITLE

Issues a new_post request, and returns both the filename and post file text.

## TODO

See [TODO](https://github.com/jbuchbinder/octopress-api/blob/master/TODO.md).

## LICENSE

BSD

