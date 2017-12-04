#!/usr/bin/env bash

rm -rf *.csv

grep -r "\.html" ./html/ |
sed -e 's/    <li><a href="//g' -e 's/"><\/a><\/li>//g' |
cut -d: -f2,3,4 |
sort |
uniq -c |
sort -n |
cut -d' ' -f8,7 |
tr ' ' ',' > unique_urls.csv


grep -r "\.html" ./html/ |
sed -e 's/    <li><a href="//g' -e 's/"><\/a><\/li>//g' |
sed -e 's/.\/html/http:\/\/sut:8080/g' -e 's/html:http/html,http/g' |
sort |
uniq -c |
cut -d' ' -f7,8 |
tr ' ' ',' > unique_links.csv

