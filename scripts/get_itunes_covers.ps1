# get_itunes_covers.ps1

$terms = @("Nayakan", "Thalapathi", "How To Name It", "Nothing But Wind")

foreach ($t in $terms) {
    $cleanName = ($t.ToLower() -replace " ", "_") + ".jpg"
    $url = "https://itunes.apple.com/search?term=" + [System.Uri]::EscapeDataString($t) + "&entity=album"
    $json = curl.exe -s $url
    try {
        $data = $json | ConvertFrom-Json
        if ($data.results.Count -gt 0) {
            $img = $data.results[0].artworkUrl100 -replace "100x100bb", "600x600bb"
            Write-Host ("[+] iTunes Cover: " + $t + " -> " + $img)
            curl.exe -s -L $img -o ("assets/images/" + $cleanName)
            curl.exe -s -L $img -o ("C:\Users\zeno7\.gemini\antigravity\brain\0a4a23c9-b4a9-4ca8-bc3c-6fb1f5d5cc12\assets\images\" + $cleanName)
        }
    } catch {
        Write-Host ("[-] Error on " + $t)
    }
}
