package utils

import (
	"log"
	"strconv"
	"time"
)

// Expire helper function to convert expire timestamp (int) into seconds since epoch
func Expire(expire int) int64 {
	tstamp := strconv.Itoa(expire)
	if len(tstamp) == 10 {
		return int64(expire)
	}
	return int64(time.Now().Unix() + int64(expire))
}

// UnixTime helper function to convert given time into Unix timestamp
func UnixTime(ts string) int64 {
	// time is unix since epoch
	if len(ts) == 10 { // unix time
		tstamp, _ := strconv.ParseInt(ts, 10, 64)
		return tstamp
	}
	// YYYYMMDD, always use 2006 as year 01 for month and 02 for date since it is predefined int Go parser
	const layout = "20060102"
	t, err := time.Parse(layout, ts)
	if err != nil {
		log.Printf("unable to parse, error %v\n", err)
		return 0
	}
	return int64(t.Unix())
}

// Unix2Time helper function to convert given time into Unix timestamp
func Unix2Time(ts int64) string {
	// YYYYMMDD, always use 2006 as year 01 for month and 02 for date since it is predefined int Go parser
	const layout = "20060102"
	t := time.Unix(ts, 0)
	return t.In(time.UTC).Format(layout)
}
