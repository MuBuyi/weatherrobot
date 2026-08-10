package attendance

import "testing"

func TestSummarize(t *testing.T) {
	records := []CheckinRecord{
		{CheckinType: "下班打卡", CheckinTime: 200},
		{CheckinType: "上班打卡", CheckinTime: 100},
		{CheckinType: "上班打卡", CheckinTime: 110},
	}
	work, off := summarize(records)
	if work != 100 || off != 200 {
		t.Fatalf("summarize() = (%d, %d), want (100, 200)", work, off)
	}
}

func TestSummarizeIgnoresOffBeforeWork(t *testing.T) {
	work, off := summarize([]CheckinRecord{{CheckinType: "下班打卡", CheckinTime: 90}, {CheckinType: "上班打卡", CheckinTime: 100}})
	if work != 100 || off != 0 {
		t.Fatalf("summarize() = (%d, %d), want (100, 0)", work, off)
	}
}
