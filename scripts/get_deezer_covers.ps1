# get_deezer_covers.ps1

$items = @("Nayakan", "Thalapathi", "Roja", "How To Name It", "Nothing But Wind")

foreach ($item in $items) {
    $cleanName = ($item.ToLower() -replace " ", "_") + ".jpg"
    $url = "https://api.deezer.com/search/album?q=" + [System.Uri]::EscapeDataString($item)
    $json = curl.exe -s $url
    try {
        $data = $json | ConvertFrom-Json
        if ($data.data.Count -gt 0) {
            $img = $data.data[0].cover_big
            Write-Host ("[+] Fetching " + $item + " -> " + $img)
            curl.exe -s -L $img -o ("assets/images/" + $cleanName)
            curl.exe -s -L $img -o ("C:\Users\zeno7\.gemini\antigravity\brain\0a4a23c9-b4a9-4ca8-bc3c-6fb1f5d5cc12\assets\images\" + $cleanName)
        }
    } catch {
        Write-Host ("[-] Error on " + $item)
    }
}
