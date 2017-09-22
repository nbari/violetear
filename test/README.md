##  Identify allocations

When setting allocfreetrace=1 to environment variable GODEBUG, you can see
stacktrace where allocation occures:

    allocfreetrace: setting allocfreetrace=1 causes every allocation to be
    profiled and a stack trace printed on each object's allocation and free.

Build test:

    go test -c

Run bench:

    GODEBUG=allocfreetrace=1 ./test.test -test.run=none -test.bench=BenchmarkRouter -test.benchtime=10ms 2>trace.log


https://methane.github.io/2015/02/reduce-allocation-in-go-code/
