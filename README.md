# pingu 

###### (ping url) checks if a url is active.

### Examples

Validate that a url returns with a 200 response status code:

    pingu check https://some.url.com

Validate that a url returns a response status code other than 200:

    pingu check --expect-status=202 https://some.url.com/tryit

Validate that a url returns a 200 response status code and the content 
contains a specific text:

    pingu check --expect-content="active" https://some.url.com/status




