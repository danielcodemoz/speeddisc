param(
  [Parameter(Mandatory = $true)][string] $Executable,
  [Parameter(Mandatory = $true)][string] $EvidenceDirectory,
  [switch] $Demo
)
$ErrorActionPreference = 'Stop'
Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public static class SpeedDiscSmoke {
    [StructLayout(LayoutKind.Sequential)]
    public struct RECT { public int Left, Top, Right, Bottom; }
    [DllImport("user32.dll", SetLastError=true)]
    public static extern bool GetWindowRect(IntPtr window, out RECT rect);
    [StructLayout(LayoutKind.Sequential)]
    public struct POINT { public int X, Y; }
    [DllImport("user32.dll", SetLastError=true)]
    public static extern IntPtr SetThreadDpiAwarenessContext(IntPtr context);
    [DllImport("user32.dll")]
    public static extern uint GetDpiForWindow(IntPtr window);
    [DllImport("user32.dll", SetLastError=true)]
    public static extern bool GetClientRect(IntPtr window, out RECT rect);
    [DllImport("user32.dll", SetLastError=true)]
    public static extern bool ClientToScreen(IntPtr window, ref POINT point);
    [DllImport("user32.dll", SetLastError=true)]
    public static extern bool SetWindowPos(IntPtr window, IntPtr insertAfter, int x, int y, int width, int height, uint flags);
    [DllImport("user32.dll")]
    public static extern bool IsWindowEnabled(IntPtr window);
    [DllImport("user32.dll")]
    public static extern int GetWindowLongW(IntPtr window, int index);
    [DllImport("user32.dll")]
    public static extern IntPtr GetNextDlgTabItem(IntPtr window, IntPtr control, bool previous);
    [DllImport("user32.dll", SetLastError=true)]
    public static extern uint GetWindowThreadProcessId(IntPtr window, out uint pid);
    [DllImport("user32.dll")]
    public static extern bool IsWindowVisible(IntPtr window);
    [DllImport("user32.dll")]
    public static extern IntPtr GetDlgItem(IntPtr window, int id);
    [DllImport("user32.dll", SetLastError=true)]
    public static extern IntPtr SendMessageTimeout(IntPtr window, uint message, UIntPtr w, IntPtr l, uint flags, uint timeout, out UIntPtr result);
    [DllImport("user32.dll", SetLastError=true)]
    public static extern bool PrintWindow(IntPtr window, IntPtr buffer, uint flags);
}
'@
$exe = (Resolve-Path -LiteralPath $Executable).Path
$evidence = (New-Item -ItemType Directory -Force -Path $EvidenceDirectory).FullName
$identity = [Security.Principal.WindowsIdentity]::GetCurrent()
$principal = [Security.Principal.WindowsPrincipal]::new($identity)
$adminToken = $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
$timer = [Diagnostics.Stopwatch]::StartNew()
$params = @{ FilePath = $exe; WorkingDirectory = $evidence; PassThru = $true }
if ($Demo) { $params.ArgumentList = @("--demo") }
$app = Start-Process @params
try {
    $handle = [IntPtr]::Zero
    while ($timer.Elapsed.TotalSeconds -lt 20) {
        $app.Refresh()
        if ($app.HasExited) { throw "GUI exited before opening a window (code $($app.ExitCode))." }
        $handle = $app.MainWindowHandle
        if ($handle -ne [IntPtr]::Zero -and [SpeedDiscSmoke]::IsWindowVisible($handle)) { break }
        Start-Sleep -Milliseconds 100
    }
    if ($handle -eq [IntPtr]::Zero -or -not [SpeedDiscSmoke]::IsWindowVisible($handle)) {
        throw 'No visible top-level GUI window appeared within 20 seconds.'
    }
    [uint32] $windowPid = 0
    $null = [SpeedDiscSmoke]::GetWindowThreadProcessId($handle, [ref]$windowPid)
    if ($windowPid -ne $app.Id -or $app.MainWindowTitle -notmatch '^SpeedDisc') {
        throw "Unexpected top-level window: PID $windowPid; title $($app.MainWindowTitle)"
    }
    if (-not $app.WaitForInputIdle(5000) -or -not $app.Responding) {
        throw 'SpeedDisc window is not processing messages.'
    }
    $list = [SpeedDiscSmoke]::GetDlgItem($handle, 102)
    $ask = { param([uint32] $message)
        [UIntPtr] $value = [UIntPtr]::Zero
        if ($list -eq [IntPtr]::Zero -or [SpeedDiscSmoke]::SendMessageTimeout($list, $message, [UIntPtr]::Zero, [IntPtr]::Zero, 2, 2000, [ref]$value) -eq [IntPtr]::Zero) {
            throw 'The native tool list did not respond.'
        }
        return $value.ToUInt64()
    }
    $topIndex = & $ask 0x018E # LB_GETTOPINDEX
    $selectedIndex = & $ask 0x0188 # LB_GETCURSEL
    if ($topIndex -ne 0 -or $selectedIndex -ne 2) {
        throw "Unexpected initial list state: top=$topIndex selected=$selectedIndex; expected top 0 and Só ver (index 2)."
    }
    $rect = [SpeedDiscSmoke+RECT]::new()
    if (-not [SpeedDiscSmoke]::GetWindowRect($handle, [ref]$rect)) { throw 'Unable to read window dimensions.' }
    $compactFact = 'Compact sizing not exercised in the default launch.'
    if ($Demo) {
        # The app is per-monitor DPI aware. Disable coordinate virtualization on this
        # caller thread while measuring its window and children in screen pixels.
        $priorDpiContext = [SpeedDiscSmoke]::SetThreadDpiAwarenessContext([IntPtr]::new(-4))
        if ($priorDpiContext -eq [IntPtr]::Zero) { throw 'Cannot use per-monitor DPI coordinates.' }
        try {
            $dpi = [SpeedDiscSmoke]::GetDpiForWindow($handle)
            if ($dpi -lt 96 -or $dpi -gt 384) { throw "Unexpected window DPI: $dpi" }
            $scale = $dpi / 96.0
            $client = [SpeedDiscSmoke+RECT]::new()
            if (-not [SpeedDiscSmoke]::GetWindowRect($handle, [ref]$rect) -or
                -not [SpeedDiscSmoke]::GetClientRect($handle, [ref]$client)) {
                throw 'Cannot measure the initial outer and client rectangles.'
            }
            $nonClientWidth = ($rect.Right - $rect.Left) - ($client.Right - $client.Left)
            $nonClientHeight = ($rect.Bottom - $rect.Top) - ($client.Bottom - $client.Top)
            if ($nonClientWidth -lt 0 -or $nonClientWidth -gt 160*$scale -or
                $nonClientHeight -lt 0 -or $nonClientHeight -gt 160*$scale) {
                throw "Unexpected frame dimensions: ${nonClientWidth}x${nonClientHeight}"
            }
            $wantedWidth = [int][Math]::Min(1040, [Math]::Round(($client.Right - $client.Left) / $scale))
            $targetWidth = [int][Math]::Round($wantedWidth * $scale) + $nonClientWidth
            $targetHeight = [int][Math]::Round(535 * $scale) + $nonClientHeight
            # SWP_NOMOVE | SWP_NOZORDER | SWP_NOACTIVATE: resize this window only.
            if (-not [SpeedDiscSmoke]::SetWindowPos($handle, [IntPtr]::Zero, 0, 0, $targetWidth, $targetHeight, 0x16)) {
                throw 'Cannot resize the SpeedDisc demo window.'
            }
            if (-not [SpeedDiscSmoke]::GetClientRect($handle, [ref]$client) -or
                -not [SpeedDiscSmoke]::GetWindowRect($handle, [ref]$rect)) {
                throw 'Cannot measure the resized SpeedDisc window.'
            }
            $actualW = [int][Math]::Round(($client.Right - $client.Left) / $scale)
            $actualH = [int][Math]::Round(($client.Bottom - $client.Top) / $scale)
            if ([Math]::Abs($actualW - $wantedWidth) -gt 5 -or [Math]::Abs($actualH - 535) -gt 5) {
                throw "Expected about ${wantedWidth}x535 logical client units; got ${actualW}x${actualH}."
            }
            $clientOrigin = [SpeedDiscSmoke+POINT]::new()
            if (-not [SpeedDiscSmoke]::ClientToScreen($handle, [ref]$clientOrigin)) {
                throw 'Cannot map the client origin to screen coordinates.'
            }
            $tabOrder = [Collections.Generic.HashSet[IntPtr]]::new()
            $next = [IntPtr]::Zero
            for ($i = 0; $i -lt 64; $i++) {
                $next = [SpeedDiscSmoke]::GetNextDlgTabItem($handle, $next, $false)
                if ($next -eq [IntPtr]::Zero -or -not $tabOrder.Add($next)) { break }
            }
            $resultHeight = 0.0
            foreach ($controlId in @(109, 104, 105)) {
                $control = [SpeedDiscSmoke]::GetDlgItem($handle, $controlId)
                $bounds = [SpeedDiscSmoke+RECT]::new()
                if ($control -eq [IntPtr]::Zero -or -not [SpeedDiscSmoke]::IsWindowVisible($control) -or
                    -not [SpeedDiscSmoke]::GetWindowRect($control, [ref]$bounds)) {
                    throw "Missing or invisible compact control #$controlId."
                }
                $controlWidth = $bounds.Right - $bounds.Left
                $controlHeight = $bounds.Bottom - $bounds.Top
                if ($controlWidth -le 0 -or $controlHeight -le 0 -or
                    $bounds.Left -lt $clientOrigin.X -or $bounds.Top -lt $clientOrigin.Y -or
                    $bounds.Right -gt ($clientOrigin.X + $client.Right) -or
                    $bounds.Bottom -gt ($clientOrigin.Y + $client.Bottom)) {
                    throw "Compact control #$controlId falls outside the client area."
                }
                if ($controlId -eq 109) {
                    $resultHeight = $controlHeight / $scale
                    if ($resultHeight -lt 55) {
                        throw "Compact result edit is only $resultHeight logical units tall; expected at least 55."
                    }
                } elseif (-not [SpeedDiscSmoke]::IsWindowEnabled($control) -or
                          -not ([SpeedDiscSmoke]::GetWindowLongW($control, -16) -band 0x00010000) -or
                          -not $tabOrder.Contains($control)) {
                    throw "Compact export control #$controlId is disabled or unreachable by Tab."
                }
            }
            $compactFact = "Compact synthetic demo: ${actualW}x${actualH} logical client at runner DPI $dpi; result height $([Math]::Round($resultHeight)) logical; both export controls inside client, enabled and in native Tab order."
        } finally {
            $null = [SpeedDiscSmoke]::SetThreadDpiAwarenessContext($priorDpiContext)
        }
    }
    $width = $rect.Right - $rect.Left
    $height = $rect.Bottom - $rect.Top
    if ($width -lt 320 -or $height -lt 240 -or $width -gt 8192 -or $height -gt 8192) {
        throw "Unexpected GUI dimensions: ${width}x${height}"
    }
    $caption = $app.MainWindowTitle
    $readyMs = $timer.ElapsedMilliseconds
    $workingBytes = $app.WorkingSet64
    $screenshotStatus = 'unavailable'
    $screenshotName = if ($Demo) { 'speeddisc-gui-compact.png' } else { 'speeddisc-gui.png' }
    # Print exactly this process's window; never capture the runner desktop or use destructive controls.
    try {
        Add-Type -AssemblyName System.Drawing.Common
        $bitmap = [Drawing.Bitmap]::new($width, $height)
        try {
            $graphics = [Drawing.Graphics]::FromImage($bitmap)
            try {
                $hdc = $graphics.GetHdc()
                try { $rendered = [SpeedDiscSmoke]::PrintWindow($handle, $hdc, 2) }
                finally { $graphics.ReleaseHdc($hdc) }
            } finally { $graphics.Dispose() }
            if ($rendered) {
                $colours = [Collections.Generic.HashSet[int]]::new()
                for ($x = 0; $x -lt $width; $x += [Math]::Max(1, [int]($width / 20))) {
                    for ($y = 0; $y -lt $height; $y += [Math]::Max(1, [int]($height / 20))) {
                        $null = $colours.Add($bitmap.GetPixel($x, $y).ToArgb())
                    }
                }
                if ($colours.Count -ge 8) {
                    $bitmap.Save((Join-Path $evidence $screenshotName), [Drawing.Imaging.ImageFormat]::Png)
                    $screenshotStatus = 'saved window-only PNG'
                }
            }
        } finally { $bitmap.Dispose() }
    } catch {
        Write-Warning "Window-only capture not available: $($_.Exception.Message)"
    }
    if ($screenshotStatus -eq 'unavailable') {
        Write-Warning 'PrintWindow did not produce a useful capture; visual review is still required.'
    }
    if (-not $app.CloseMainWindow() -or -not $app.WaitForExit(10000)) {
        throw 'SpeedDisc did not close cleanly on WM_CLOSE.'
    }
    if ($app.ExitCode -ne 0) { throw "SpeedDisc closed with exit code $($app.ExitCode)." }
    $facts = @(
        'Window launch, visible process-owned window, responsiveness and normal close: PASS'
        "Window: $caption / ${width}x${height}"
        "Initial list: top=$topIndex, selected=$selectedIndex (Só ver)"
        "User token administrator: $adminToken"
        "Time to responsive window on this runner (includes polling): ${readyMs} ms"
        "Working set on this runner: $workingBytes bytes"
        "Screenshot: $screenshotStatus; synthetic demo: $Demo"
        $compactFact
        'Compact sizing at runner DPI is not a 125%-DPI hardware or multi-monitor test.'
        'No maintenance operation, elevation request, system change or export was triggered.'
        'Windows Server smoke is not Windows 10/11 standard-user, USB, UAC, SmartScreen or performance acceptance.'
    )
    $facts | Set-Content -Encoding utf8 -LiteralPath (Join-Path $evidence 'gui-smoke.txt')
    if ($env:GITHUB_STEP_SUMMARY) { $facts | Add-Content -Encoding utf8 -LiteralPath $env:GITHUB_STEP_SUMMARY }
    $facts | Write-Output
} finally {
    $app.Refresh()
    if (-not $app.HasExited) {
        Stop-Process -Id $app.Id -Force
        $null = $app.WaitForExit(5000)
    }
    $app.Dispose()
}
