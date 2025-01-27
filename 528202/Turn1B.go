package main

Foundation

func benchmarkFunctions(goVersion: String, functions: [() -> Void]) {
let startTime = DispatchTime.now()

for function in functions {
function()
}

let endTime = DispatchTime.now()
let timeElapsed = Double(endTime.uptimeNanoseconds - startTime.uptimeNanoseconds) / 1_000_000

print("Benchmark Results for Go Version $goVersion):")
print("Execution Time: $timeElapsed) milliseconds")
}