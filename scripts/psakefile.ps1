Task DownloadModules {
    Exec { go mod download }
}

Task TidyModules -depends DownloadModules {
    Exec { go mod tidy }
}

Task Build -depends TidyModules {
    Exec { go build } -workingDirectory ..\cmd\famulus
}
