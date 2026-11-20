param(
  [Parameter(Mandatory = $true)][string] $Executable,
  [Parameter(Mandatory = $true)][string] $EvidenceDirectory,
  [switch] $Demo,
  [switch] $Running
)
$ErrorActionPreference = 'Stop'
if ($Demo -and $Running) { throw 'Choose the completed demo (-Demo) or running demo (-Running), not both.' }
$fixture = $Demo -or $Running
Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public static class SpeedDiscSmoke {
    [StructLayout(LayoutKind.Sequential)]
    public struct RECT { public int Left, Top, Right, Bottom; }
    [StructLayout(LayoutKind.Sequential)]
    public struct HIGHCONTRAST { public uint Size, Flags; public IntPtr Scheme; }
    [StructLayout(LayoutKind.Sequential)]
    public struct SCROLLBARINFO {
        public uint Size;
        public RECT Bounds;
        public int LineButton, ThumbTop, ThumbBottom, Reserved;
        public uint State0, State1, State2, State3, State4, State5;
    }
    [DllImport("user32.dll", SetLastError=true)]
    public static extern bool GetScrollBarInfo(IntPtr window, int idObject, ref SCROLLBARINFO info);
    [DllImport("user32.dll", EntryPoint="SystemParametersInfoW", ExactSpelling=true, SetLastError=true)]
    public static extern bool QueryHighContrast(uint action, uint size, ref HIGHCONTRAST info, uint flags);
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
    [DllImport("user32.dll", EntryPoint="SendMessageTimeoutW", ExactSpelling=true, SetLastError=true)]
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
if ($Running) { $params.ArgumentList = @("--demo-running", "--lang=en") }
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
    if ($Running -and $app.MainWindowTitle -notmatch 'Portable workbench') {
        throw "Running demo did not launch in English: $($app.MainWindowTitle)"
    }
    $readControl = { param([IntPtr] $control)
        # WM_GETTEXT is marshalled by the window manager for standard controls.
        # GetWindowText cannot read another process's edit/button text reliably.
        $capacity = 1024
        $bytes = [Runtime.InteropServices.Marshal]::AllocHGlobal($capacity * 2)
        try {
            [UIntPtr] $count = [UIntPtr]::Zero
            [Runtime.InteropServices.Marshal]::WriteInt16($bytes, 0)
            if ([SpeedDiscSmoke]::SendMessageTimeout($control, 0x000D, [UIntPtr]::new([uint32]$capacity), $bytes, 2, 2000, [ref]$count) -eq [IntPtr]::Zero) {
                throw 'The native control did not respond to WM_GETTEXT.'
            }
            $length = $count.ToUInt64()
            if ($length -ge $capacity) { throw "Invalid native text length ($length)." }
            [Runtime.InteropServices.Marshal]::PtrToStringUni($bytes, [int]$length)
        } finally {
            [Runtime.InteropServices.Marshal]::FreeHGlobal($bytes)
        }
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
    $scrollFact = 'Branded rail exercise runs in the completed compact demo only.'
    if ($fixture) {
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
            $controlIDs = @(109, 104, 105)
            if ($Running) { $controlIDs += 110 }
            foreach ($controlId in $controlIDs) {
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
                switch ($controlId) {
                    109 {
                        $resultHeight = $controlHeight / $scale
                        if ($resultHeight -lt 55) {
                            throw "Compact result edit is only $resultHeight logical units tall; expected at least 55."
                        }
                        $visibleText = & $readControl $control
                        if ($Demo -and $visibleText -notmatch '(?i)Demonstração|fictíci|Disco principal') {
                            throw "Portuguese synthetic report text not found: $visibleText"
                        }
                        if ($Running -and ($visibleText -notmatch '(?i)step' -or $visibleText -notmatch 'Checking Windows components' -or $visibleText -notmatch '2\s*(?:/|of)\s*3' -or $visibleText -notmatch '42\s*%')) {
                            throw "Expected English running text with step 2 of 3 and 42%: $visibleText"
                        }
                    }
                    110 {
                        # PBM_GETPOS takes no cross-process pointer and returns the current position.
                        [UIntPtr] $position = [UIntPtr]::Zero
                        if ([SpeedDiscSmoke]::SendMessageTimeout($control, 0x0408, [UIntPtr]::Zero, [IntPtr]::Zero, 2, 2000, [ref]$position) -eq [IntPtr]::Zero -or $position.ToUInt64() -ne 42) {
                            throw "Expected running progress 42%; native bar reports $($position.ToUInt64())%."
                        }
                    }
                    default {
                        $format = if ($controlId -eq 104) { 'TXT' } else { 'PDF' }
                        $expected = if ($Running) { "Export $format" } else { "Exportar $format" }
                        if ((& $readControl $control) -ne $expected) { throw "Expected localised button '$expected'." }
                        if ($Running) {
                            if ([SpeedDiscSmoke]::IsWindowEnabled($control) -or $tabOrder.Contains($control)) {
                                throw "Export #$controlId is active while the app displays a running operation."
                            }
                        } elseif (-not [SpeedDiscSmoke]::IsWindowEnabled($control) -or
                                  -not ([SpeedDiscSmoke]::GetWindowLongW($control, -16) -band 0x00010000) -or
                                  -not $tabOrder.Contains($control)) {
                            throw "Compact export control #$controlId is disabled or unreachable by Tab."
                        }
                    }
                }
            }
            $check = if ($Running) { 'English running detail step 2 of 3 at 42%, native progress visible at 42%, both exports visible but disabled' } else { 'Portuguese report readable, both localised exports visible, enabled and in native Tab order' }
            $compactFact = "Compact synthetic demo: ${actualW}x${actualH} logical client at runner DPI $dpi; result height $([Math]::Round($resultHeight)) logical; $check."
            if ($Demo) {
                $hc = [SpeedDiscSmoke+HIGHCONTRAST]::new()
                $hc.Size = [Runtime.InteropServices.Marshal]::SizeOf($hc)
                if (-not [SpeedDiscSmoke]::QueryHighContrast(0x42, $hc.Size, [ref]$hc, 0)) {
                    throw 'Cannot query the high contrast setting.'
                }
                $highContrast = ($hc.Flags -band 1) -ne 0
                $sendPointerless = {
                    param([IntPtr] $hwnd, [uint32] $message, [uint32] $wparam = 0, [int] $lparam = 0)
                    [UIntPtr] $reply = [UIntPtr]::Zero
                    if ([SpeedDiscSmoke]::SendMessageTimeout($hwnd, $message, [UIntPtr]::new($wparam), [IntPtr]::new($lparam), 2, 2000, [ref]$reply) -eq [IntPtr]::Zero) {
                        throw "Native scroll message failed: $message"
                    }
                    $reply.ToUInt64()
                }
                foreach ($testPair in @(@(102,112,0x018E), @(109,113,0x00CE))) {
                    $target = [SpeedDiscSmoke]::GetDlgItem($handle, $testPair[0])
                    $rail = [SpeedDiscSmoke]::GetDlgItem($handle, $testPair[1])
                    $bar = [SpeedDiscSmoke+SCROLLBARINFO]::new()
                    $bar.Size = [Runtime.InteropServices.Marshal]::SizeOf($bar)
                    if (-not [SpeedDiscSmoke]::GetScrollBarInfo($target, -5, [ref]$bar) -or $bar.Bounds.Bottom -le $bar.Bounds.Top -or ($bar.State0 -band 0x18001)) {
                        throw "Expected a real, scrollable vertical native bar for control $($testPair[0])."
                    }
                    if ($highContrast) {
                        if ($rail -ne [IntPtr]::Zero -and [SpeedDiscSmoke]::IsWindowVisible($rail)) {
                            throw 'High contrast must leave the native scrollbar unobstructed.'
                        }
                        continue
                    }
                    $railWait = [Diagnostics.Stopwatch]::StartNew()
                    $railBounds = [SpeedDiscSmoke+RECT]::new()
                    do {
                        if ($rail -ne [IntPtr]::Zero -and [SpeedDiscSmoke]::IsWindowVisible($rail) -and
                            [SpeedDiscSmoke]::GetWindowRect($rail, [ref]$railBounds) -and
                            $railBounds.Left -eq $bar.Bounds.Left -and $railBounds.Top -eq $bar.Bounds.Top -and
                            $railBounds.Right -eq $bar.Bounds.Right -and $railBounds.Bottom -eq $bar.Bounds.Bottom) { break }
                        Start-Sleep -Milliseconds 40
                    } while ($railWait.ElapsedMilliseconds -lt 2000)
                    if ($railWait.ElapsedMilliseconds -ge 2000) {
                        throw "Branded rail $($testPair[1]) does not match actual native scrollbar geometry."
                    }
                    if ($tabOrder.Contains($rail)) { throw 'A visual rail must not become a new Tab stop.' }
                    $queryTop = [uint32]$testPair[2]
                    $before = & $sendPointerless $target $queryTop
                    $point = (($railBounds.Bottom - $railBounds.Top - 1) -shl 16) -bor 2
                    $null = & $sendPointerless $rail 0x0201 1 $point # WM_LBUTTONDOWN on blank track below thumb
                    $null = & $sendPointerless $rail 0x0202 0 $point # WM_LBUTTONUP
                    $after = & $sendPointerless $target $queryTop
                    if ($after -le $before) { throw "Branded rail did not move control $($testPair[0]): $before -> $after." }
                    $null = & $sendPointerless $target 0x0115 6 # WM_VSCROLL / SB_TOP
                    if ((& $sendPointerless $target $queryTop) -ne 0) { throw 'Could not reset scroll position for screenshot.' }
                }
                $scrollFact = if ($highContrast) { 'High contrast active: both native bars unobstructed.' } else { 'Branded rails match both native scrollbar bounds; track clicks moved list and result, and neither rail changed Tab order.' }
            }

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
    $screenshotName = if ($Running) { 'speeddisc-gui-running-en-compact.png' } elseif ($Demo) { 'speeddisc-gui-report-pt-compact.png' } else { 'speeddisc-gui.png' }
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
        "Screenshot: $screenshotStatus; synthetic completed demo: $Demo; synthetic running demo: $Running"
        $compactFact
        $scrollFact
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
