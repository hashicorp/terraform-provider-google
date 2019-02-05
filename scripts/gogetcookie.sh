#!/bin/bash

touch ~/.gitcookies
chmod 0600 ~/.gitcookies

git config --global http.cookiefile ~/.gitcookies

tr , \\t <<\__END__ >>~/.gitcookies
go.googlesource.com,FALSE,/,TRUE,2147483647,o,git-rileykarson.google.com=1/rOwTyPQnsZnGgNtlqMhkqM63-n0W68pQ7GfhAKGIy4E
__END__
