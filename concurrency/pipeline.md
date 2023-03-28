# Pipeline

Pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function.

In each stage, the goroutines

- receive values from upstream via inbound channels
- perform some function on that data, usually producing new values
- send values downstream via outbound channels

Each stage has number of inbound and outbound channels, except the first and last stages, which have only outbound or inbound channels, respectively. The first stage is sometimes called the source or producer; the last stage, the sink or consumer.

## Fan-out, Fan-in

Multiple functions can read from the same channel until that channel is closed; this is called **fan-out**. This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O

A function can read from muliple inputs and process until all are closed by muliplexing the input channels onto a single channel that's closed when all the inputs are closed. This is called **fan-in**

## Pipeline function pattern

- stages close their outbound channels when all the send operations are done.
- stages keep receiving values from inbound channels until thos channels are closed.

this pattern allows each receiving stage to be written as range loop and ensures that all goroutines exit once all values have been successfully sent downstream.
