//go:build prob_generate

package fuzzer

// 0.1% chance to trigger an a generate request each time we fetch a new request to execute from the priority queue
const PROB_GENERATE_PROC = 0.001
