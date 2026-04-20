// Package aggregator provides streaming aggregation over structured JSON log lines.
//
// An Aggregator tracks a single numeric field across multiple log lines and
// computes one of four operations: count, sum, min, or max.
//
// Example:
//
//	a, err := aggregator.New("duration_ms", aggregator.OpSum)
//	if err != nil {
//		log.Fatal(err)
//	}
//	a.Add([]byte(`{"duration_ms": 42}`))
//	a.Add([]byte(`{"duration_ms": 58}`))
//	fmt.Println(a.Result()) // 100
//
// Aggregators are safe for concurrent use.
package aggregator
