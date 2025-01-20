import Foundation

func readFileSegment(_ filePath: String, _ offset: Int, _ length: Int) -> Data? {
do {
let fileHandle = try FileHandle(forReadingFrom: URL(fileURLWithPath: filePath))
defer {
fileHandle.closeFile()
}

let data = try fileHandle.readData(ofLength: length)

return data
} catch {
print("Error reading file: \(error.localizedDescription)")
return nil
}
}

let filePath = "/path/to/your/large/file"
let offset = 0 // Starting offset to read from
let length = 1024 // Number of bytes to read

if let data = readFileSegment(filePath, offset, length) {
print("Read \(data.count) bytes successfully.")
} else {
print("Failed to read file segment.")
}