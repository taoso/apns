# apns
go apns client using the http2 api

May be the first APNS's go client using the [HTTP/2 API](https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/APNsProviderAPI.html)

# Usage
First, put all the tokens into a plain text file, named tokens

```
token1
token2
...
```

Second,
```
cat tokens|./apns -m "alert to push" -c "path to pem"
```

The output is

```
-token1+|{"reason":"BadDeviceToken"}
+token2|
...
```

the `+` means success, while the `-` failure, and the content after the `|` is the apns error response.

Enjoy youself :)
