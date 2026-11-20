param(
    [Parameter(Mandatory = $true)][string] $Executable,
    [Parameter(Mandatory = $true)][string] $EvidenceDirectory
)
$ErrorActionPreference = 'Stop'
Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
public static class SpeedDiscExportTest {
    private delegate bool EnumProc(IntPtr hwnd, IntPtr data);
    [DllImport("user32.dll")] private static extern bool EnumWindows(EnumProc cb, IntPtr data);
    [DllImport("user32.dll")] private static extern bool EnumChildWindows(IntPtr hwnd, EnumProc cb, IntPtr data);
    [DllImport("user32.dll")] private static extern uint GetWindowThreadProcessId(IntPtr hwnd, out uint pid);
    [DllImport("user32.dll")] private static extern IntPtr GetWindow(IntPtr hwnd, uint command);
    [DllImport("user32.dll")] private static extern bool IsWindowVisible(IntPtr hwnd);
    [DllImport("user32.dll")] public static extern bool IsWindow(IntPtr hwnd);
    [DllImport("user32.dll")] public static extern IntPtr GetDlgItem(IntPtr hwnd, int id);
    [DllImport("user32.dll", CharSet = CharSet.Unicode)] private static extern int GetClassName(IntPtr hwnd, StringBuilder b, int capacity);
    [DllImport("user32.dll", CharSet = CharSet.Unicode)] private static extern int GetWindowText(IntPtr hwnd, StringBuilder b, int capacity);
    [DllImport("user32.dll", EntryPoint = "SendMessageTimeoutW", ExactSpelling = true)]
    private static extern IntPtr Send(IntPtr hwnd, uint message, UIntPtr wparam, IntPtr lparam, uint flags, uint milliseconds, out UIntPtr result);
    [DllImport("user32.dll", EntryPoint = "PostMessageW", ExactSpelling = true)]
    private static extern bool Post(IntPtr hwnd, uint message, UIntPtr wparam, IntPtr lparam);

    public static IntPtr OwnedDialog(uint pid, IntPtr owner, string title) {
        IntPtr found = IntPtr.Zero;
        EnumWindows(delegate(IntPtr hwnd, IntPtr ignored) {
            uint actualPid;
            GetWindowThreadProcessId(hwnd, out actualPid);
            if (actualPid != pid || GetWindow(hwnd, 4) != owner || !IsWindowVisible(hwnd)) return true; // GW_OWNER
            var text = new StringBuilder(256);
            GetWindowText(hwnd, text, text.Capacity);
            if (text.ToString() != title) return true;
            found = hwnd;
            return false;
        }, IntPtr.Zero);
        return found;
    }

    public static string Text(IntPtr hwnd) {
        const uint capacity = 1024;
        IntPtr p = Marshal.AllocHGlobal((int)capacity * 2);
        try {
            UIntPtr count;
            if (Send(hwnd, 0x000D, new UIntPtr(capacity), p, 2, 1500, out count) == IntPtr.Zero) return "";
            int len = (int)count.ToUInt64();
            return len >= 0 && len < capacity ? Marshal.PtrToStringUni(p, len) : "";
        } finally { Marshal.FreeHGlobal(p); }
    }

    public static IntPtr FilenameEdit(IntPtr dialog, string suffix) {
        IntPtr found = IntPtr.Zero;
        EnumChildWindows(dialog, delegate(IntPtr hwnd, IntPtr ignored) {
            var cls = new StringBuilder(64);
            GetClassName(hwnd, cls, cls.Capacity);
            if (!IsWindowVisible(hwnd) || !String.Equals(cls.ToString(), "Edit", StringComparison.OrdinalIgnoreCase)) return true;
            string text = Text(hwnd);
            if (!text.StartsWith("SpeedDisc-", StringComparison.Ordinal) || !text.EndsWith(suffix, StringComparison.OrdinalIgnoreCase)) return true;
            found = hwnd;
            return false;
        }, IntPtr.Zero);
        return found;
    }

    public static bool SetText(IntPtr hwnd, string text) {
        IntPtr p = Marshal.StringToHGlobalUni(text);
        try {
            UIntPtr result;
            return Send(hwnd, 0x000C, UIntPtr.Zero, p, 2, 1500, out result) != IntPtr.Zero && result.ToUInt64() != 0;
        } finally { Marshal.FreeHGlobal(p); }
    }

    public static bool PostButton(IntPtr window, uint id) {
        IntPtr button = GetDlgItem(window, (int)id);
        // BN_CLICKED is zero. Post instead of sending so modal dialogs do not block the controller.
        return button != IntPtr.Zero && IsWindowVisible(button) && Post(window, 0x0111, new UIntPtr(id), button);
    }
}
'@

$exe = (Resolve-Path -LiteralPath $Executable).Path
$evidence = (New-Item -ItemType Directory -Force -Path $EvidenceDirectory).FullName
$scratch = (New-Item -ItemType Directory -Path (Join-Path $evidence ('exports-' + [Guid]::NewGuid().ToString('N')))).FullName
$app = Start-Process -FilePath $exe -WorkingDirectory $scratch -ArgumentList '--demo' -PassThru
try {
    $timer = [Diagnostics.Stopwatch]::StartNew()
    do {
        $app.Refresh()
        if ($app.HasExited) { throw 'Synthetic SpeedDisc GUI exited at launch.' }
        $main = $app.MainWindowHandle
        if ($main -ne [IntPtr]::Zero -and $app.MainWindowTitle -match '^SpeedDisc') { break }
        Start-Sleep -Milliseconds 70
    } while ($timer.Elapsed.TotalSeconds -lt 10)
    if ($main -eq [IntPtr]::Zero -or -not $app.WaitForInputIdle(5000)) {
        throw 'Synthetic SpeedDisc GUI did not open a responsive window.'
    }
    $pid32 = [uint32]$app.Id
    $findDialog = {
        param([string] $title)
        $limit = [Diagnostics.Stopwatch]::StartNew()
        do {
            $dialog = [SpeedDiscExportTest]::OwnedDialog($pid32, $main, $title)
            if ($dialog -ne [IntPtr]::Zero) { return $dialog }
            Start-Sleep -Milliseconds 70
        } while ($limit.Elapsed.TotalSeconds -lt 8)
        throw "Owned dialog did not appear: $title"
    }
    # Script acts only on known report controls and dialogs of this --demo PID.
    foreach ($format in @('txt','pdf')) {
        $buttonId = if ($format -eq 'txt') { [uint32]104 } else { [uint32]105 }
        if (-not [SpeedDiscExportTest]::PostButton($main, $buttonId)) { throw "Cannot invoke export $format." }
        $dialog = & $findDialog 'Guardar relatório SpeedDisc'
        $edit = [SpeedDiscExportTest]::FilenameEdit($dialog, ".$format")
        if ($edit -eq [IntPtr]::Zero) { throw "Cannot identify the preset $format filename field." }
        $path = Join-Path $scratch "relatório sintético.$format"
        if (-not [SpeedDiscExportTest]::SetText($edit, $path) -or [SpeedDiscExportTest]::Text($edit) -cne $path) {
            throw "Could not set and read back synthetic $format path."
        }
        if (-not [SpeedDiscExportTest]::PostButton($dialog, 1)) { throw "Cannot submit the $format save dialog." }
        $confirmation = & $findDialog 'SpeedDisc'
        if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw "No $format file at exact Unicode path." }
        if (-not [SpeedDiscExportTest]::PostButton($confirmation, 1)) { throw "Cannot dismiss $format confirmation." }
        [byte[]] $bytes = [IO.File]::ReadAllBytes($path)
        if ($format -eq 'txt') {
            $content = [IO.File]::ReadAllText($path, [Text.Encoding]::UTF8)
            if ($content -notmatch 'Demonstração' -or $content -notmatch 'Disco principal') {
                throw 'The saved TXT did not contain the expected synthetic report.'
            }
        } else {
            if ($bytes.Length -lt 1024 -or [Text.Encoding]::ASCII.GetString($bytes, 0, 5) -ne '%PDF-') {
                throw 'The saved PDF did not contain a valid PDF header.'
            }
            $end = [Math]::Min(40, $bytes.Length)
            if ([Text.Encoding]::ASCII.GetString($bytes, $bytes.Length - $end, $end) -notmatch '%%EOF\s*$') {
                throw 'The saved PDF did not end with its required EOF marker.'
            }
        }
        $limit = [Diagnostics.Stopwatch]::StartNew()
        while ([SpeedDiscExportTest]::IsWindow($confirmation) -and $limit.Elapsed.TotalSeconds -lt 5) { Start-Sleep -Milliseconds 40 }
        if ([SpeedDiscExportTest]::IsWindow($confirmation)) { throw "Could not close $format confirmation." }
        "SAVE $format through original native owned dialog: PASS ($($bytes.Length) bytes)" | Write-Output
    }
    if (-not $app.CloseMainWindow() -or -not $app.WaitForExit(7000) -or $app.ExitCode -ne 0) {
        throw 'Synthetic SpeedDisc GUI did not close cleanly after exports.'
    }
    Write-Output 'Synthetic report only. No maintenance, elevation, accounts or privileges were used or changed.'
} finally {
    $app.Refresh()
    if (-not $app.HasExited) {
        Stop-Process -Id $app.Id -Force
        $null = $app.WaitForExit(5000)
    }
    $app.Dispose()
}
