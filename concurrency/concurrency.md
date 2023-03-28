# Go Concurrency Patterns

## What is concurrency

composition of independently executing computations.
Way to structure software, particulary as a way to write clean code that interacts well with the real world.

## Goroutines

It's and independently executing function, launched by a go statement.

It has its own call stack, which grows and shrinks as requred.

It is very cheap.

It is not a thread.

## Channels

A channels in Go provides a connection between two goroutines, allowing them to communicate.

## Synchronization

when the main function executes <-c, it will wait for a value to be sent.

Similarly, when the boring function executes c <- value, it waits for a receiver to be ready.

A sender and receiver must both be ready to play their part in the communication. Otherwise we wait until the are.

Thus channels both communicate and synchronize.

## Patterns

### Generator: function that returns a channel

Channels are first-class values.

## Select

A control structure unique to concurrency
The reason channels and goroutines are built into the language.
