// Package log implements a simple structure logging package. It defines a type,
// Context, with methods for formatting output. It also has a predefined
// 'standard' Context accessible through helper functions Info[f], Warning[f],
// Error[f] and Debug[f], which are easier to use than creating a Context
// manually. That Context writes to standard error and prints the date and time
// of each logged message.
//
// In fact, this is another "bike" for output logs, only with blackjack and
// hookers.
package log
