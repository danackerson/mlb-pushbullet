[![Circle CI](https://circleci.com/gh/danackerson/mlb-pushbullet.svg?style=shield&circle-token=3ad6694a5592b15aef77eeb7051a7b6c61d1c56f)](https://circleci.com/gh/danackerson/mlb-pushbullet)

Pushbullet: 
curl --header 'Access-Token: o.h8UMop2p50YMlbvsYpTXs6dQeqsAPTwc' https://api.pushbullet.com/v2/users/me

0. Get "upload_url":
curl --header 'Access-Token: o.h8UMop2p50YMlbvsYpTXs6dQeqsAPTwc' \
     --header 'Content-Type: application/json' \
     --data-binary '{"file_name":"asset_450.mp4","file_type":"video/mp4"}' \
     --request POST \
     https://api.pushbullet.com/v2/upload-request
{"data":{"acl":"public-read","awsaccesskeyid":"AKIAJJIUQPUDGPM4GD3W","content-type":"video/mp4","key":"vv12H5VKRKJ8F2B3rQb6gQHTPddDWdhu/asset_450.mp4","policy":"eyKjb25kaXRpb25zIjTE6MzcuMjM0MTMwWiJ9","signature":"UX5s1uIy1ov6+xlj58JY7rGFKcs="},"file_name":"asset_450.mp4","file_type":"video/mp4","file_url":"https://dl2.pushbulletusercontent.com/vv12H5VKRKJ8F2B3rQb6gQHTPddDWdhu/asset_450.mp4","upload_url":"https://upload.pushbullet.com/upload-legacy/Bf6TlYt4i3YIuCNR9Doxm9PjjG6ADhlS"}

1. POST file to "upload_url" as multipart/form-data!
curl -i -X POST https://upload.pushbullet.com/upload-legacy/eG4Pr08WVP2iKF4AqAHxaQ7sjGKUiwQQ \
  -F file=@asset_450K.mp4
  HTTP/1.1 413 Request Entity Too Large

  =========

log.txt (30KB)
0. Prepare Upload
curl --header 'Access-Token: o.h8UMop2p50YMlbvsYpTXs6dQeqsAPTwc' \
   --header 'Content-Type: application/json' \
   --data-binary '{"file_name":"log.txt","file_type":"text/plain"}' \
   --request POST \
   https://api.pushbullet.com/v2/upload-request
{"data":{"acl":"public-read","awsaccesskeyid":"AKIAJJIUQPUDGPM4GD3W","content-type":"text/plain; charset=utf-8","key":"a586XA5d3qZtmezkJBR8ANmBwiEAF5Hy/log.txt","policy":"eyKjb25kaXRpb25zIjTE6MzcuMjM0MTMwWiJ9","signature":"UX5s1uIy1ov6+xlj58JY7rGFKcs="},"file_name":"log.txt","file_type":"text/plain; charset=utf-8","file_url":"https://dl2.pushbulletusercontent.com/a586XA5d3qZtmezkJBR8ANmBwiEAF5Hy/log.txt","upload_url":"https://upload.pushbullet.com/upload-legacy/rdav5s70y1dzijZlXVkG7vov4XJrvE1E"}

1. Upload
curl -i -X POST https://upload.pushbullet.com/upload-legacy/rdav5s70y1dzijZlXVkG7vov4XJrvE1E -F file=@log.txt
HTTP/1.1 100 Continue
HTTP/1.1 204 No Content
X-Ratelimit-Limit: 16384
X-Ratelimit-Remaining: 16384
X-Ratelimit-Reset: 1462093176
Date: Sun, 01 May 2016 08:22:40 GMT
Content-Type: text/html
Server: Google Frontend
Content-Length: 0
Alt-Svc: quic=":443"; ma=2592000; v="33,32,31,30,29,28,27,26,25"

2. Push Notification
curl --header 'Access-Token: o.dNFnokvY4xyxTQIQaIugzFEaA7nRwkGm' \
     --header 'Content-Type: application/json' \
     --data-binary '{"file_name":"log.txt", "file_type":"text/plain", "file_url":"https://dl2.pushbulletusercontent.com/a586XA5d3qZtmezkJBR8ANmBwiEAF5Hy/log.txt","type":"file"}' \
     --request POST \
     https://api.pushbullet.com/v2/pushes
{"active":true,"iden":"ujvjwFibNtYsjAl8kzYqE8","created":1.462091264486548e+09,"modified":1.4620912645366223e+09,"type":"file","dismissed":false,"direction":"self","sender_iden":"ujvjwFibNtY","sender_email":"dan.ackerson@stylight.com","sender_email_normalized":"dan.ackerson@stylight.com","sender_name":"Dan Ackerson","receiver_iden":"ujvjwFibNtY","receiver_email":"dan.ackerson@stylight.com","receiver_email_normalized":"dan.ackerson@stylight.com","file_name":"log.txt","file_type":"text/plain","file_url":"https://dl2.pushbulletusercontent.com/a586XA5d3qZtmezkJBR8ANmBwiEAF5Hy/log.txt"}

3. Voila! On my Android...

========
MLB Howto
0. http://gd2.mlb.com/components/game/mlb/year_2016/month_04/day_20/grid_ce.xml
1. look for "condensed_game"
2. get $id: 605442983
3. http://m.mlb.com/gen/multimedia/detail/9/8/3/605442983.xml (where 9/8/3 are last 3 digits of $id)
and:
<url playback_scenario="FLASH_450K_400X224" ... =>
http://mediadownloads.mlb.com/mlbam/mp4/2016/04/21/605442983/1461270037182/asset_450K.mp4
<url playback_scenario="FLASH_1200K_640X360" ... =>
http://mediadownloads.mlb.com/mlbam/mp4/2016/04/21/605442983/1461270037182/asset_1200K.mp4
<url playback_scenario="FLASH_1800K_960X540" ... =>
http://mediadownloads.mlb.com/mlbam/mp4/2016/04/21/605442983/1461270037182/asset_1800K.mp4
<url playback_scenario="FLASH_2500K_1280X720" ... =>
http://mediadownloads.mlb.com/mlbam/mp4/2016/04/21/605442983/1461270037182/asset_2500K.mp4
