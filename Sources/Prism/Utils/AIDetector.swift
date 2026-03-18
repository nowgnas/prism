import Foundation

/// Lightweight AI process detector.
/// Polls `ps` every 3 seconds — minimal CPU/memory impact.
enum AIDetector {
    static let knownNames: Set<String> = [
        "claude", "aider", "llm", "ollama",
        "sgpt", "gpt4all", "chatgpt", "gemini", "copilot",
    ]

    /// Returns names of child processes for the given parent PID.
    static func childNames(of ppid: pid_t) -> [String] {
        let task = Process()
        task.executableURL = URL(fileURLWithPath: "/bin/ps")
        task.arguments = ["-ax", "-o", "ppid=,comm="]
        let pipe = Pipe()
        task.standardOutput = pipe
        task.standardError = Pipe()
        guard (try? task.run()) != nil else { return [] }
        let data = pipe.fileHandleForReading.readDataToEndOfFile()
        task.waitUntilExit()
        guard let output = String(data: data, encoding: .utf8) else { return [] }

        var names: [String] = []
        for line in output.split(separator: "\n") {
            let parts = line.split(separator: " ", maxSplits: 1, omittingEmptySubsequences: true)
            guard parts.count == 2,
                  let parentPID = pid_t(parts[0].trimmingCharacters(in: .whitespaces)),
                  parentPID == ppid else { continue }
            let comm = String(parts[1]).trimmingCharacters(in: .whitespaces)
            names.append(URL(fileURLWithPath: comm).lastPathComponent.lowercased())
        }
        return names
    }

    static func isAI(_ processName: String) -> Bool {
        knownNames.contains(processName.lowercased())
    }
}
