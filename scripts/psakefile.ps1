Task Init {
    $trashFolder = Join-Path -Path ".." -ChildPath ".trash"
    $trashSubFolder = Get-Date -Format 'yyyyMMddHHmmss'
    $script:trashFolder = Join-Path -Path $trashFolder -ChildPath $trashSubFolder
    New-Item -Path $script:trashFolder -ItemType Directory | Out-Null
    $script:trashFolder = Resolve-Path -Path $script:trashFolder
    $script:rootFolder = Resolve-Path -Path ".." -Relative
}

Task UpdateVendors -depends Init {
    Exec { go mod vendor }
}

Task DownloadModules -depends UpdateVendors {
    Exec { go mod download }
}

Task Clean -Depends UpdateVendors {
    Exec { go clean } -workingDirectory $script:rootFolder
}

Task Format -Depends Clean {
    Exec { go fmt ./cmd/famulus/ } -workingDirectory $script:rootFolder
    Exec { go fmt ./internal/app/famulus/ } -workingDirectory $script:rootFolder
    Exec { go fmt ./internal/pkg/people/ } -workingDirectory $script:rootFolder
    Exec { go fmt ./internal/pkg/phone/ } -workingDirectory $script:rootFolder
    Exec { go fmt ./internal/pkg/scraper/ } -workingDirectory $script:rootFolder
    Exec { go fmt ./pkg/famulus/cmd/collect } -workingDirectory $script:rootFolder
    Exec { go fmt ./pkg/famulus/cmd/root } -workingDirectory $script:rootFolder
    Exec { go fmt ./test/ } -workingDirectory $script:rootFolder
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
    $script:publishWinx64OutputFile = Join-Path -Path $script:publishWinx64Folder -ChildPath "famulus.exe"

    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    Exec { go build -o $script:publishWinx64OutputFile ./cmd/famulus } -workingDirectory $script:rootFolder
}

Task BuildWinx86 -Depends PreBuild {
    $script:publishWinx86Folder = Join-Path -Path $script:publishFolder -ChildPath "winx86"
    $script:publishWinx86OutputFile = Join-Path -Path $script:publishWinx86Folder -ChildPath "famulus.exe"
    
    $env:GOOS = "windows"
    $env:GOARCH = "386"
    Exec { go build -o $script:publishWinx86OutputFile ./cmd/famulus } -workingDirectory $script:rootFolder
}

Task BuildLinux64 -Depends PreBuild {
    $script:publishLinux64Folder = Join-Path -Path $script:publishFolder -ChildPath "linux64"
    $script:publishLinux64OutputFile = Join-Path -Path $script:publishLinux64Folder -ChildPath "famulus"

    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    Exec { go build -o $script:publishLinux64OutputFile ./cmd/famulus } -workingDirectory $script:rootFolder
}

Task Build -Depends BuildWinx64, BuildWinx86, BuildLinux64

Task Test -depends Build {
    $env:GOOS = ""
    $env:GOARCH = ""
    Exec { go test ./test/ } -workingDirectory $script:rootFolder
}

Task InstallLocally -depends Test {
    Exec { go install ./cmd/famulus } -workingDirectory $script:rootFolder
}

Task CollectArtifacts -depends Test, BuildWinx64, BuildWinx86, BuildLinux64 {
    $script:artifactsFolder = Join-Path -Path $script:trashFolder -ChildPath 'artifacts'
    New-Item -Path $script:artifactsFolder -ItemType Directory | Out-Null

    $rootCommandFile = Join-Path -Path $script:rootFolder -ChildPath 'pkg\famulus\cmd\root\root.go'
    $rootCommandFileContent = Get-Content -Path $rootCommandFile
    $appVersionLine = $rootCommandFileContent | Where-Object { $_.Contains('AppVersion') } | Where-Object { $_.Contains('=') }
    $appVersion = ($appVersionLine -split '=')[1].Trim()
    $appVersion = $appVersion.TrimStart('"')
    $appVersion = $appVersion.TrimEnd('"')
    Assert -conditionToCheck ($null -ne $appVersion) -failureMessage 'Version is not provided'

    $script:archiveWinx64 = Join-Path -Path $script:artifactsFolder -ChildPath "famulus-$appVersion-win-x64.zip"
    Compress-Archive -Path $script:publishWinx64OutputFile -DestinationPath $script:archiveWinx64

    $script:archiveWinx86 = Join-Path -Path $script:artifactsFolder -ChildPath "famulus-$appVersion-win-x86.zip"
    Compress-Archive -Path $script:publishWinx86OutputFile -DestinationPath $script:archiveWinx86

    $script:archiveLinux64 = Join-Path -Path $script:artifactsFolder -ChildPath "famulus-$appVersion-linux-x64.zip"
    Compress-Archive -Path $script:publishLinux64OutputFile -DestinationPath $script:archiveLinux64
}