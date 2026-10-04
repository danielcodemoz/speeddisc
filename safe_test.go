package main

import "testing"

func TestAllowClean(t *testing.T) {
	ok := []struct {
		k CleanKind
		p string
	}{
		{CleanTemp, `C:\Users\Ana\AppData\Local\Temp`},
		{CleanTemp, `C:\Windows\Temp`},
		{CleanTemp, `C:\Users\Ana\AppData\Local\Temp\..\Temp`},
		{CleanWU, `C:\Windows\SoftwareDistribution\Download`},
		{CleanDO, `C:\Windows\ServiceProfiles\NetworkService\AppData\Local\Microsoft\Windows\DeliveryOptimization\Cache`},
		{CleanWER, `C:\ProgramData\Microsoft\Windows\WER\ReportQueue`},
		{CleanWER, `C:\Users\Ana\AppData\Local\Microsoft\Windows\WER\Temp`},
		{CleanThumb, `C:\Users\Ana\AppData\Local\Microsoft\Windows\Explorer`},
	}
	for _, tc := range ok {
		if err := AllowClean(tc.k, tc.p); err != nil {
			t.Fatalf("AllowClean(%d, %s) = %v", tc.k, tc.p, err)
		}
	}
	bad := []struct {
		k CleanKind
		p string
	}{
		{CleanTemp, `C:\`},
		{CleanTemp, `C:\Windows`},
		{CleanTemp, `C:\Windows\System32`},
		{CleanTemp, `C:\Windows\System32\Temp`},
		{CleanTemp, `C:\Windows\Temp\..\System32`},
		{CleanTemp, `C:\Users\Ana\Desktop`},
		{CleanWU, `C:\Windows\SoftwareDistribution`},
		{CleanWU, `C:\Windows\SoftwareDistribution\DataStore`},
		{CleanWU, `C:\Windows\Prefetch`},
		{CleanDO, `C:\Windows\ServiceProfiles\NetworkService\AppData\Local\Microsoft\Windows\DeliveryOptimization`},
		{CleanDO, `C:\Users\Ana\Microsoft\Windows\DeliveryOptimization\Cache`},
		{CleanWER, `C:\ProgramData\Microsoft\Windows\WER`},
		{CleanThumb, `C:\Windows\System32\Microsoft\Windows\Explorer`},
		{CleanTemp, ``},
	}
	for _, tc := range bad {
		if err := AllowClean(tc.k, tc.p); err == nil {
			t.Fatalf("AllowClean(%d, %s) should fail", tc.k, tc.p)
		}
	}
}

func TestStartupDir(t *testing.T) {
	if !IsStartupDir(`C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup`) {
		t.Fatal("common startup")
	}
	if !IsStartupDir(`C:\Users\Ana\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`) {
		t.Fatal("user startup")
	}
	if IsStartupDir(`C:\Users\Ana\Desktop`) {
		t.Fatal("desktop")
	}
}
