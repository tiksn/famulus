Task UpdateVendors {
    Exec { go mod vendor }
}

Task DownloadModules -depends UpdateVendors {
    Exec { go mod download }
}

Task TidyModules -depends DownloadModules {
    Exec { go mod tidy }
}

Task Build -depends TidyModules {
    Exec { go build } -workingDirectory ..\cmd\famulus
}
