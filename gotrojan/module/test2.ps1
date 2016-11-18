[Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}
$webclient = new-object System.Net.WebClient
#$credCache = new-object System.Net.CredentialCache
#$creds = new-object System.Net.NetworkCredential("guest","guest")
#$credCache.Add("https://xxx/", "Basic", $creds)
$webclient.Credentials = new-object System.Net.NetworkCredential("guest","guest")
try{
$webclient.DownloadFile("https://xxx", "C:\img.png")
}
catch [Exception]
{
Write-Output $_.Exception.Message
}
"hello!!"
