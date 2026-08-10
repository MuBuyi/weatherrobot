package attendance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"wechatrobot/internal/weather"
)

const wecomAPIBase = "https://qyapi.weixin.qq.com/cgi-bin"

type Config struct {
	CorpID       string
	Secret       string
	UserIDs      []string
	StateFile    string
	WorkDuration time.Duration
	LunchBreak   time.Duration
	Location     *time.Location
	HTTPTimeout  time.Duration
}

type Manager struct {
	cfg       Config
	client    *http.Client
	mu        sync.Mutex
	token     string
	tokenTill time.Time
}

type tokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type checkinResponse struct {
	ErrCode int             `json:"errcode"`
	ErrMsg  string          `json:"errmsg"`
	Data    []CheckinRecord `json:"checkindata"`
}

type CheckinRecord struct {
	UserID        string `json:"userid"`
	GroupName     string `json:"groupname,omitempty"`
	CheckinType   string `json:"checkin_type"`
	ExceptionType string `json:"exception_type,omitempty"`
	CheckinTime   int64  `json:"checkin_time"`
	LocationTitle string `json:"location_title,omitempty"`
}

type State struct {
	Days map[string]map[string]*EmployeeDay `json:"days"`
}

type EmployeeDay struct {
	UserID          string          `json:"userid"`
	WorkCheckinTime int64           `json:"work_checkin_time,omitempty"`
	OffCheckinTime  int64           `json:"off_checkin_time,omitempty"`
	RemindedAt      int64           `json:"reminded_at,omitempty"`
	Records         []CheckinRecord `json:"records,omitempty"`
}

func New(cfg Config) (*Manager, error) {
	if cfg.CorpID == "" || cfg.Secret == "" || len(cfg.UserIDs) == 0 {
		return nil, errors.New("企业微信打卡配置不完整")
	}
	if cfg.Location == nil {
		cfg.Location = time.Local
	}
	if cfg.StateFile == "" {
		cfg.StateFile = "data/attendance-state.json"
	}
	if cfg.WorkDuration <= 0 {
		cfg.WorkDuration = 8 * time.Hour
	}
	if cfg.LunchBreak < 0 {
		return nil, errors.New("午休时长不能小于 0")
	}
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = 10 * time.Second
	}
	return &Manager{cfg: cfg, client: &http.Client{Timeout: cfg.HTTPTimeout}}, nil
}

func (m *Manager) Poll(now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now = now.In(m.cfg.Location)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, m.cfg.Location)
	records, err := m.fetchCheckins(dayStart, now)
	if err != nil {
		return err
	}

	state, err := m.loadState()
	if err != nil {
		return err
	}
	dayKey := dayStart.Format("2006-01-02")
	if state.Days[dayKey] == nil {
		state.Days[dayKey] = make(map[string]*EmployeeDay)
	}
	byUser := make(map[string][]CheckinRecord)
	for _, record := range records {
		byUser[record.UserID] = append(byUser[record.UserID], record)
	}

	for _, userID := range m.cfg.UserIDs {
		records := byUser[userID]
		sort.Slice(records, func(i, j int) bool { return records[i].CheckinTime < records[j].CheckinTime })
		employee := state.Days[dayKey][userID]
		if employee == nil {
			employee = &EmployeeDay{UserID: userID}
			state.Days[dayKey][userID] = employee
		}
		employee.Records = records
		employee.WorkCheckinTime, employee.OffCheckinTime = summarize(records)
		if employee.WorkCheckinTime == 0 || employee.OffCheckinTime != 0 || employee.RemindedAt != 0 {
			continue
		}
		reminderAt := time.Unix(employee.WorkCheckinTime, 0).In(m.cfg.Location).Add(m.cfg.WorkDuration + m.cfg.LunchBreak)
		if now.Before(reminderAt) {
			continue
		}
		content := fmt.Sprintf("⏰ 下班打卡提醒\n你今天于 %s 完成上班打卡，当前有效工作时间已满 8 小时（已扣除 1.5 小时午休），请记得打卡下班。", time.Unix(employee.WorkCheckinTime, 0).In(m.cfg.Location).Format("15:04"))
		if err := weather.SendWecomMessage(content, []string{userID}); err != nil {
			return fmt.Errorf("提醒员工 %s 失败: %w", userID, err)
		}
		employee.RemindedAt = now.Unix()
	}

	m.removeOldDays(state, dayStart.AddDate(0, 0, -90))
	return m.saveState(state)
}

func summarize(records []CheckinRecord) (int64, int64) {
	var work, off int64
	for _, record := range records {
		if strings.Contains(record.CheckinType, "上班") && (work == 0 || record.CheckinTime < work) {
			work = record.CheckinTime
		}
		if strings.Contains(record.CheckinType, "下班") && record.CheckinTime > off {
			off = record.CheckinTime
		}
	}
	if off != 0 && work != 0 && off < work {
		off = 0
	}
	return work, off
}

func (m *Manager) fetchCheckins(start, end time.Time) ([]CheckinRecord, error) {
	token, err := m.accessToken()
	if err != nil {
		return nil, err
	}
	payload := map[string]any{"opencheckindatatype": 1, "starttime": start.Unix(), "endtime": end.Unix(), "useridlist": m.cfg.UserIDs}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := wecomAPIBase + "/checkin/getcheckindata?access_token=" + token
	resp, err := m.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("请求打卡记录失败: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("打卡接口 HTTP 状态码 %d", resp.StatusCode)
	}
	var result checkinResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("解析打卡接口响应失败: %w", err)
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("打卡接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}
	return result.Data, nil
}

func (m *Manager) accessToken() (string, error) {
	if m.token != "" && time.Now().Before(m.tokenTill) {
		return m.token, nil
	}
	req, err := http.NewRequest(http.MethodGet, wecomAPIBase+"/gettoken", nil)
	if err != nil {
		return "", err
	}
	query := req.URL.Query()
	query.Set("corpid", m.cfg.CorpID)
	query.Set("corpsecret", m.cfg.Secret)
	req.URL.RawQuery = query.Encode()
	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取 access token 失败: %w", err)
	}
	defer resp.Body.Close()
	var result tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析 access token 响应失败: %w", err)
	}
	if result.ErrCode != 0 || result.AccessToken == "" {
		return "", fmt.Errorf("获取 access token 失败: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}
	m.token = result.AccessToken
	m.tokenTill = time.Now().Add(time.Duration(result.ExpiresIn)*time.Second - 5*time.Minute)
	return m.token, nil
}

func (m *Manager) loadState() (*State, error) {
	state := &State{Days: make(map[string]map[string]*EmployeeDay)}
	data, err := os.ReadFile(m.cfg.StateFile)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取打卡状态失败: %w", err)
	}
	if len(data) == 0 {
		return state, nil
	}
	if err := json.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("解析打卡状态失败: %w", err)
	}
	if state.Days == nil {
		state.Days = make(map[string]map[string]*EmployeeDay)
	}
	return state, nil
}

func (m *Manager) saveState(state *State) error {
	if err := os.MkdirAll(filepath.Dir(m.cfg.StateFile), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	temporary := m.cfg.StateFile + ".tmp"
	if err := os.WriteFile(temporary, data, 0600); err != nil {
		return err
	}
	if err := os.Chmod(temporary, 0600); err != nil {
		return err
	}
	return os.Rename(temporary, m.cfg.StateFile)
}

func (m *Manager) removeOldDays(state *State, cutoff time.Time) {
	for day := range state.Days {
		parsed, err := time.ParseInLocation("2006-01-02", day, m.cfg.Location)
		if err == nil && parsed.Before(cutoff) {
			delete(state.Days, day)
		}
	}
}
