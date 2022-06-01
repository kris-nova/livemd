/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package twitch

//Request URL: https://gql.twitch.tv/gql
//Request Method: POST
//Status Code: 200 OK
//Remote Address: 199.232.78.167:443
//Referrer Policy: strict-origin-when-cross-origin

//POST /gql HTTP/1.1
//Accept: */*
//Accept-Encoding: gzip, deflate, br
//Accept-Language: en-US
//Authorization: OAuth xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//Client-Id: kimne78kx3ncx6brgo4mv6wki5h1ko
//Client-Session-Id: e4998efa01431223
//Client-Version: 9318b809-f01c-4431-a8d1-322b76e2f93a
//Connection: keep-alive
//Content-Length: 611
//Content-Type: text/plain;charset=UTF-8
//Host: gql.twitch.tv
//Origin: https://dashboard.twitch.tv
//Referer: https://dashboard.twitch.tv/
//Sec-Fetch-Dest: empty
//Sec-Fetch-Mode: cors
//Sec-Fetch-Site: same-site
//Sec-GPC: 1
//User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36
//X-Device-Id: Bdu3nkc4TX3Mdu9zk7zIREntCg4WRMuC

//[{operationName: "VideoManager_Utils_Tracking_Video", variables: {videoID: "1314727530"},…},…]
//0: {operationName: "VideoManager_Utils_Tracking_Video", variables: {videoID: "1314727530"},…}
//extensions: {,…}
//operationName: "VideoManager_Utils_Tracking_Video"
//variables: {videoID: "1314727530"}
//1: {operationName: "YoutubeExportModal_ExportVideoToYoutube",…}
//extensions: {,…}
//operationName: "YoutubeExportModal_ExportVideoToYoutube"
//variables: {input: {videoID: "1314727530",…}}
