package main

import (
	"fmt"
)

// TODO: implement REST API production rdy
// * TLS?
// * CORS?
// * database sqlite?
// * env config
// * logger stdout?
// * rate limiter
// * doc + README

// TODO: /version
// * print API version (for test purpose)

// TODO: /fizzbuzz (public)
// * five parameters, 3 int (int1, int2, limit) & 2 string (str1, str2)
// * int1
// * int2
// * limit (go from 1 to limit)
// * return list from 1 to limit where:
//      - all multiples of int1 are replaced by str1
//      - all multiples of int2 are replaced by str2
//      - all multiples of int1+int2 replaced by str1+str2
// * with tests

// TODO: /stats
// * return:
//      - the most submited sequence
//      - the number of hits
// * with tests
