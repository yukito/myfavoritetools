#!/usr/bin/env python
import urllib

PROXIES = {
    'http' : 'http://localhost:6666',
}

opener = urllib.FancyURLopener(PROXIES)
with open("testurl.txt") as fp:
   for urlinfo in fp.readlines():
      method, url = urlinfo.split(" ")
      if method == "GET":
         data = opener.open(url).read()
