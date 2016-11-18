$url = "https://www.google.co.jp/images/nav_logo242_hr.png" 
$file = "C:\web.png" 
$client = New-Object System.Net.WebClient 
$client.DownloadFile($url, $file)
