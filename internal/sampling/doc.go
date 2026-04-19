// Package sampling implements pluggable log-line sampling strategies
// for logpipe pipelines.
//
// Two strategies are provided:
//
//   - Rate sampling: keeps every Nth line in a deterministic round-robin
//     fashion, useful for predictable volume reduction.
//
//   - Random sampling: keeps each line independently with a given
//     probability, useful for statistically representative samples.
//
// Use sampling.New(rate) or sampling.NewRandom(prob) to obtain a Sampler,
// then call Sampler.Sample() per line to decide whether to forward it.
package sampling
