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

Task UpdateVendors -depends Clean {
    Exec { go mod vendor }
}

Task DownloadModules -depends UpdateVendors {
    Exec { go mod download }
}

Task Format -Depends DownloadModules {
    Exec { go fmt .\cmd\famulus\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\app\famulus\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\pkg\people\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\pkg\phone\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\internal\pkg\scraper\ } -workingDirectory $script:rootFolder
    Exec { go fmt .\pkg\famulus\cmd\collect } -workingDirectory $script:rootFolder
    Exec { go fmt .\pkg\famulus\cmd\root } -workingDirectory $script:rootFolder
    Exec { go fmt .\test\ } -workingDirectory $script:rootFolder
}

Task TidyModules -depends DownloadModules, Format {
    Exec { go mod tidy }
}

Task PreBuild -Depends TidyModules {
    $script:publishFolder = Join-Path -Path $script:trashFolder -ChildPath "bin"
    
    New-Item -Path $script:publishFolder -ItemType Directory | Out-Null
}

Task BuildWinx64 -Depends PreBuild {
    $script:publishWinx64Folder = Join-Path -Path $script:publishFolder -ChildPath "winx64"
    $outputFile = Join-Path -Path $script:publishWinx64Folder -ChildPath "famulus.exe"

    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile .\cmd\famulus } -workingDirectory $script:rootFolder
}

Task BuildWinx86 -Depends PreBuild {
    $script:publishWinx86Folder = Join-Path -Path $script:publishFolder -ChildPath "winx86"
    $outputFile = Join-Path -Path $script:publishWinx86Folder -ChildPath "famulus.exe"
    
    $env:GOOS = "windows"
    $env:GOARCH = "386"
    Exec { go build -o $outputFile .\cmd\famulus } -workingDirectory $script:rootFolder
}

Task BuildLinux64 -Depends PreBuild {
    $script:publishLinux64Folder = Join-Path -Path $script:publishFolder -ChildPath "linux64"
    $outputFile = Join-Path -Path $script:publishLinux64Folder -ChildPath "famulus"

    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile .\cmd\famulus } -workingDirectory $script:rootFolder
}

Task Build -Depends BuildWinx64, BuildWinx86, BuildLinux64

Task Test -depends Build {
    $env:GOOS = ""
    $env:GOARCH = ""
    Exec { go test .\test\ } -workingDirectory $script:rootFolder
}
