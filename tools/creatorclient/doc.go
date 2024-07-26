// Package creatorclient is a small command line tool to test the creator service. Pipe a JSON object to it with the
// following structure:
//
//	{
//	  "Creators": [
//	    {
//	      "id": "usr_2O5YsvLlxU6d4r5ALdB03y6bf9Q",
//	      "email": "xU6d4r5@ALdB03y6bf9Q.com"
//	    },
//	    {
//	      "id": "usr_2O5YsrLH9ViyyxLzqFNuSXxBDpV",
//	      "email": "9ViyyxL@zqFNuSXxBDpV.com"
//	    }
//	  ],
//	  "Products": [
//	    {
//	      "id": "prod_2O5Yst3NQSf8b6xBQCVgi4KJw6p",
//	      "creatorId": "usr_2O5YsvLlxU6d4r5ALdB03y6bf9Q",
//	      "createTime": "2023-04-06T21:01:59.752638+02:00"
//	    }
//	  ]
//	}
//
// For example:
// cat ./data/example1.json | go run ./tools/creatorclient --creators_svc localhost:9070
//
// output:
// The top creators are:
// 1: NfDynPx@oCsT5hsulDwM.com (products: 10, most recent creation: 2023-04-06T21:01:59.752638+02:00)
// 2: hnTwMre@zr0GUSVXxq7x.com (products: 10, most recent creation: 2023-04-05T13:54:59.752746+02:00)
// 3: ac6kaPa@F1TwzsLrWe5j.com (products: 10, most recent creation: 2023-04-05T06:52:59.752688+02:00)
// ---
package main
