import Foundation
import AVFoundation
import AppKit

let args = CommandLine.arguments
let source = URL(fileURLWithPath: args[1])
let destination = URL(fileURLWithPath: args[2])
let images = try FileManager.default.contentsOfDirectory(at: source, includingPropertiesForKeys: nil).filter { $0.pathExtension == "png" }.sorted { $0.lastPathComponent < $1.lastPathComponent }
guard let first = images.first, let image = NSImage(contentsOf: first) else { fatalError("No preview frames") }
var bounds = CGRect(origin: .zero, size: image.size)
guard let firstCG = image.cgImage(forProposedRect: &bounds, context: nil, hints: nil) else { fatalError("Image decoding") }
let width = firstCG.width, height = firstCG.height
let writer = try AVAssetWriter(outputURL: destination, fileType: .mp4)
let input = AVAssetWriterInput(mediaType: .video, outputSettings: [AVVideoCodecKey: AVVideoCodecType.h264, AVVideoWidthKey: width, AVVideoHeightKey: height])
let adapter = AVAssetWriterInputPixelBufferAdaptor(assetWriterInput: input, sourcePixelBufferAttributes: [kCVPixelBufferPixelFormatTypeKey as String: kCVPixelFormatType_32ARGB, kCVPixelBufferWidthKey as String: width, kCVPixelBufferHeightKey as String: height])
writer.add(input)
guard writer.startWriting() else { fatalError("Video start failed") }
writer.startSession(atSourceTime: .zero)
for (index, url) in images.enumerated() {
  while !input.isReadyForMoreMediaData { Thread.sleep(forTimeInterval: 0.01) }
  try autoreleasepool {
    guard let img = NSImage(contentsOf: url) else { throw NSError(domain: "Preview", code: 1) }
    var rect = CGRect(origin: .zero, size: img.size)
    guard let cg = img.cgImage(forProposedRect: &rect, context: nil, hints: nil) else { throw NSError(domain: "Preview", code: 2) }
    var buffer: CVPixelBuffer?
    guard CVPixelBufferPoolCreatePixelBuffer(nil, adapter.pixelBufferPool!, &buffer) == kCVReturnSuccess, let pixel = buffer else { throw NSError(domain: "Preview", code: 3) }
    CVPixelBufferLockBaseAddress(pixel, [])
    let context = CGContext(data: CVPixelBufferGetBaseAddress(pixel), width: width, height: height, bitsPerComponent: 8, bytesPerRow: CVPixelBufferGetBytesPerRow(pixel), space: CGColorSpaceCreateDeviceRGB(), bitmapInfo: CGImageAlphaInfo.noneSkipFirst.rawValue)!
    context.draw(cg, in: CGRect(x: 0, y: 0, width: width, height: height))
    CVPixelBufferUnlockBaseAddress(pixel, [])
    guard adapter.append(pixel, withPresentationTime: CMTime(value: Int64(index), timescale: 10)) else { throw NSError(domain: "Preview", code: 4) }
  }
}
input.markAsFinished()
let completion = DispatchSemaphore(value: 0)
writer.finishWriting { completion.signal() }
completion.wait()
guard writer.status == .completed else { fatalError("Video encode failed: \(String(describing: writer.error))") }
print("Encoded \(images.count) frames to \(destination.path)")
