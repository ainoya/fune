/*
Package integration implements tests built upon embedded fune, and focus on
fune correctness.

Features/goals of the integration tests:
1. test the whole code base except command-line parsing.
2. check internal data, including raft, store and etc.
3. based on goroutines, which is faster than process.
4. mainly tests user behavior and user-facing API.
*/

package integration
