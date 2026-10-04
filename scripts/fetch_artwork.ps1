# fetch_artwork.ps1 - Complete artwork and metadata pipeline for albums and artists

param (
    [string]$OutputDir = "assets/images"
)

New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null

$albumQueries = @(
    @{ name="nayakan"; query="Nayakan" },
    @{ name="thalapathi"; query="Thalapathi" },
    @{ name="roja"; query="Roja" },
    @{ name="how_to_name_it"; query="How To Name It" },
    @{ name="nothing_but_wind"; query="Nothing But Wind" }
)

Write-Host "[+] Fetching Album Covers from Deezer Search API..."
foreach ($a in $albumQueries) {
    $apiUrl = "https://api.deezer.com/search/album?q=" + $a.query
    $jsonRaw = curl.exe -s $apiUrl
    try {
        $data = $jsonRaw | ConvertFrom-Json
        if ($data.data.Count -gt 0) {
            $imgUrl = $data.data[0].cover_big
            $target = Join-Path $OutputDir ($a.name + ".jpg")
            Write-Host ("[+] Album Cover: " + $a.name + " -> " + $imgUrl)
            curl.exe -s -L $imgUrl -o $target
        }
    } catch {
        Write-Host ("[-] Failed downloading album cover: " + $a.name)
    }
}

$artistQueries = @(
    @{ name="ilaiyaraaja"; query="Ilaiyaraaja" },
    @{ name="ar_rahman"; query="A.R.+Rahman" },
    @{ name="tyagaraja"; query="Tyagaraja" },
    @{ name="ms_subbulakshmi"; query="M.S.+Subbulakshmi" }
)

Write-Host "[+] Fetching Artist Portraits from Deezer Search API..."
foreach ($art in $artistQueries) {
    $apiUrl = "https://api.deezer.com/search/artist?q=" + $art.query
    $jsonRaw = curl.exe -s $apiUrl
    try {
        $data = $jsonRaw | ConvertFrom-Json
        if ($data.data.Count -gt 0) {
            $imgUrl = $data.data[0].picture_big
            $target = Join-Path $OutputDir ($art.name + ".jpg")
            Write-Host ("[+] Artist Portrait: " + $art.name + " -> " + $imgUrl)
            curl.exe -s -L $imgUrl -o $target
        }
    } catch {
        Write-Host ("[-] Failed downloading artist photo: " + $art.name)
    }
}

Write-Host ("[+] Complete Downloaded Assets in " + $OutputDir)
Get-ChildItem $OutputDir
