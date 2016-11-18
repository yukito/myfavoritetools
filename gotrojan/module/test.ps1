$url="https://www.yk-tb.com/contents/defcon2015/babysfirst/1/README"
$webClient = new-object System.Net.WebClient
$webClient.DownloadFile($url, "C:\README.txt")
