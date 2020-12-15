Task Init {
    $trashFolder = Join-Path -Path ".." -ChildPath ".trash"
    $trashSubFolder = Get-Date -Format 'yyyyMMddHHmmss'
    $script:trashFolder = Join-Path -Path $trashFolder -ChildPath $trashSubFolder
    New-Item -Path $script:trashFolder -ItemType Directory | Out-Null
    $script:trashFolder = Resolve-Path -Path $script:trashFolder
    $script:rootFolder = Resolve-Path -Path ".." -Relative
}

Task Clean -Depends Init {
    Exec { go clean } -workingDirectory $script:rootFolder
}

Task Format -Depends Clean {
    Exec { go fmt .\cmd\famulus\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\app\famulus\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\pkg\people\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\pkg\phone\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\pkg\scraper\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\pkg\famulus\cmd\collect } -workingDirectory $script:rootFolder
    Exec { go fmt .\pkg\famulus\cmd\root } -workingDirectory $script:rootFolder
    Exec { go fmt .\test\ } -workingDirectory $script:rootFolder
}

Task UpdateVendors -depends Clean {
    Exec { go mod vendor }
}

Task DownloadModules -depends UpdateVendors {
    Exec { go mod download }
}

Task TidyModules -depends DownloadModules, Format {
    Exec { go mod tidy }
}

Task Build -depends TidyModules {
    Exec { go build } -workingDirectory ..\cmd\famulus
}

Task Test -depends Build {
    Exec { go test } -workingDirectory ..\test\
}
